// eternal/main.go - Main entry point for the Eternal application

package main

import (
	"bufio"
	"context"
	"embed"
	"errors"
	"eternal/pkg/llm"
	"eternal/pkg/sd"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/afero"
)

// Embed static files and binaries
//
//go:embed public/* pkg/llm/local/bin/* pkg/sd/sdcpp/build/bin/*
var embedfs embed.FS
var currentProject Project
var comfyUICmd *exec.Cmd

// WebSocketMessage represents the structure of a WebSocket message
type WebSocketMessage struct {
	ChatMessage string                 `json:"chat_message"`
	Model       string                 `json:"model"`
	Headers     map[string]interface{} `json:"HEADERS"`
}

// Tool represents a tool with its name and enabled status
type Tool struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func main() {
	flag.BoolVar(&devMode, "devmode", false, "Run the application in development mode")
	flag.Parse()

	displayBanner()

	// DISABLED due to bug in CUDA
	// Print host information as pterm table
	// hostInfo, err := GetHostInfo()
	// if err != nil {
	// 	pterm.Error.Println("Error getting host information:", err)
	// } else {
	// 	// Convert memory to GB
	// 	hostInfo.Memory.Total = hostInfo.Memory.Total / 1024 / 1024 / 1024
	// 	// Convert ints to strings for pterm table
	// 	pterm.DefaultTable.WithData(pterm.TableData{
	// 		{"OS", hostInfo.OS},
	// 		{"Architecture", hostInfo.Arch},
	// 		{"CPU Cores", fmt.Sprintf("%d", hostInfo.CPUs)},
	// 		{"Memory (GB)", fmt.Sprintf("%d", hostInfo.Memory.Total)},
	// 		{"GPU Model", hostInfo.GPUs[0].Model},
	// 		{"GPU Cores", hostInfo.GPUs[0].TotalNumberOfCores},
	// 		{"Metal Support", hostInfo.GPUs[0].MetalSupport},
	// 	}).Render()
	// }

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		pterm.Error.Println("Error loading config:", err)
		os.Exit(1)
	}

	// Set defaults
	// Set default assistant role
	config.CurrentRoleInstructions = config.AssistantRoles[0].Instructions

	// Initialize tools based on config
	tools := initializeTools(config)

	// If the tool is enabled, print the tool name
	for _, tool := range tools {
		if tool.Enabled {
			pterm.Info.Println("Enabled tool:", tool.Name)
		}
	}

	// Create data directory if it doesn't exist
	if err := createDataDirectory(config.DataPath); err != nil {
		pterm.Error.Println("Error creating data directory:", err)
		os.Exit(1)
	} else {
		// Delete all of the files in the web/public/tmp directory
		tmpDir := filepath.Join(config.DataPath, "web", "public", "tmp")
		if err := os.RemoveAll(tmpDir); err != nil {
			pterm.Error.Println("Error deleting tmp directory:", err)
		}
	}

	// Initialize server
	if err := initializeServer(config.DataPath); err != nil {
		pterm.Error.Println("Error initializing server:", err)
		os.Exit(1)
	}

	pterm.Warning.Println("Server initialized")

	// Initialize database
	if err := initializeDatabase(config); err != nil {
		pterm.Error.Println("Failed to initialize database:", err)
		os.Exit(1)
	}

	currentProject = config.DefaultProjectConfig

	// Create the default project if it doesn't exist
	err = sqliteDB.CreateProject(&currentProject)
	if err != nil {
		pterm.Warning.Println("Default project already exists")
	}

	// List all projects and print to terminal
	projects, err := sqliteDB.ListProjects()
	if err != nil {
		pterm.Error.Println("Failed to list projects:", err)
		os.Exit(1)
	}

	// Convert projects to [][]string
	var projectData [][]string
	for _, project := range projects {
		projectData = append(projectData, []string{project.Name, project.Description})
	}

	// Print the projects as a pterm table
	pterm.DefaultTable.WithData(projectData).WithHasHeader().WithStyle(pterm.NewStyle(pterm.FgCyan)).Render()

	// Initialize search index
	if err := initializeSearchIndex(config.DataPath); err != nil {
		pterm.Error.Println("Failed to initialize search index:", err)
		os.Exit(1)
	}

	// Load model parameters
	modelParams, err := loadModelParams(config)
	if err != nil {
		pterm.Error.Println("Failed to load model data to database:", err)
		os.Exit(1)
	}

	// Prepare data for the pterm table including headers
	tableData := pterm.TableData{
		{"Model Name", "Context Size", "Downloaded"},
	}

	// Loop through model parameters and add each to the table
	for _, param := range modelParams {
		tableData = append(tableData, []string{param.Name, fmt.Sprintf("%d", param.Options.CtxSize), fmt.Sprintf("%t", param.Downloaded)})
	}

	// Print the model parameters as a pterm table
	pterm.DefaultTable.WithData(tableData).WithHasHeader().WithStyle(pterm.NewStyle(pterm.FgCyan)).Render()

	// Load image models
	err = DownloadDefaultImageModel(config)
	if err != nil {
		pterm.Error.Println("Failed to download default image model:", err)
	}

	// Setup context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start ComfyUI
	if err := startComfyUI(ctx, config); err != nil {
		pterm.Error.Println("Failed to start ComfyUI:", err)
		os.Exit(1)
	}

	pterm.Info.Printf("Serving frontend on: %s:%s\n", config.ControlHost, config.ControlPort)
	pterm.Info.Println("Press Ctrl+C to stop")

	// Run frontend server
	runFrontendServer(ctx, config, modelParams)

	pterm.Warning.Println("Shutdown signal received")

	// Stop ComfyUI
	if err := stopComfyUI(ctx); err != nil {
		pterm.Error.Println("Failed to stop ComfyUI:", err)
	}

	os.Exit(0)
}

// displayBanner displays the application banner
func displayBanner() {
	_ = pterm.DefaultBigText.WithLetters(putils.LettersFromString("ETERNAL")).Render()
}

// loadConfig loads the application configuration from a file
func loadConfig() (*AppConfig, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current path: %w", err)
	}

	configPath := filepath.Join(currentPath, "config.yml")
	pterm.Info.Println("Loading config:", configPath)

	return LoadConfig(osFS, configPath)
}

// initializeTools initializes the tools based on the configuration
func initializeTools(config *AppConfig) []Tool {
	var tools []Tool
	if config.Tools.WebGet.Enabled {
		tools = append(tools, Tool{Name: "webget", Enabled: true})
	}
	if config.Tools.WebSearch.Enabled {
		tools = append(tools, Tool{Name: "websearch", Enabled: true})
	}
	return tools
}

// createDataDirectory creates the data directory if it doesn't exist
func createDataDirectory(dataPath string) error {
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return os.Mkdir(dataPath, 0755)
	}
	return nil
}

// initializeServer initializes the server
func initializeServer(dataPath string) error {
	_, err := InitServer(dataPath)
	return err
}

// initializeDatabase initializes the SQLite database
func initializeDatabase(config *AppConfig) error {
	var err error
	sqliteDB, err = NewSQLiteDB(config.DataPath)
	if err != nil {
		return err
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
		return err
	}

	pterm.Warning.Println("Database initialized")

	return nil
}

func setCurrentProject(projectName string) (Project, error) {
	var project Project
	if err := sqliteDB.First(projectName, &project); err != nil {
		return Project{}, err
	}

	// Set the application context to the current project
	currentProject = project
	return project, nil
}

// initializeSearchIndex initializes the search index
func initializeSearchIndex(dataPath string) error {
	searchDB := fmt.Sprintf("%s/search.bleve", dataPath)

	if _, err := os.Stat(searchDB); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		searchIndex, err = bleve.New(searchDB, mapping)
		if err != nil {
			return err
		}
	} else {
		searchIndex, err = bleve.Open(searchDB)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadModelParams loads the model parameters from the configuration
func loadModelParams(config *AppConfig) ([]ModelParams, error) {
	var modelParams []ModelParams
	for _, model := range config.LanguageModels {
		if model.Downloads != nil {
			fileName := strings.Split(model.Downloads[0], "/")
			model.LocalPath = fmt.Sprintf("%s/models/%s/%s", config.DataPath, model.Name, fileName[len(fileName)-1])
		}

		var downloaded bool
		if _, err := os.Stat(model.LocalPath); err == nil {
			downloaded = true
		}

		modelParams = append(modelParams, ModelParams{
			Name:       model.Name,
			Homepage:   model.Homepage,
			GGUFInfo:   model.GGUF,
			Downloaded: downloaded,
			Options: &llm.GGUFOptions{
				Model:         model.LocalPath,
				Prompt:        model.Prompt,
				CtxSize:       model.Ctx,
				Temp:          0.7,
				RepeatPenalty: 1.1,
			},
		})
	}

	if err := LoadModelDataToDB(sqliteDB, modelParams); err != nil {
		return nil, err
	}
	return modelParams, nil
}

// loadImageModels loads the image models from the configuration
func loadImageModels(config *AppConfig) ([]ImageModel, error) {
	var imageModels []ImageModel
	for _, model := range config.ImageModels {
		// Print the download state
		pterm.Info.Printf("Image model: %s, Downloaded: %t\n", model.Name, model.Downloaded)

		if model.Downloads != nil {
			fileName := strings.Split(model.Downloads[0], "/")
			pterm.Info.Printf("File name: %s\n", fileName[len(fileName)-1])
			// /Users/arturoaquino/.eternal-v1/sd/ComfyUI-master/models/checkpoints
			model.LocalPath = fmt.Sprintf("%s/sd/ComfyUI-master/models/checkpoints/%s", config.DataPath, fileName[len(fileName)-1])
		}

		var downloaded bool
		if _, err := os.Stat(model.LocalPath); err == nil {
			downloaded = true
		} else {
			pterm.Warning.Printf("Image model not found: %s\n", model.LocalPath)
		}

		imageModels = append(imageModels, ImageModel{
			Name:       model.Name,
			Homepage:   model.Homepage,
			Prompt:     model.Prompt,
			Downloaded: downloaded,
			Options: &sd.SDParams{
				Model:  model.LocalPath,
				Prompt: model.Prompt,
			},
		})
	}

	if err := LoadImageModelDataToDB(sqliteDB, imageModels); err != nil {
		return nil, err
	}
	return imageModels, nil
}

// runFrontendServer runs the frontend server
func runFrontendServer(ctx context.Context, config *AppConfig, modelParams []ModelParams) {
	basePath := filepath.Join(config.DataPath, "web")
	baseFs := afero.NewBasePathFs(afero.NewOsFs(), basePath)
	httpFs := afero.NewHttpFs(baseFs)
	engine := html.NewFileSystem(httpFs, ".html")

	app := fiber.New(fiber.Config{
		AppName:               "Eternal v0.1.0",
		BodyLimit:             100 * 1024 * 1024, // 100MB, to allow for larger file uploads
		DisableStartupMessage: true,
		ServerHeader:          "Eternal",
		PassLocalsToViews:     true,
		Views:                 engine,
		StrictRouting:         true,
		StreamRequestBody:     true,
	})

	// Setup CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	// Serve static files
	app.Use("/public", filesystem.New(filesystem.Config{
		Root:   httpFs,
		Index:  "index.html",
		Browse: true,
	}))

	app.Static("/", "public")

	// Setup routes
	setupRoutes(app, config, modelParams)

	// Handle graceful shutdown
	go func() {
		<-ctx.Done() // Wait for the context to be cancelled

		if devMode {
			// delete the search index and database
			if err := os.RemoveAll(filepath.Join(config.DataPath, "search.bleve")); err != nil {
				log.Fatalf("Failed to delete search index: %v", err)
			}

			if err := os.RemoveAll(filepath.Join(config.DataPath, "eternaldata.db")); err != nil {
				log.Fatalf("Failed to delete database: %v", err)
			}

			// Loop through the config models and delete the cache
			for _, model := range modelParams {
				if model.Downloaded {
					cachePath := filepath.Join(config.DataPath, "models", model.Name, "cache")

					// First check if the cache file exists
					if _, err := os.Stat(cachePath); err == nil {
						pterm.Warning.Printf("Deleting cache: %s\n", cachePath)

						if err := os.RemoveAll(cachePath); err != nil {
							log.Fatalf("Failed to delete cache: %v", err)
						}
					}
				}
			}
		}

		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
	}()

	addr := fmt.Sprintf("%s:%s", config.ControlHost, config.ControlPort)
	if err := app.Listen(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Frontend server failed: %v", err)
	}

	pterm.Info.Println("Server gracefully shutdown")
}

func startComfyUI(ctx context.Context, config *AppConfig) error {
	comfyPort := config.ServiceHosts["image"]["image_host_1"].Port
	cmdArgs := []string{
		"main.py",
		"--listen",
		"--port", comfyPort,
	}

	comfyUIPath := filepath.Join(config.DataPath, "sd/ComfyUI-master")
	comfyUICmd = exec.CommandContext(ctx, "python", cmdArgs...)
	comfyUICmd.Dir = comfyUIPath

	// Set up pipes for stdout and stderr
	stdout, err := comfyUICmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderr, err := comfyUICmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start ComfyUI
	if err := comfyUICmd.Start(); err != nil {
		return fmt.Errorf("failed to start ComfyUI: %v", err)
	}

	// Monitor ComfyUI output
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Info("ComfyUI: ", scanner.Text())
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Error("ComfyUI Error: ", scanner.Text())
		}
	}()

	// Wait for ComfyUI to be ready
	if err := waitForComfyUI(comfyPort); err != nil {
		return fmt.Errorf("ComfyUI failed to start: %v", err)
	}

	return nil
}

func waitForComfyUI(comfyPort string) error {
	comfyUrl := fmt.Sprintf("http://localhost:%s", comfyPort)
	client := &http.Client{Timeout: 1 * time.Second}
	for i := 0; i < 30; i++ {
		resp, err := client.Get(comfyUrl)
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timeout waiting for ComfyUI to start")
}

func stopComfyUI(ctx context.Context) error {
	if comfyUICmd == nil || comfyUICmd.Process == nil {
		return nil
	}

	// Send SIGTERM
	if err := comfyUICmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to ComfyUI: %v", err)
	}

	// Wait for the process to exit with a timeout
	done := make(chan error, 1)
	go func() {
		done <- comfyUICmd.Wait()
	}()

	select {
	case <-time.After(10 * time.Second):
		// Force kill if it doesn't exit within the timeout
		if err := comfyUICmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill ComfyUI process: %v", err)
		}
		return errors.New("ComfyUI did not exit gracefully, forcefully terminated")
	case err := <-done:
		if err != nil {
			return fmt.Errorf("ComfyUI exited with error: %v", err)
		}
		return nil
	}
}
