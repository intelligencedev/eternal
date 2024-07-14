// utils.go
package main

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

// logFatalError logs an error message and exits the program with a status code of 1.
// It takes a message string and an error as parameters.
func logFatalError(msg string, err error) {
	pterm.Error.Println(msg, err)
	os.Exit(1)
}

// deleteFile attempts to delete the file or directory at the specified path.
// If an error occurs during deletion, it logs an error message.
func deleteFile(path string) {
	if err := os.RemoveAll(path); err != nil {
		pterm.Error.Printf("Error deleting file: %s\n", path)
	}
}

// fileExists checks if a file or directory exists at the specified path.
// It returns true if the file exists, and false otherwise.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// displayBanner renders a banner with the text "ETERNAL" using pterm.
func displayBanner() {
	_ = pterm.DefaultBigText.WithLetters(putils.LettersFromString("ETERNAL")).Render()
}

// displayProjects displays a table of projects using pterm.
// It takes a slice of Project structs, where each project has a Name and Description.
func displayProjects(projects []Project) {
	var projectData [][]string
	for _, project := range projects {
		projectData = append(projectData, []string{project.Name, project.Description})
	}

	pterm.DefaultTable.WithData(projectData).WithHasHeader().WithStyle(pterm.NewStyle(pterm.FgCyan)).Render()
}
