package llm

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/pterm/pterm"
)

var (
	downloadProgressMap = make(map[string]DownloadProgress)
	progressMutex       sync.Mutex
	TurnCounter         = 0
)

// DownloadProgress structure to hold download progress information
type DownloadProgress struct {
	Total   int64 `json:"total"`
	Current int64 `json:"current"`
}

// WebModel defines the interface for a model that interacts over the web.
type WebModel interface {
	Connect(endpoint string) (*http.Response, error)
}

// ModelManager defines the interface for managing local models.
type ModelManager interface {
	GetConfig() error
	Download() error
	Delete() error
}

type Model struct {
	Name      string   `yaml:"name"`
	Homepage  string   `yaml:"homepage"`
	Prompt    string   `yaml:"prompt"`
	Ctx       int      `yaml:"ctx"`
	Roles     []string `yaml:"roles"`
	Tags      []string `yaml:"tags,omitempty"`
	GGUF      string   `yaml:"gguf,omitempty"`
	Downloads []string `yaml:"downloads,omitempty"`
	LocalPath string   `yaml:"localPath,omitempty"`
}

type ProgressReader struct {
	Reader        io.Reader
	ProgressBar   *pterm.ProgressbarPrinter
	TotalRead     int64
	ContentLength int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.TotalRead += int64(n)

	// Calculate progress percentage
	//progressPercentage := int64(100.0 * float64(pr.TotalRead) / float64(pr.ContentLength))

	// Update the progress map safely
	progressMutex.Lock()
	downloadProgressMap["sse-progress"] = DownloadProgress{
		Total:   pr.ContentLength,
		Current: pr.TotalRead,
	}
	progressMutex.Unlock()

	pr.ProgressBar.Add(n)

	return n, err
}

func Download(url string, localPath string) error {
	dir := filepath.Dir(localPath)
	// Ensure the directory exists
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open the local file for writing, create it if not exists
	out, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer out.Close()

	// Find out how much has already been downloaded
	fi, err := out.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	size := fi.Size()

	// If already downloaded, no need to download again
	if size > 0 {
		fmt.Printf("Resuming download from byte %d...\n", size)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Range header to request the portion of the file we don't have yet
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", size))

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start file download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status getting file: %s", resp.Status)
	}

	// Seek to the end of the file to start appending data
	_, err = out.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	pterm.Info.Printf("Downloading model:\nURL: %s\nFile: %s\n", url, localPath)

	// Initialize the progress bar
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(int(resp.ContentLength)).WithTitle("Downloading").Start()

	// Wrap the response body in a custom reader that updates the progress bar
	progressReader := &ProgressReader{
		Reader:        resp.Body,
		ProgressBar:   progressBar,
		TotalRead:     size,
		ContentLength: resp.ContentLength + size,
	}

	// Copy the remaining data to the file, updating progress along the way
	_, err = io.Copy(out, progressReader)
	if err != nil {
		return err
	}

	// Ensure the progress bar reflects the complete download
	progressBar.Total = (int(progressReader.TotalRead))

	// Finish the progress bar
	progressBar.Stop()

	pterm.Success.Println("Download completed successfully.")

	return nil
}

func GetDownloadProgress(key string) string {
	progressMutex.Lock()
	defer progressMutex.Unlock()

	// Return the progress for a specific key
	if progress, ok := downloadProgressMap[key]; ok {
		// Calculate progress percentage
		return fmt.Sprintf("%d%%", int64(100.0*float64(progress.Current)/float64(progress.Total)))
	}

	// Fallback if no progress is found
	return "0%"
}

func (m *Model) Delete() error {
	if err := os.Remove(m.LocalPath); err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	return nil
}

func IncrementTurnCounter() {
	TurnCounter++
}

func GetTurnCounter() int {
	return TurnCounter
}
