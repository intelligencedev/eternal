// Package main provides the initialization and setup functions for the application.
package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/pterm/pterm"
)

// embedFS embeds the necessary files for the application.
//
//go:embed public/* pkg/llm/local/bin/*
var embedFS embed.FS

// initializeApplication initializes the application with the given configuration.
func initializeApplication(config *AppConfig) {
	createDataDirectory(config.DataPath)
	initializeServer(config.DataPath)
	initializeDatabase(config)
	initializeDefaultProject(config)
	initializeSearchIndex(config.DataPath)
	downloadDefaultImageModel(config)
}

// createDataDirectory creates the data directory and removes the temporary directory if it exists.
func createDataDirectory(dataPath string) {
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		logFatalError("Error creating data directory", err)
	}

	tmpDir := filepath.Join(dataPath, "web", "public", "tmp")
	if err := os.RemoveAll(tmpDir); err != nil {
		pterm.Error.Println("Error deleting tmp directory:", err)
	}
}

// initializeServer initializes the server by setting up necessary directories and files.
func initializeServer(dataPath string) {
	if _, err := initServer(dataPath); err != nil {
		logFatalError("Error initializing server", err)
	}
	pterm.Warning.Println("Server initialized")
}

// initializeDatabase initializes the SQLite database and performs auto-migration.
func initializeDatabase(config *AppConfig) {
	var err error
	sqliteDB, err = NewSQLiteDB(config.DataPath)
	if err != nil {
		logFatalError("Failed to initialize database", err)
	}

	err = sqliteDB.AutoMigrate(
		&Project{},
		&ModelParams{},
		&ImageModel{},
		&SelectedModels{},
		&Chat{},
		&URLTracking{},
		&Assistant{},
	)
	if err != nil {
		logFatalError("Failed to auto-migrate database", err)
	}

	pterm.Warning.Println("Database initialized")
}

// initializeDefaultProject initializes the default project based on the configuration.
func initializeDefaultProject(config *AppConfig) {
	currentProject = config.DefaultProjectConfig
	err := sqliteDB.CreateProject(&currentProject)
	if err != nil {
		pterm.Warning.Println("Default project already exists")
	}

	projects, err := sqliteDB.ListProjects()
	if err != nil {
		logFatalError("Failed to list projects", err)
	}

	displayProjects(projects)
}

// initializeSearchIndex initializes the search index using Bleve.
func initializeSearchIndex(dataPath string) {
	searchDB := filepath.Join(dataPath, "search.bleve")

	var err error
	if _, err := os.Stat(searchDB); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		searchIndex, err = bleve.New(searchDB, mapping)
		if err != nil {
			logFatalError("Failed to create search index", err)
		}
	} else {
		searchIndex, err = bleve.Open(searchDB)
		if err != nil {
			logFatalError("Failed to open search index", err)
		}
	}

	if err != nil {
		logFatalError("Failed to initialize search index", err)
	}
}

// loadModelParams loads the model parameters from the configuration and saves them to the database.
func loadModelParams(config *AppConfig) ([]ModelParams, error) {
	var modelParams []ModelParams
	for _, model := range config.LanguageModels {
		modelParam := createModelParam(&model, config.DataPath)
		modelParams = append(modelParams, modelParam)
	}

	if err := LoadModelDataToDB(sqliteDB, modelParams); err != nil {
		return nil, err
	}

	displayModelParams(modelParams)
	return modelParams, nil
}

// displayModelParams displays the model parameters in a table format.
func displayModelParams(modelParams []ModelParams) {
	tableData := [][]string{{"Model Name", "Context Size", "Downloaded"}}
	for _, param := range modelParams {
		tableData = append(tableData, []string{
			param.Name,
			fmt.Sprintf("%d", param.Options.CtxSize),
			fmt.Sprintf("%t", param.Downloaded),
		})
	}
	pterm.DefaultTable.WithData(tableData).WithHasHeader().WithStyle(pterm.NewStyle(pterm.FgCyan)).Render()
}

// downloadDefaultImageModel downloads the default image model based on the configuration.
func downloadDefaultImageModel(config *AppConfig) {
	if err := DownloadDefaultImageModel(config); err != nil {
		pterm.Error.Println("Failed to download default image model:", err)
	}
}

// initServer initializes the server by setting up necessary directories and files.
func initServer(configPath string) (string, error) {
	if err := setupDirectory(configPath, "web", "public"); err != nil {
		return "", err
	}

	if err := setupDirectory(configPath, "gguf", "pkg/llm/local/bin"); err != nil {
		return "", err
	}

	if err := setupComfyUI(configPath); err != nil {
		return "", err
	}

	if err := setupImpactPack(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

// setupDirectory creates a directory and copies files into it.
func setupDirectory(configPath, dirName, srcDir string) error {
	dirPath := filepath.Join(configPath, dirName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	if err := CopyFiles(embedFS, srcDir, dirPath); err != nil {
		return fmt.Errorf("failed to copy files to %s: %w", dirPath, err)
	}

	return setExecutablePermissions(dirPath)
}

// installPythonRequirements installs Python requirements from a requirements.txt file.
func installPythonRequirements(reqPath string) error {
	cmd := exec.Command("pip3", "install", "-r", reqPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// setExecutablePermissions sets executable permissions on all files in a directory.
func setExecutablePermissions(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return os.Chmod(path, 0755)
	})
}
