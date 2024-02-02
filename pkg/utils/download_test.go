package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	socket "github.com/gorilla/websocket"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// MockFileSystem sets up an afero in-memory file system for testing
func MockFileSystem() afero.Fs {
	appFs := afero.NewMemMapFs()
	return appFs
}

// TestDownloadHandler tests the download handler with WebSocket connection
func TestDownloadHandler(t *testing.T) {
	// Create a Fiber app instance
	app := fiber.New()

	// Setup routes
	app.Get("/download", DownloadHandler)
	app.Use("/ws", websocket.New(WebSocketDownloadHandler))

	// Start Fiber app in a goroutine
	go app.Listen(":3000")
	time.Sleep(100 * time.Millisecond) // give some time for the server to start

	// Prepare a WebSocket client
	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/ws"}
	c, _, err := socket.DefaultDialer.Dial(u.String(), nil)
	assert.Nil(t, err, "Dialing to websocket server should not fail")
	defer c.Close()

	// Declare downloadChan variable
	downloadChan := make(chan Progress)

	// Mock a download and listen for messages
	go func() {
		err := Download("http://localhost:8000/cloudbox_api_schema.json", "public/downloads/cloudbox_api_schema.json", downloadChan)
		assert.Nil(t, err, "Download should not return an error")
	}()

	// Read messages from WebSocket
	_, message, err := c.ReadMessage()
	assert.Nil(t, err, "Reading from WebSocket should not fail")
	assert.NotNil(t, message, "Message should not be nil")

	// Further tests can be added here to check the specific content of the WebSocket messages
	// and validate the progress updates sent through it.
}

// TestEnsureDir tests the ensureDir function for directory creation
func TestEnsureDir(t *testing.T) {
	fs := MockFileSystem()
	aferoFs := &afero.Afero{Fs: fs}

	testDir := "."

	err := ensureDir(testDir)
	assert.Nil(t, err, "ensureDir should not return an error for valid directory path")

	exists, err := aferoFs.DirExists(testDir)
	assert.Nil(t, err, "DirExists should not return an error")
	assert.True(t, exists, "Directory should exist after ensureDir")
}

// TestHTTPGet tests the httpGet function for HTTP requests
func TestHTTPGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Test successful GET request
	resp, err := httpGet(server.URL + "/success")
	assert.Nil(t, err, "httpGet should not return an error for a successful request")
	assert.NotNil(t, resp, "Response should not be nil for a successful request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be OK for a successful request")
	resp.Body.Close() // Close the response body to avoid resource leaks

	// Test failed GET request
	resp, err = httpGet(server.URL + "/notfound")
	assert.NotNil(t, err, "httpGet should return an error for a failed request")
	assert.Nil(t, resp, "Response should be nil for a failed request")
}

// TestDownloadWithProgress tests the downloadWithProgress function for downloading with progress updates
func TestDownloadWithProgress(t *testing.T) {
	progressChannel := make(chan Progress, 10) // Buffered channel to avoid blocking

	// Mock server response
	resp := &http.Response{
		Request:       nil,
		StatusCode:    200,
		Header:        make(http.Header),
		Body:          io.NopCloser(bytes.NewBufferString("test")),
		ContentLength: 40960,
	}

	go func() {
		err := downloadWithProgress(resp, 4096, progressChannel)
		assert.Nil(t, err, "downloadWithProgress should not return error")
		// Handle progress updates in a non-blocking way
		for progress := range progressChannel {
			checkProgressValue(progress, t, progressChannel, "Progress value not expected.")
		}
		close(progressChannel) // Close the channel when done
	}()

}

func checkProgressValue(expected Progress, t *testing.T, progressChannel <-chan Progress, errorMessage string) {
	select {
	case progress := <-progressChannel:
		// Validate progress
		assert.Equal(t, expected, progress, errorMessage)
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for progress update:", errorMessage)
	}
}
