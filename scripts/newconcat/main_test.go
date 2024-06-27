// main_test.go
package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestFiles(t *testing.T, baseDir string, files map[string]string) {
	for name, content := range files {
		err := ioutil.WriteFile(filepath.Join(baseDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}
}

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

func TestProcessPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files and directories
	files := map[string]string{
		"file1.txt": "Content of file1",
		"file2.go":  "Content of file2",
		"file3.jpg": "Content of file3",
	}
	subDir := filepath.Join(tempDir, "subdir")
	os.Mkdir(subDir, 0755)
	subFiles := map[string]string{
		"subfile1.txt": "Content of subfile1",
		"subfile2.go":  "Content of subfile2",
	}

	setupTestFiles(t, tempDir, files)
	setupTestFiles(t, subDir, subFiles)

	// Test non-recursive processing
	fileContents := make(map[string]string)
	fileStructure := make(map[string][]string)
	err = processPath(context.Background(), tempDir, []string{".txt", ".go"}, false, fileContents, fileStructure)
	if err != nil {
		t.Fatalf("processPath failed: %v", err)
	}

	expectedFiles := []string{"file1.txt", "file2.go"}
	for _, fileName := range expectedFiles {
		if _, exists := fileContents[filepath.Join(tempDir, fileName)]; !exists {
			t.Errorf("Expected file %s not found in fileContents", fileName)
		}
	}

	if _, exists := fileContents[filepath.Join(subDir, "subfile1.txt")]; exists {
		t.Errorf("Unexpected file subfile1.txt found in fileContents")
	}

	// Test recursive processing
	fileContents = make(map[string]string)
	fileStructure = make(map[string][]string)
	err = processPath(context.Background(), tempDir, []string{".txt", ".go"}, true, fileContents, fileStructure)
	if err != nil {
		t.Fatalf("processPath failed: %v", err)
	}

	expectedFiles = []string{"file1.txt", "file2.go"}
	expectedSubFiles := []string{"subfile1.txt", "subfile2.go"}

	for _, fileName := range expectedFiles {
		if _, exists := fileContents[filepath.Join(tempDir, fileName)]; !exists {
			t.Errorf("Expected file %s not found in fileContents", fileName)
		}
	}

	for _, fileName := range expectedSubFiles {
		if _, exists := fileContents[filepath.Join(subDir, fileName)]; !exists {
			t.Errorf("Expected file %s not found in fileContents", fileName)
		}
	}
}

func TestGenerateStructureString(t *testing.T) {
	fileStructure := map[string][]string{
		"/path/to/dir":        {"file1.txt", "file2.go"},
		"/path/to/dir/subdir": {"subfile1.txt", "subfile2.go"},
	}
	fileContents := map[string]string{
		"/path/to/dir/file1.txt":           "Content of file1",
		"/path/to/dir/file2.go":            "Content of file2",
		"/path/to/dir/subdir/subfile1.txt": "Content of subfile1",
		"/path/to/dir/subdir/subfile2.go":  "Content of subfile2",
	}

	expected := "/path/to/dir/\n├── file1.txt\n├── file2.go\n/path/to/dir/subdir/\n├── subfile1.txt\n├── subfile2.go\n"
	formatter := &PlainTextFormatter{}
	result := formatter.Format(fileStructure, fileContents)
	if result != expected {
		t.Errorf("generateStructureString() = %v; want %v", result, expected)
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

	setupTestFiles(t, tempDir, files)

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
