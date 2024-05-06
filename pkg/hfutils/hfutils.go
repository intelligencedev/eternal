package hfutils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// ConcurrentDownloadManager handles downloading a file in parts concurrently.
type ConcurrentDownloadManager struct {
	FileName       string
	URL            string
	Destination    string
	NumParts       int
	TempDir        string
	Sha256Checksum string
	TotalBytes     int64
	TotalLength    int64      // Total length of the file to be downloaded
	BytesMutex     sync.Mutex // Mutex to protect TotalBytes
}

func (dm *ConcurrentDownloadManager) Download() error {
	resp, err := http.Head(dm.URL)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %s", resp.Status)
	}

	lengthStr := resp.Header.Get("Content-Length")
	length, err := strconv.ParseInt(lengthStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid content length: %w", err)
	}
	dm.TotalLength = length // Store the total length

	partSize := int(length) / dm.NumParts
	var wg sync.WaitGroup
	var downloadErr error
	mutex := &sync.Mutex{}

	// Start progress reporting in a goroutine
	go dm.PrintProgress()

	for i := 0; i < dm.NumParts; i++ {
		wg.Add(1)
		start := i * partSize
		end := start + partSize
		if i == dm.NumParts-1 {
			end = int(length)
		}

		go func(partNum, start, end int) {
			defer wg.Done()
			err := dm.downloadPart(partNum, start, end)
			if err != nil {
				mutex.Lock()
				if downloadErr == nil {
					downloadErr = err
				}
				mutex.Unlock()
			}
		}(i, start, end)
	}

	wg.Wait()

	if downloadErr != nil {
		return downloadErr
	}

	if err := dm.mergeParts(); err != nil {
		return fmt.Errorf("failed to merge parts: %w", err)
	}

	if dm.Sha256Checksum != "" {
		if err := dm.verifyChecksum(); err != nil {
			return err
		}
	}

	return nil
}

func (dm *ConcurrentDownloadManager) downloadPart(partNum, start, end int) error {
	req, err := http.NewRequest("GET", dm.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	partPath := filepath.Join(dm.TempDir, fmt.Sprintf("part-%d", partNum))
	partFile, err := os.Create(partPath)
	if err != nil {
		return err
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024) // 32 KB buffer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := partFile.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			dm.BytesMutex.Lock()
			dm.TotalBytes += int64(n)
			dm.BytesMutex.Unlock()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

func (dm *ConcurrentDownloadManager) mergeParts() error {
	finalFile, err := os.Create(dm.Destination)
	if err != nil {
		return err
	}
	defer finalFile.Close()

	for i := 0; i < dm.NumParts; i++ {
		partPath := filepath.Join(dm.TempDir, fmt.Sprintf("part-%d", i))
		partFile, err := os.Open(partPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(finalFile, partFile); err != nil {
			partFile.Close()
			return err
		}

		partFile.Close()
		os.Remove(partPath) // Cleanup part file after merge
	}
	return nil
}

func (dm *ConcurrentDownloadManager) verifyChecksum() error {
	finalPath := filepath.Join(dm.TempDir, dm.FileName)
	file, err := os.Open(finalPath)
	if err != nil {
		return err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if actualChecksum != dm.Sha256Checksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", dm.Sha256Checksum, actualChecksum)
	}

	return nil
}

func (dm *ConcurrentDownloadManager) PrintProgress() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		dm.BytesMutex.Lock()
		downloaded := dm.TotalBytes
		dm.BytesMutex.Unlock()

		percent := float64(downloaded) / float64(dm.TotalLength) * 100
		fmt.Printf("\rDownloading... %.2f%% complete (%d of %d bytes)", percent, downloaded, dm.TotalLength)

		if downloaded >= dm.TotalLength {
			break
		}
	}
}
