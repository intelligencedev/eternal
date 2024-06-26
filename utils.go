package main

import (
	"embed"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync/atomic"

	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/afero"
)

var (
	LocalFs        = new(afero.OsFs)
	MemFs          = afero.NewMemMapFs()
	messageCounter int64
)

func InitServer(configPath string) (string, error) {

	// WEB FILES
	webPath := filepath.Join(configPath, "web")
	err := os.MkdirAll(webPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", webPath, err)
	}
	err = CopyFiles(embedfs, "public", webPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy files: %v", err)
	}

	// GGUF FILES
	ggufPath := filepath.Join(configPath, "gguf")
	err = os.MkdirAll(ggufPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", ggufPath, err)
	}
	err = CopyFiles(embedfs, "pkg/llm/local/bin", ggufPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy files: %v", err)
	}

	files, err := os.ReadDir(ggufPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %v", ggufPath, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			err = os.Chmod(filepath.Join(ggufPath, file.Name()), 0755)
			if err != nil {
				return "", fmt.Errorf("failed to set executable permission on file %s: %v", file.Name(), err)
			}
		}
	}

	// IMG GEN
	imgGenPath := filepath.Join(configPath, "sd")
	err = os.MkdirAll(imgGenPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", imgGenPath, err)
	}

	err = CopyFiles(embedfs, "pkg/sd/sdcpp/build/bin", imgGenPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy files: %v", err)
	}

	files, err = os.ReadDir(imgGenPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %v", imgGenPath, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			err = os.Chmod(filepath.Join(imgGenPath, file.Name()), 0755)
			if err != nil {
				return "", fmt.Errorf("failed to set executable permission on file %s: %v", file.Name(), err)
			}
		}
	}

	return configPath, nil
}

func EnsureDataPath(config *AppConfig) error {
	if _, err := os.Stat(config.DataPath); os.IsNotExist(err) {
		return LocalFs.MkdirAll(config.DataPath, os.ModePerm)
	}
	return nil
}

func CopyFiles(fsys embed.FS, srcDir, destDir string) error {
	fileEntries, err := fsys.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", srcDir, err)
	}

	for _, entry := range fileEntries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Create the directory and copy its contents
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", destPath, err)
			}
			if err := CopyFiles(fsys, srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy the file
			fileData, err := fsys.ReadFile(srcPath)
			if err != nil {
				log.Errorf("failed to read file %s: %v", srcPath, err)
				continue // Skip to the next file
			}
			if err := os.WriteFile(destPath, fileData, 0755); err != nil {
				return fmt.Errorf("failed to write file %s: %v", destPath, err)
			}
		}
	}
	return nil
}

// Increments and returns a counter that gets appended to the id for frontend chat elements
func IncrementTurn() int64 {
	return atomic.AddInt64(&messageCounter, 1)
}

// findURLInText searches for a URL in a given text and returns it if found.
// It returns nil if no valid URL is found.
func URLParse(text string) *url.URL {
	// Define a regular expression for finding URLs
	// This is a simple regex for demonstration; it might not cover all URL cases
	re := regexp.MustCompile(`https?://[^\s]+`)

	// Find a URL using the regex
	found := re.FindString(text)
	if found == "" {
		// No URL found
		return nil
	}

	// Parse the URL to validate it and return *url.URL
	parsedURL, err := url.Parse(found)
	if err != nil {
		// The URL is not valid
		return nil
	}

	return parsedURL
}
