// main.go
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

var (
	currentProject Project
	comfyUICmd     *exec.Cmd
	devMode        bool
)

// main is the entry point of the application.
func main() {
	flag.BoolVar(&devMode, "devmode", false, "Run the application in development mode")
	flag.Parse()

	displayBanner()

	config, err := loadConfig(localFs, "config.yml")
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

// runFrontendServer starts the Fiber frontend server.
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

// createFiberApp initializes and returns a new Fiber application.
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

// handleGracefulShutdown handles the graceful shutdown of the application.
func handleGracefulShutdown(ctx context.Context, app *fiber.App, config *AppConfig, modelParams []ModelParams) {
	<-ctx.Done()

	if devMode {
		cleanupDevMode(config, modelParams)
	}

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}

// cleanupDevMode performs cleanup tasks specific to development mode.
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
