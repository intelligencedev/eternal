// main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// FileProcessor interface for processing files
type FileProcessor interface {
	Process(path string) error
	GetContent() string
}

// TextFileProcessor processes text files
type TextFileProcessor struct {
	content string
}

func (tfp *TextFileProcessor) Process(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	tfp.content = string(content)
	return nil
}

func (tfp *TextFileProcessor) GetContent() string {
	return tfp.content
}

// Formatter interface for formatting output
type Formatter interface {
	Format(fileStructure map[string][]string, fileContents map[string]string) string
}

// PlainTextFormatter formats output as plain text
type PlainTextFormatter struct{}

func (ptf *PlainTextFormatter) Format(fileStructure map[string][]string, fileContents map[string]string) string {
	var builder strings.Builder
	for path, files := range fileStructure {
		builder.WriteString(path + "/\n")
		for _, file := range files {
			filePath := filepath.Join(path, file)
			if _, exists := fileContents[filePath]; exists {
				builder.WriteString("├── " + file + "\n")
			}
		}
	}
	return builder.String()
}

func main() {
	// Define command-line flags
	paths := flag.String("paths", "", "Comma-separated list of paths")
	types := flag.String("types", "", "Comma-separated list of file types")
	outputFile := flag.String("output", "output.txt", "Output file path")
	recursive := flag.Bool("recursive", false, "Traverse subdirectories")
	ignorePattern := flag.String("ignore", "", "Pattern to ignore in file names")
	flag.Parse()

	// Validate input
	if *paths == "" || *types == "" {
		log.Fatal("Paths and types must be provided")
	}

	// Split paths and types into slices
	pathList := strings.Split(*paths, ",")
	typeList := strings.Split(*types, ",")

	// Create a map to store file contents
	fileContents := make(map[string]string)
	// Create a map to store file structure
	fileStructure := make(map[string][]string)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Loop over each path
	for _, path := range pathList {
		err := processPath(ctx, path, typeList, *recursive, *ignorePattern, fileContents, fileStructure)
		if err != nil {
			log.Printf("Error processing path %s: %v\n", path, err)
		}
	}

	// Generate the file/folder structure string
	formatter := &PlainTextFormatter{}
	structureString := formatter.Format(fileStructure, fileContents)

	// Concatenate file contents with delimiters
	var concatenated strings.Builder
	concatenated.WriteString(structureString)
	for filePath, content := range fileContents {
		concatenated.WriteString(fmt.Sprintf("--- BEGIN %s ---\n", filePath))
		concatenated.WriteString(content)
		concatenated.WriteString(fmt.Sprintf("--- END %s ---\n\n", filePath))
	}

	// Write concatenated content to the output file
	err := ioutil.WriteFile(*outputFile, []byte(concatenated.String()), 0644)
	if err != nil {
		log.Fatalf("Error writing to output file: %v\n", err)
	}

	fmt.Printf("Concatenated content written to %s\n", *outputFile)
}

func processPath(ctx context.Context, path string, typeList []string, recursive bool, ignorePattern string, fileContents map[string]string, fileStructure map[string][]string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range files {
			filePath := filepath.Join(path, file.Name())

			if file.IsDir() {
				if recursive {
					err := processPath(ctx, filePath, typeList, recursive, ignorePattern, fileContents, fileStructure)
					if err != nil {
						return err
					}
				}
			} else if hasMatchingExtension(filePath, typeList) && !shouldIgnore(file.Name(), ignorePattern) {
				processor := &TextFileProcessor{}
				err := processor.Process(filePath)
				if err != nil {
					return err
				}

				// Store file contents in the map
				fileContents[filePath] = processor.GetContent()
				fileStructure[path] = append(fileStructure[path], file.Name())
			}
		}

		return nil
	}
}

func hasMatchingExtension(filePath string, extensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, e := range extensions {
		if strings.ToLower(e) == ext {
			return true
		}
	}
	return false
}

func shouldIgnore(fileName, ignorePattern string) bool {
	if ignorePattern == "" {
		return false
	}
	return strings.Contains(fileName, ignorePattern)
}
