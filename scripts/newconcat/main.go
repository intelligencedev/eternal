// go run main.go -paths="/path1,/path2" -types=".txt,.go" -output="result.txt" -recursive=true

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func main() {
	// Define command-line flags
	paths := flag.String("paths", "", "Comma-separated list of paths")
	types := flag.String("types", "", "Comma-separated list of file types")
	outputFile := flag.String("output", "output.txt", "Output file path")
	recursive := flag.Bool("recursive", false, "Traverse subdirectories")
	flag.Parse()

	// Split paths and types into slices
	pathList := strings.Split(*paths, ",")
	typeList := strings.Split(*types, ",")

	// Create a map to store file contents
	fileContents := make(map[string]string)

	// Loop over each path
	for _, path := range pathList {
		err := processPath(path, typeList, *recursive, fileContents)
		if err != nil {
			fmt.Printf("Error processing path %s: %v\n", path, err)
		}
	}

	// Concatenate file contents with delimiters
	var concatenated strings.Builder
	for filePath, content := range fileContents {
		concatenated.WriteString(fmt.Sprintf("--- BEGIN %s ---\n", filePath))
		concatenated.WriteString(content)
		concatenated.WriteString(fmt.Sprintf("--- END %s ---\n\n", filePath))
	}

	// Write concatenated content to the output file
	err := ioutil.WriteFile(*outputFile, []byte(concatenated.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}

	fmt.Printf("Concatenated content written to %s\n", *outputFile)
}

func processPath(path string, typeList []string, recursive bool, fileContents map[string]string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())

		if file.IsDir() {
			if recursive {
				err := processPath(filePath, typeList, recursive, fileContents)
				if err != nil {
					return err
				}
			}
		} else if hasMatchingExtension(filePath, typeList) {
			// Read file contents
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			// Store file contents in the map
			fileContents[filePath] = string(content)
		}
	}

	return nil
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
