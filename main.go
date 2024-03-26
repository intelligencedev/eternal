package main

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	database "eternal/pkg/database"
	"eternal/pkg/llm"
	"eternal/pkg/llm/claude"
	"eternal/pkg/llm/openai"
	"eternal/pkg/sd"
	"eternal/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/afero"
	"github.com/valyala/fasthttp"
)

//go:embed public/* pkg/llm/local/bin/* pkg/sd/sdcpp/build/bin/*
var embedfs embed.FS

var (
	osFS afero.Fs = afero.NewOsFs()
	// memFS afero.Fs = afero.NewMemMapFs()

	chatTurn = 1
	//sqliteDB *database.SQLiteDB

	tools []Tool
)

type WebSocketMessage struct {
	ChatMessage string                 `json:"chat_message"`
	Model       string                 `json:"model"`
	Headers     map[string]interface{} `json:"HEADERS"`
}

// Tool represents a functionality within Eternal.
type Tool struct {
	Name    string
	Enabled bool
}

func main() {
	pterm.DefaultBigText.WithLetters(putils.LettersFromString("ETERNAL")).Render()

	// Load configuration
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current path: %v", err)
	}
	configPath := filepath.Join(currentPath, "config.yml")
	config, err := LoadConfig(osFS, configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create data directory if it doesn't exist
	if err := EnsureDataPath(config); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	// Initialize server files and database
	if _, err := InitServer(config.DataPath); err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}
	sqliteDB, err := database.NewDB(config.DataPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if err := sqliteDB.AutoMigrate(&database.ModelParams{}, &database.ImageModel{}, &database.SelectedModels{}, &database.Chat{}); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// Load model data into database
	modelParams, err := loadModelData(config)
	if err != nil {
		log.Fatalf("Failed to load model data: %v", err)
	}
	if err := database.LoadModelDataToDB(sqliteDB, modelParams); err != nil {
		log.Fatalf("Failed to load model data to database: %v", err)
	}

	// Load image model data into database
	imageModels, err := loadImageModelData(config)
	if err != nil {
		log.Fatalf("Failed to load image model data: %v", err)
	}
	if err := database.LoadImageModelDataToDB(sqliteDB, imageModels); err != nil {
		log.Fatalf("Failed to load image model data to database: %v", err)
	}

	// Initialize tools
	tools = []Tool{
		{Name: "websearch", Enabled: false},
		{Name: "imagegen", Enabled: false},
	}

	// Start the server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	runFrontendServer(ctx, config, modelParams, sqliteDB)
}

// loadModelData loads model data from the configuration and checks if models are downloaded.
func loadModelData(config *AppConfig) ([]database.ModelParams, error) {
	var modelParams []database.ModelParams
	for _, model := range config.LanguageModels {
		if model.Downloads != nil {
			fileName := filepath.Base(model.Downloads[0])
			model.LocalPath = filepath.Join(config.DataPath, "models", model.Name, fileName)
		}

		var downloaded bool
		if _, err := os.Stat(model.LocalPath); err == nil {
			downloaded = true
		}

		modelParams = append(modelParams, database.ModelParams{
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
	return modelParams, nil
}

// loadImageModelData loads image model data from the configuration and checks if models are downloaded.
func loadImageModelData(config *AppConfig) ([]database.ImageModel, error) {
	var imageModels []database.ImageModel
	for _, model := range config.ImageModels {
		if model.Downloads != nil {
			fileName := filepath.Base(model.Downloads[0])
			model.LocalPath = filepath.Join(config.DataPath, "models", model.Name, fileName)
		}

		var downloaded bool
		if _, err := os.Stat(model.LocalPath); err == nil {
			downloaded = true
		}

		imageModels = append(imageModels, database.ImageModel{
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
	return imageModels, nil
}

// runFrontendServer starts the Fiber web server and handles routes.
func runFrontendServer(ctx context.Context, config *AppConfig, modelParams []database.ModelParams, sqliteDB *database.SQLiteDB) {
	// Create HTTP file system
	basePath := filepath.Join(config.DataPath, "web")
	baseFs := afero.NewBasePathFs(osFS, basePath)
	httpFs := afero.NewHttpFs(baseFs)
	engine := html.NewFileSystem(httpFs, ".html")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Eternal v0.1.0",
		BodyLimit:             100 * 1024 * 1024, // 100MB
		DisableStartupMessage: true,
		ServerHeader:          "Eternal",
		PassLocalsToViews:     true,
		Views:                 engine,
		StrictRouting:         true,
		StreamRequestBody:     true,
	})

	// Enable CORS
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

	// Define routes
	app.Get("/", handleIndex)
	app.Get("/config", handleConfig(config))
	app.Post("/tool/:toolName", handleToolToggle)
	app.Get("/openai/models", handleOpenAIModels(config))
	app.Get("/modeldata/:modelName", handleModelData(sqliteDB))
	app.Put("/modeldata/:modelName/downloaded", handleModelDownloaded(sqliteDB))
	app.Post("/modelcards", handleModelCards(sqliteDB))
	app.Post("/model/select", handleModelSelect(sqliteDB))
	app.Get("/models/selected", handleSelectedModels(sqliteDB))
	app.Post("/model/download", handleModelDownload(config, sqliteDB))
	app.Post("/imgmodel/download", handleImageModelDownload(config))
	app.Post("/chattemplates", handleChatTemplates(config))
	app.Post("/chatsubmit", handleChatSubmit(config, sqliteDB))
	app.Get("/chats", handleGetChats(sqliteDB))
	app.Get("/chats/:id", handleGetChatByID(sqliteDB))
	app.Put("/chats/:id", handleUpdateChat(sqliteDB))
	app.Delete("/chats/:id", handleDeleteChat(sqliteDB))
	app.Get("/dpsearch", handleDPSearch)
	app.Get("/sseupdates", handleSSEUpdates)
	app.Get("/ws", websocket.New(handleWebSocket(config, sqliteDB)))
	app.Get("/wsoai", websocket.New(handleOpenAIWebSocket(config, sqliteDB)))
	app.Get("/wsclaude", websocket.New(handleClaudeWebSocket(config, sqliteDB)))

	// Start the server
	go func() {
		<-ctx.Done()
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

// handleIndex serves the main index page.
func handleIndex(c *fiber.Ctx) error {
	return c.Render("templates/index", fiber.Map{})
}

// handleConfig returns the application configuration as JSON.
func handleConfig(config *AppConfig) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(config)
	}
}

// handleToolToggle toggles the enabled state of a tool.
func handleToolToggle(c *fiber.Ctx) error {
	toolName := c.Params("toolName")

	// Find the tool and toggle its state
	for i, tool := range tools {
		if tool.Name == toolName {
			tools[i].Enabled = !tool.Enabled
			return c.JSON(tools[i])
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("Tool not found")
}

// handleOpenAIModels retrieves and filters OpenAI models.
func handleOpenAIModels(config *AppConfig) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		client := openai.NewClient(config.OAIKey)
		modelsResponse, err := openai.GetModels(client)
		if err != nil {
			log.Errorf(err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		// Filter models starting with "gpt"
		var gptModels []string
		for _, model := range modelsResponse.Data {
			if strings.HasPrefix(model.ID, "gpt") {
				gptModels = append(gptModels, model.ID)
			}
		}

		return c.JSON(fiber.Map{
			"object": "list",
			"data":   gptModels,
		})
	}
}

// handleModelData retrieves model data from the database.
func handleModelData(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var model database.ModelParams
		modelName := c.Params("modelName")
		err := db.First("name = ?", modelName).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get model: %v", err)})
		}
		return c.JSON(model)
	}
}

// handleModelDownloaded updates the downloaded status of a model in the database.
func handleModelDownloaded(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		modelName := c.Params("modelName")
		var payload struct {
			Downloaded bool `json:"downloaded"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		err := db.UpdateDownloadedByName(modelName, payload.Downloaded)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update model: %v", err)})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Model 'Downloaded' status updated successfully",
		})
	}
}

// handleModelCards renders the model cards template with data from the database.
func handleModelCards(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var modelParams []database.ModelParams
		err := db.Find(&modelParams)
		if err != nil {
			log.Errorf("Database error: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}
		return c.Render("templates/model", fiber.Map{"models": modelParams})
	}
}

// handleModelSelect adds or removes a model from the user's selection.
func handleModelSelect(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var selection database.SelectedModels
		if err := c.BodyParser(&selection); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Bad request")
		}

		switch selection.Action {
		case "add":
			if err := database.AddSelectedModel(db, selection.ModelName); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		case "remove":
			if err := database.RemoveSelectedModel(db, selection.ModelName); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// handleSelectedModels returns the list of selected model names.
func handleSelectedModels(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		selectedModels, err := database.GetSelectedModels(db)
		if err != nil {
			log.Errorf("Error getting selected models: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		var selectedModelNames []string
		for _, model := range selectedModels {
			selectedModelNames = append(selectedModelNames, model.ModelName)
		}

		return c.JSON(selectedModelNames)
	}
}

// handleModelDownload handles model download requests.
func handleModelDownload(config *AppConfig, db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		modelName := c.Query("model")
		if modelName == "" {
			log.Errorf("Missing parameters for download")
			return c.Status(fiber.StatusBadRequest).SendString("Missing parameters")
		}

		var downloadURL string
		for _, model := range config.LanguageModels {
			if model.Name == modelName {
				downloadURL = model.Downloads[0]
				break
			}
		}

		modelFileName := filepath.Base(downloadURL)
		modelPath := filepath.Join(config.DataPath, "models", modelName, modelFileName)

		// Check if the file exists and partially downloaded
		var partialDownload bool
		if info, err := os.Stat(modelPath); err == nil {
			// Check if the file size is less than the expected size (if available)
			if info.Size() > 0 {
				// Assuming here that we can check the expected file size somehow,
				// e.g., from a database or a config file. If not available, we
				// still try to resume assuming partial download.
				expectedSize, err := llm.GetExpectedFileSize(downloadURL)
				if err != nil {
					log.Errorf("Error getting expected file size: %v", err)
				}
				partialDownload = info.Size() < expectedSize
			}
		}

		// Download or resume the download
		go func() {
			var err error
			if partialDownload {
				pterm.Info.Printf("Resuming download for model: %s\n", modelName)
				err = llm.Download(downloadURL, modelPath)
			} else {
				pterm.Info.Printf("Starting download for model: %s\n", modelName)
				err = llm.Download(downloadURL, modelPath)
			}
			if err != nil {
				log.Errorf("Error in download: %v", err)
			} else {
				// Update the model's downloaded state in the database
				err = db.UpdateByName(modelName, true)
				if err != nil {
					log.Errorf("Failed to update model downloaded state: %v", err)
				}
			}
		}()

		progressErr := fmt.Sprintf("<div class='w-100' id='progress-download-%s' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message' hx-trigger='load'></div>", modelName)
		return c.SendString(progressErr)
	}
}

// handleImageModelDownload handles image model download requests.
func handleImageModelDownload(config *AppConfig) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		modelName := c.Query("model")
		if modelName == "" {
			log.Errorf("Missing parameters for download")
			return c.Status(fiber.StatusBadRequest).SendString("Missing parameters")
		}

		var downloadURL string
		for _, model := range config.ImageModels {
			if model.Name == modelName {
				downloadURL = model.Downloads[0]
				break
			}
		}

		modelFileName := filepath.Base(downloadURL)
		modelPath := filepath.Join(config.DataPath, "models", modelName, modelFileName)

		// Check if the modelPath does not exist and download it if it doesn't
		if _, err := os.Stat(modelPath); err != nil {
			// Start the download in a goroutine
			go func() {
				if err := sd.Download(downloadURL, modelPath); err != nil {
					log.Errorf("Error in download: %v", err)
				}
			}()
		}

		progressErr := "<div name='sse-messages' class='w-100' id='sse-messages' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message'></div>"
		return c.SendString(progressErr)
	}
}

// handleChatTemplates renders the chat templates.
func handleChatTemplates(config *AppConfig) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		modelsFile := fmt.Sprintf("%v/chat-templates.json", config)
		chatTemplates, err := os.ReadFile(modelsFile)
		if err != nil {
			log.Errorf(err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		var chatTemplate []llm.ChatPromptTemplate
		err = json.Unmarshal(chatTemplates, &chatTemplate)
		if err != nil {
			log.Errorf(err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		return c.Render("templates/chattemplates", fiber.Map{"templates": chatTemplate})
	}
}

// handleChatSubmit handles chat submission and renders the chat template.
func handleChatSubmit(config *AppConfig, db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userPrompt := c.FormValue("userprompt")
		var wsRoute string

		selectedModels, err := database.GetSelectedModels(db)
		if err != nil {
			log.Errorf("Error getting selected models: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		if len(selectedModels) > 0 {
			firstModelName := selectedModels[0].ModelName
			switch {
			case strings.HasPrefix(firstModelName, "openai-"):
				pterm.Println("OpenAI model selected")
				wsRoute = "/wsoai"
			case strings.HasPrefix(firstModelName, "claude-"):
				pterm.Println("Claude model selected")
				wsRoute = "/wsclaude"
			default:
				pterm.Println("Local model selected")
				wsRoute = "/ws"
			}
		} else {
			return c.JSON(fiber.Map{"error": "No models selected"})
		}

		turnID := IncrementTurn()
		return c.Render("templates/chat", fiber.Map{
			"username":  config.CurrentUser,
			"message":   userPrompt,
			"assistant": config.AssistantName,
			"model":     selectedModels[0].ModelName,
			"turnID":    turnID,
			"wsRoute":   wsRoute,
			"hosts":     config.ServiceHosts["llm"],
		})
	}
}

// handleGetChats retrieves all chats from the database.
func handleGetChats(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		chats, err := database.GetChats(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chats"})
		}
		return c.Status(fiber.StatusOK).JSON(chats)
	}
}

// handleGetChatByID retrieves a chat by its ID from the database.
func handleGetChatByID(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		chat, err := database.GetChatByID(db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chat"})
		}

		return c.Status(fiber.StatusOK).JSON(chat)
	}
}

// handleUpdateChat updates an existing chat in the database.
func handleUpdateChat(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		chat := new(database.Chat)
		if err := c.BodyParser(chat); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		err = database.UpdateChat(db, id, chat.Prompt, chat.Response, chat.ModelName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update chat"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// handleDeleteChat deletes a chat from the database.
func handleDeleteChat(db *database.SQLiteDB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		err = database.DeleteChat(db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete chat"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// handleDPSearch performs a search on DuckDuckGo and returns the URLs.
func handleDPSearch(c *fiber.Ctx) error {
	urls := []string{}
	query := c.Query("q")
	res := web.SearchDuckDuckGo(query)
	if len(res) == 0 {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving search results")
	}

	// Remove YouTube results
	urls = append(urls, res...)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"urls": urls})
}

// handleSSEUpdates sends download progress updates via Server-Sent Events.
func handleSSEUpdates(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			progress := llm.GetDownloadProgress("sse-progress")
			msg := fmt.Sprintf("data: <div class='progress specific-h-25 m-4' role='progressbar' aria-label='download' aria-valuenow='%s' aria-valuemin='0' aria-valuemax='100'><div class='progress-bar progress-bar-striped progress-bar-animated' style='width: %s;'></div></div><div class='text-center fs-6'>Please refresh this page when the download completes.</br> Downloading...%s</div>\n\n", progress, progress, progress)
			if _, err := w.WriteString(msg); err != nil {
				pterm.Printf("Error writing to stream: %v", err)
				break
			}
			if err := w.Flush(); err != nil {
				pterm.Printf("Error flushing writer: %v", err)
				break
			}
			time.Sleep(2 * time.Second)
		}
	}))

	return nil
}

// handleWebSocket handles WebSocket connections for local GGML models.
func handleWebSocket(config *AppConfig, db *database.SQLiteDB) func(c *websocket.Conn) {
	pterm.Println("Handling WebSocket connection")
	return func(c *websocket.Conn) {
		if c == nil {
			pterm.Error.Println("WebSocket connection is nil")
			return
		}
		defer c.Close()

		// Read initial message
		_, message, err := c.ReadMessage()
		if err != nil {
			pterm.PrintOnError(err)
			return
		}

		// Unmarshal JSON message
		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Error unmarshalling JSON"))
			return
		}

		// Process chat message and tools
		chatMessage := wsMessage.ChatMessage
		document := processTools(chatMessage, config)
		chatMessage = fmt.Sprintf("%s%s", document, chatMessage)

		// Get model details
		var model database.ModelParams
		result := db.First(wsMessage.Model, &model)
		if result.Error != nil {
			log.Errorf("Error getting model %s: %v", wsMessage.Model, result.Error)
			return
		}

		// Create LLM options
		modelOpts := &llm.GGUFOptions{
			Model:         model.Options.Model,
			Prompt:        chatMessage,
			CtxSize:       model.Options.CtxSize,
			Temp:          0.7,
			RepeatPenalty: 1.1,
		}

		// Send completion to WebSocket
		if err := llm.MakeCompletionWebSocket(*c, chatTurn, modelOpts, config.DataPath); err != nil {
			handleWebSocketError(c, chatTurn, wsMessage.Model, err, db)
			return
		}

		chatTurn++
		c.Close()
	}
}

// handleOpenAIWebSocket handles WebSocket connections for OpenAI models.
func handleOpenAIWebSocket(config *AppConfig, db *database.SQLiteDB) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		if c == nil {
			pterm.Error.Println("WebSocket connection is nil")
			return
		}
		defer c.Close()

		// Read initial message
		_, message, err := c.ReadMessage()
		if err != nil {
			pterm.PrintOnError(err)
			return
		}

		// Unmarshal JSON message
		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Error unmarshalling JSON"))
			return
		}

		// Process chat message and tools
		chatMessage := wsMessage.ChatMessage
		document := processTools(chatMessage, config)
		chatMessage = fmt.Sprintf("%s%s", document, chatMessage)

		// Create chat prompt template
		cpt := llm.GetSystemTemplate(chatMessage)

		// Send completion to WebSocket
		if err := openai.StreamCompletionToWebSocket(c, chatTurn, wsMessage.Model, cpt.Messages, 0.7, config.OAIKey); err != nil {
			handleWebSocketError(c, chatTurn, wsMessage.Model, err, db)
			return
		}

		chatTurn++
		c.Close()
	}
}

// handleClaudeWebSocket handles WebSocket connections for Claude models.
func handleClaudeWebSocket(config *AppConfig, db *database.SQLiteDB) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		if c == nil {
			pterm.Error.Println("WebSocket connection is nil")
			return
		}
		defer c.Close()

		// Read initial message
		_, message, err := c.ReadMessage()
		if err != nil {
			pterm.PrintOnError(err)
			return
		}

		// Unmarshal JSON message
		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Error unmarshalling JSON"))
			return
		}

		// Process chat message and tools
		chatMessage := wsMessage.ChatMessage
		document := processTools(chatMessage, config)
		chatMessage = fmt.Sprintf("%s%s", document, chatMessage)

		// Create chat prompt template
		cpt := llm.GetSystemTemplate(chatMessage)

		// Send completion to WebSocket
		if err := claude.StreamCompletionToWebSocket(c, chatTurn, wsMessage.Model, cpt.Messages, config.AnthropicKey); err != nil {
			handleWebSocketError(c, chatTurn, wsMessage.Model, err, db)
			return
		}

		chatTurn++
		c.Close()
	}
}

// processTools extracts URLs and applies tools to the chat message.
func processTools(chatMessage string, config *AppConfig) string {
	var document string

	// Retrieve page content from prompt URLs
	url := web.ExtractURLs(chatMessage)
	if len(url) > 0 {
		pterm.Info.Println("Extracted URLs: ", url)
		document, _ = web.WebGetHandler(url[0])
		document = fmt.Sprintf("%s\nUse the previous information as reference for the following:\n", document)
		chatMessage = fmt.Sprintf("%s%s", document, chatMessage)
	}

	// Apply tools
	for _, tool := range tools {
		pterm.Info.Println("Processing tool: ", tool)
		switch tool.Name {
		case "imagegen":
			if tool.Enabled {
				document = processImageGenTool(chatMessage, config)
			}
		case "websearch":
			if tool.Enabled {
				document = processWebSearchTool(chatMessage)
			}
		}
	}

	return document
}

// processImageGenTool generates an image using the SD tool.
func processImageGenTool(chatMessage string, config *AppConfig) string {
	pterm.Info.Println("Generating image...")
	sdParams := &sd.SDParams{Prompt: chatMessage}
	sd.Text2Image(config.DataPath, sdParams)

	imgHost := config.ServiceHosts["image"]["image_host_1"]
	imgHostURL := fmt.Sprintf("http://%s:%s", imgHost.Host, imgHost.Port)
	timestamp := time.Now().UnixNano()
	imgElement := fmt.Sprintf("<img class='rounded-2' src='%s/public/img/sd_out.png?%d' />", imgHostURL, timestamp)
	return fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>", fmt.Sprint(chatTurn), imgElement)
}

// processWebSearchTool performs a web search and adds results to the document.
func processWebSearchTool(chatMessage string) string {
	var document string
	urls := web.SearchDuckDuckGo(chatMessage)
	for _, url := range urls {
		pterm.Info.Printf("Retrieving %s\n", url)
		document, _ = web.WebGetHandler(url)
		chatMessage = fmt.Sprintf("%s\nUse the previous information as reference to respond to the following query:\n%s", document, chatMessage)
	}
	return document
}

// handleWebSocketError handles errors during WebSocket communication.
func handleWebSocketError(c *websocket.Conn, chatID int, modelName string, err error, db *database.SQLiteDB) {
	pterm.PrintOnError(err)

	// Store the chat in the database
	chat := &database.Chat{
		Prompt:    llm.GetSystemTemplate("").Messages[0].Content,
		Response:  err.Error(),
		ModelName: modelName,
	}

	pterm.Warning.Print("Storing chat in database...")
	if _, err := database.CreateChat(db, chat.Prompt, chat.Response, chat.ModelName); err != nil {
		pterm.Error.Println("Error storing chat in database:", err)
	}

	// Send error message to the client
	formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>", fmt.Sprint(chatID), err.Error())
	if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
		pterm.Error.Println("WebSocket write error:", err)
	}

	chatTurn++
}
