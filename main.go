package main

import (
	"bufio"
	"context"
	"embed"
	"errors"
	"eternal/internal/python"
	"eternal/pkg/llm"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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

//go:embed public/* pkg/llm/local/bin/* pkg/sd/sdcpp/build/bin/*
var embedfs embed.FS

var (
	currentProject Project
	comfyUICmd     *exec.Cmd
	//devMode        bool
	//sqliteDB       *SQLiteDB
	//searchIndex    bleve.Index
)

type WebSocketMessage struct {
	ChatMessage string                 `json:"chat_message"`
	Model       string                 `json:"model"`
	Headers     map[string]interface{} `json:"HEADERS"`
}

type Tool struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func main() {
	flag.BoolVar(&devMode, "devmode", false, "Run the application in development mode")
	flag.Parse()

	displayBanner()

	config, err := loadConfig()
	if err != nil {
		logFatalError("Error loading config", err)
	}

	initializeApplication(config)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := startComfyUI(ctx, config); err != nil {
		logFatalError("Failed to start ComfyUI", err)
	}

	pterm.Info.Printf("Serving frontend on: %s:%s\n", config.ControlHost, config.ControlPort)
	pterm.Info.Println("Press Ctrl+C to stop")

	modelParams, err := loadModelParams(config)
	if err != nil {
		logFatalError("Failed to load model parameters", err)
	}

	runFrontendServer(ctx, config, modelParams)

	pterm.Warning.Println("Shutdown signal received")

	if err := stopComfyUI(ctx); err != nil {
		pterm.Error.Println("Failed to stop ComfyUI:", err)
	}

	os.Exit(0)
}

func displayBanner() {
	_ = pterm.DefaultBigText.WithLetters(putils.LettersFromString("ETERNAL")).Render()
}

func loadConfig() (*AppConfig, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current path: %w", err)
	}

	configPath := filepath.Join(currentPath, "config.yml")
	pterm.Info.Println("Loading config:", configPath)

	return LoadConfig(osFS, configPath)
}

func initializeApplication(config *AppConfig) {
	initializeTools(config)
	createDataDirectory(config.DataPath)
	initializeServer(config.DataPath)
	initializeDatabase(config)
	initializeDefaultProject(config)
	initializeSearchIndex(config.DataPath)
	downloadDefaultImageModel(config)
}

func initializeTools(config *AppConfig) {
	tools := []Tool{
		{Name: "webget", Enabled: config.Tools.WebGet.Enabled},
		{Name: "websearch", Enabled: config.Tools.WebSearch.Enabled},
	}

	for _, tool := range tools {
		if tool.Enabled {
			pterm.Info.Println("Enabled tool:", tool.Name)
		}
	}
}

func createDataDirectory(dataPath string) {
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		logFatalError("Error creating data directory", err)
	}

	tmpDir := filepath.Join(dataPath, "web", "public", "tmp")
	if err := os.RemoveAll(tmpDir); err != nil {
		pterm.Error.Println("Error deleting tmp directory:", err)
	}
}

func initializeServer(dataPath string) {
	if _, err := InitServer(dataPath); err != nil {
		logFatalError("Error initializing server", err)
	}
	pterm.Warning.Println("Server initialized")
}

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

func displayProjects(projects []Project) {
	var projectData [][]string
	for _, project := range projects {
		projectData = append(projectData, []string{project.Name, project.Description})
	}

	pterm.DefaultTable.WithData(projectData).WithHasHeader().WithStyle(pterm.NewStyle(pterm.FgCyan)).Render()
}

func initializeSearchIndex(dataPath string) {
	searchDB := filepath.Join(dataPath, "search.bleve")

	var err error
	if _, err := os.Stat(searchDB); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		searchIndex, err = bleve.New(searchDB, mapping)
	} else {
		searchIndex, err = bleve.Open(searchDB)
	}

	if err != nil {
		logFatalError("Failed to initialize search index", err)
	}
}

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

func createModelParam(model *llm.Model, dataPath string) ModelParams {
	var localPath string
	if model.Downloads != nil {
		fileName := filepath.Base(model.Downloads[0])
		localPath = filepath.Join(dataPath, "models", model.Name, fileName)
	}

	downloaded := fileExists(localPath)

	return ModelParams{
		Name:       model.Name,
		Homepage:   model.Homepage,
		GGUFInfo:   model.GGUF,
		Downloaded: downloaded,
		Options: &llm.GGUFOptions{
			Model:         localPath,
			Prompt:        model.Prompt,
			CtxSize:       model.Ctx,
			Temp:          0.7,
			RepeatPenalty: 1.1,
		},
	}
}

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

func downloadDefaultImageModel(config *AppConfig) {
	if err := DownloadDefaultImageModel(config); err != nil {
		pterm.Error.Println("Failed to download default image model:", err)
	}
}

func runFrontendServer(ctx context.Context, config *AppConfig, modelParams []ModelParams) {
	app := createFiberApp(config)
	setupRoutes(app, config, modelParams)

	go handleGracefulShutdown(ctx, app, config, modelParams)

	addr := fmt.Sprintf("%s:%s", config.ControlHost, config.ControlPort)
	if err := app.Listen(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Frontend server failed: %v", err)
	}

	pterm.Info.Println("Server gracefully shutdown")
}

func createFiberApp(config *AppConfig) *fiber.App {
	basePath := filepath.Join(config.DataPath, "web")
	baseFs := afero.NewBasePathFs(afero.NewOsFs(), basePath)
	httpFs := afero.NewHttpFs(baseFs)
	engine := html.NewFileSystem(httpFs, ".html")

	app := fiber.New(fiber.Config{
		AppName:               "Eternal v0.1.0",
		BodyLimit:             100 * 1024 * 1024,
		DisableStartupMessage: true,
		ServerHeader:          "Eternal",
		PassLocalsToViews:     true,
		Views:                 engine,
		StrictRouting:         true,
		StreamRequestBody:     true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	app.Use("/public", filesystem.New(filesystem.Config{
		Root:   httpFs,
		Index:  "index.html",
		Browse: true,
	}))

	app.Static("/", "public")

	return app
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App, config *AppConfig, modelParams []ModelParams) {
	<-ctx.Done()

	if devMode {
		cleanupDevMode(config, modelParams)
	}

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}

func cleanupDevMode(config *AppConfig, modelParams []ModelParams) {
	deleteFile(filepath.Join(config.DataPath, "search.bleve"))
	deleteFile(filepath.Join(config.DataPath, "eternaldata.db"))

	for _, model := range modelParams {
		if model.Downloaded {
			cachePath := filepath.Join(config.DataPath, "models", model.Name, "cache")
			if fileExists(cachePath) {
				pterm.Warning.Printf("Deleting cache: %s\n", cachePath)
				deleteFile(cachePath)
			}
		}
	}
}

func startComfyUI(ctx context.Context, config *AppConfig) error {
	comfyPort := config.ServiceHosts["image"]["image_host_1"].Port
	comfyUIPath := filepath.Join(config.DataPath, "sd/ComfyUI-master")

	if err := installComfyUIRequirements(comfyUIPath); err != nil {
		return err
	}

	if err := runComfyUI(ctx, comfyUIPath, comfyPort); err != nil {
		return err
	}

	return waitForComfyUI(comfyPort)
}

func installComfyUIRequirements(comfyUIPath string) error {
	requirementsPath := filepath.Join(comfyUIPath, "requirements.txt")
	output, err := python.ExecuteScript("-m", "pip", "install", "-r", requirementsPath)
	if err != nil {
		return fmt.Errorf("failed to install ComfyUI requirements: %v\nOutput: %s", err, output)
	}
	pterm.Info.Println("ComfyUI requirements installed")
	return nil
}

func runComfyUI(ctx context.Context, comfyUIPath, comfyPort string) error {
	cmdArgs := []string{
		filepath.Join(comfyUIPath, "main.py"),
		"--listen",
		"--port", comfyPort,
		"--force-fp16",
		"--use-split-cross-attention",
	}

	comfyUICmd = exec.CommandContext(ctx, "python", cmdArgs...)
	comfyUICmd.Dir = comfyUIPath

	stdout, err := comfyUICmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderr, err := comfyUICmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	if err := comfyUICmd.Start(); err != nil {
		return fmt.Errorf("failed to start ComfyUI: %v", err)
	}

	go monitorComfyUIOutput(stdout, stderr)

	return nil
}

func monitorComfyUIOutput(stdout, stderr io.Reader) {
	go scanAndLog(stdout, "ComfyUI: ", log.Info)
	go scanAndLog(stderr, "ComfyUI: ", log.Error)
}

func scanAndLog(r io.Reader, prefix string, logFunc func(...interface{})) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		logFunc(prefix, scanner.Text())
	}
}

func waitForComfyUI(comfyPort string) error {
	comfyUrl := fmt.Sprintf("http://localhost:%s", comfyPort)
	client := &http.Client{Timeout: 1 * time.Second}
	for i := 0; i < 120; i++ {
		if resp, err := client.Get(comfyUrl); err == nil {
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

	if err := comfyUICmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to ComfyUI: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- comfyUICmd.Wait()
	}()

	select {
	case <-time.After(10 * time.Second):
		if err := comfyUICmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill ComfyUI process: %v", err)
		}
		return errors.New("ComfyUI did not exit gracefully, forcefully terminated")
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to wait for ComfyUI process: %v", err)
		}
		return nil
	}
}

func logFatalError(msg string, err error) {
	pterm.Error.Println(msg, err)
	os.Exit(1)
}

func deleteFile(path string) {
	if err := os.RemoveAll(path); err != nil {
		pterm.Error.Printf("Error deleting file: %s\n", path)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
