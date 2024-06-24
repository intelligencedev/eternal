package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHasMatchingExtension(t *testing.T) {
	tests := []struct {
		filePath   string
		extensions []string
		expected   bool
	}{
		{"test.txt", []string{".txt", ".go"}, true},
		{"test.go", []string{".txt", ".go"}, true},
		{"test.jpg", []string{".txt", ".go"}, false},
		{"test.TXT", []string{".txt", ".go"}, true},
		{"test", []string{".txt", ".go"}, false},
	}

	for _, test := range tests {
		result := hasMatchingExtension(test.filePath, test.extensions)
		if result != test.expected {
			t.Errorf("hasMatchingExtension(%s, %v) = %v; want %v", test.filePath, test.extensions, result, test.expected)
		}
	}
}

func TestMainFunctionality(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	files := map[string]string{
		"file1.txt": "Content of file1",
		"file2.go":  "Content of file2",
		"file3.jpg": "Content of file3",
	}

	for name, content := range files {
		err := ioutil.WriteFile(filepath.Join(tempDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Set up command-line arguments
	os.Args = []string{"cmd", "-paths=" + tempDir, "-types=.txt,.go", "-output=test_output.txt"}

	// Run main function
	main()

	// Check if output file was created
	outputContent, err := ioutil.ReadFile("test_output.txt")
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check if output contains expected content
	expectedFiles := []string{"file1.txt", "file2.go"}
	for _, fileName := range expectedFiles {
		if !strings.Contains(string(outputContent), "--- BEGIN "+filepath.Join(tempDir, fileName)) {
			t.Errorf("Output doesn't contain expected content for %s", fileName)
		}
		if !strings.Contains(string(outputContent), "Content of "+strings.TrimSuffix(fileName, filepath.Ext(fileName))) {
			t.Errorf("Output doesn't contain expected content for %s", fileName)
		}
		if !strings.Contains(string(outputContent), "--- END "+filepath.Join(tempDir, fileName)) {
			t.Errorf("Output doesn't contain expected content for %s", fileName)
		}
	}

	// Check if output doesn't contain unexpected content
	if strings.Contains(string(outputContent), "file3.jpg") {
		t.Errorf("Output contains unexpected content for file3.jpg")
	}

	// Clean up
	os.Remove("test_output.txt")
}
