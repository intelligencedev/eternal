package main

import (
	"context"
	"embed"
	"errors"
	"eternal/pkg/llm"
	"eternal/pkg/sd"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/blevesearch/bleve/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/afero"
)

var (
	//go:embed public/* pkg/llm/local/bin/* pkg/sd/sdcpp/build/bin/*
	embedfs embed.FS

	osFS afero.Fs = afero.NewOsFs()

	chatTurn = 1
	sqliteDB *SQLiteDB

	searchIndex bleve.Index
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
	displayBanner()

	config, err := loadConfig()
	if err != nil {
		pterm.Error.Println("Error loading config:", err)
		os.Exit(1)
	}

	tools := initializeTools(config)
	log.Infof("Enabled tools: %v", tools)

	if err := createDataDirectory(config.DataPath); err != nil {
		pterm.Error.Println("Error creating data directory:", err)
		os.Exit(1)
	}

	if err := initializeServer(config.DataPath); err != nil {
		pterm.Error.Println("Error initializing server:", err)
		os.Exit(1)
	}

	if err := initializeDatabase(config.DataPath); err != nil {
		pterm.Error.Println("Failed to initialize database:", err)
		os.Exit(1)
	}

	if err := initializeSearchIndex(config.DataPath); err != nil {
		pterm.Error.Println("Failed to initialize search index:", err)
		os.Exit(1)
	}

	modelParams, err := loadModelParams(config)
	if err != nil {
		pterm.Error.Println("Failed to load model data to database:", err)
		os.Exit(1)
	}

	imageModels, err := loadImageModels(config)
	if err != nil {
		pterm.Error.Println("Failed to load image model data to database:", err)
		os.Exit(1)
	}

	log.Infof("Loaded %d language models and %d image models", len(modelParams), len(imageModels))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pterm.Info.Printf("Serving frontend on: %s:%s\n", config.ControlHost, config.ControlPort)
	pterm.Info.Println("Press Ctrl+C to stop")

	runFrontendServer(ctx, config, modelParams)

	pterm.Warning.Println("Shutdown signal received")
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

func createDataDirectory(dataPath string) error {
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return os.Mkdir(dataPath, 0755)
	}
	return nil
}

func initializeServer(dataPath string) error {
	_, err := InitServer(dataPath)
	return err
}

func initializeDatabase(dataPath string) error {
	var err error
	sqliteDB, err = NewSQLiteDB(dataPath)
	if err != nil {
		return err
	}

	return sqliteDB.AutoMigrate(&ModelParams{}, &ImageModel{}, &SelectedModels{}, &Chat{})
}

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

func loadImageModels(config *AppConfig) ([]ImageModel, error) {
	var imageModels []ImageModel
	for _, model := range config.ImageModels {
		if model.Downloads != nil {
			fileName := strings.Split(model.Downloads[0], "/")
			model.LocalPath = fmt.Sprintf("%s/models/%s/%s", config.DataPath, model.Name, fileName[len(fileName)-1])
		}

		var downloaded bool
		if _, err := os.Stat(model.LocalPath); err == nil {
			downloaded = true
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

	setupRoutes(app, config, modelParams)

	go func() {
		<-ctx.Done() // Wait for the context to be cancelled
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

func setupRoutes(app *fiber.App, config *AppConfig, modelParams []ModelParams) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("templates/index", fiber.Map{})
	})

	app.Get("/config", func(c *fiber.Ctx) error {
		return c.JSON(config)
	})

	app.Get("/flow", func(c *fiber.Ctx) error {
		return c.Render("templates/flow", fiber.Map{})
	})

	app.Post("/upload", handleUpload(config))

	app.Post("/tool/:toolName", handleToolToggle(config))

	app.Get("/openai/models", handleOpenAIModels(config))

	app.Get("/modeldata/:modelName", handleModelData())

	app.Put("/modeldata/:modelName/downloaded", handleModelDownloadUpdate())

	app.Post("/modelcards", handleModelCards(modelParams))

	app.Post("/model/select", handleModelSelect())

	app.Get("/models/selected", handleSelectedModels())

	app.Post("/model/download", handleModelDownload(config))

	app.Post("/imgmodel/download", handleImgModelDownload(config))

	app.Post("/api/v1/role/:name", handleRoleSelection(config))

	app.Post("/chatsubmit", handleChatSubmit(config))

	app.Get("/chats", handleGetChats())

	app.Get("/chats/:id", handleGetChatByID())

	app.Put("/chats/:id", handleUpdateChat())

	app.Delete("/chats/:id", handleDeleteChat())

	app.Get("/dpsearch", handleDPSearch())

	app.Get("/sseupdates", handleSSEUpdates())

	app.Get("/ws", websocket.New(handleWebSocket(config)))

	app.Get("/wsoai", websocket.New(handleOpenAIWebSocket(config)))

	app.Get("/wsanthropic", websocket.New(handleAnthropicWebSocket(config)))

	app.Get("/wsgoogle", websocket.New(handleGoogleWebSocket(config)))
}
