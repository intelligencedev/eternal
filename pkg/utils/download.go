package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

// Progress represents the download progress
type Progress struct {
	Percentage float64 `json:"percentage"`
}

// Download downloads the model from the specified URL and saves it to the specified path.
// It sends progress updates through the provided channel.
func Download(url, localPath string, progressChan chan<- Progress) error {
	// Create directory
	if err := ensureDir(filepath.Dir(localPath)); err != nil {
		return err
	}

	// Initialize HTTP request
	resp, err := httpGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create output file
	out, err := fs.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download with progress, sending updates through the channel
	return downloadWithProgress(resp, resp.ContentLength, progressChan)
}

// DownloadHandler creates a WebSocket endpoint for download progress updates.
func DownloadHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) { // Check if it's a WebSocket request
		return c.Next() // Proceed to the WebSocket handler
	}
	return fiber.ErrUpgradeRequired // If not WebSocket, send error
}

// WebSocketDownloadHandler handles WebSocket connections for download progress.
func WebSocketDownloadHandler(c *websocket.Conn) {
	// Create a progress channel
	progressChan := make(chan Progress)
	defer close(progressChan)

	// Start the download in a new goroutine
	go func() {
		err := Download("http://localhost:8000/sd_output.png", ".", progressChan)
		if err != nil {
			// Handle error, maybe send an error message via WebSocket before closing
			_ = c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			c.Close()
			return
		}
	}()

	// Listen for progress updates and send them to the WebSocket client
	for progress := range progressChan {
		jsonProgress, err := json.Marshal(progress)
		if err != nil {
			// Handle error
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, jsonProgress); err != nil {
			// Handle WebSocket write error
			break
		}
	}
}

func ensureDir(dir string) error {
	if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

func httpGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to start file download: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close() // ensure the body is closed on error
		return nil, fmt.Errorf("bad status getting file: %s", resp.Status)
	}
	return resp, nil
}

func downloadWithProgress(resp *http.Response, totalSize int64, progressChan chan<- Progress) error {
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(int(resp.ContentLength)).WithTitle("Downloading model").Start()

	buffer := make([]byte, 4096) // Buffer size

	// Send progress updates
	var downloaded int64
	// Open file for writing
	file, err := fs.Create("public/downloads/test")
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer func(file afero.File) {
		err := file.Close()
		if err != nil {
			err = fmt.Errorf("failed to close file: %w", err)
			fmt.Println(err)
			return
		}
	}(file)

	// Read from response body
	for {
		bytesRead, readErr := resp.Body.Read(buffer)
		if bytesRead > 0 {
			// Write only the number of bytes read to the file
			_, writeErr := file.Write(buffer[:bytesRead])
			if writeErr != nil {
				return fmt.Errorf("failed to write data to file: %w", writeErr)
			}
			progressBar.Add(bytesRead)
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return fmt.Errorf("failed to read data: %w", readErr)
		}

		// Calculate and send progress
		downloaded += int64(bytesRead)
		progress := Progress{Percentage: float64(downloaded) / float64(totalSize)}
		progressChan <- progress
	}

	_, err = progressBar.Stop()
	if err != nil {
		return err
	}

	return nil
}
