package main

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"eternal/pkg/embeddings"
	"eternal/pkg/llm"
	"eternal/pkg/llm/anthropic"
	"eternal/pkg/llm/google"
	"eternal/pkg/llm/openai"
	"eternal/pkg/sd"
	"eternal/pkg/web"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

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
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

var (
	//go:embed public/* pkg/llm/local/bin/* pkg/sd/sdcpp/build/bin/*
	embedfs embed.FS

	osFS  afero.Fs = afero.NewOsFs()
	memFS afero.Fs = afero.NewMemMapFs()

	chatTurn    = 1
	sqliteDB    *SQLiteDB
	searchIndex bleve.Index

	tools []Tool
)

type WebSocketMessage struct {
	ChatMessage string                 `json:"chat_message"`
	Model       string                 `json:"model"`
	Headers     map[string]interface{} `json:"HEADERS"`
}

// Define a Tool struct
type Tool struct {
	Name    string
	Enabled bool
}

func main() {
	_ = pterm.DefaultBigText.WithLetters(putils.LettersFromString("ETERNAL")).Render()

	// LOG SETTINGS
	//log.SetOutput(io.Discard)

	// Log configuration
	//log.SetLevel(log.LevelDebug)

	//zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// TODO: Check if external dependencies are installed and if not, install them
	// Such as Chromium, Docker, etc. For now, only Chromium is required for the web tool.

	// CONFIG
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current path: %v", err)
	}

	configPath := filepath.Join(currentPath, "config.yml")

	pterm.Info.Println("Loading config:", configPath)

	config, err := LoadConfig(osFS, configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		os.Exit(1)
	}

	// Populate tools based on the configuration
	var tools []Tool

	// Print the tools enabled in the config

	if config.Tools.WebGet.Enabled {
		tools = append(tools, Tool{Name: "webget", Enabled: true})
	}
	if config.Tools.WebSearch.Enabled {
		tools = append(tools, Tool{Name: "websearch", Enabled: true})
	}

	pterm.Info.Sprintf("GPU Layers: %s\n", config.ServiceHosts["llm"]["llm_host_1"].GgufGPULayers)

	if _, err := os.Stat(config.DataPath); os.IsNotExist(err) {
		err = os.Mkdir(config.DataPath, 0755)
		if err != nil {
			pterm.Error.Println("Error creating data directory:", err)
			os.Exit(1)
		}
	}

	_, err = InitServer(config.DataPath)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	sqliteDB, err = NewSQLiteDB(config.DataPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	err = sqliteDB.AutoMigrate(&ModelParams{}, &ImageModel{}, &SelectedModels{}, &Chat{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	searchDB := fmt.Sprintf("%s/search.bleve", config.DataPath)

	// If the database exists, open it, else create a new one
	if _, err := os.Stat(searchDB); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		searchIndex, err = bleve.New(searchDB, mapping)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		searchIndex, err = bleve.Open(searchDB)
		if err != nil {
			log.Fatalf("Failed to open search index: %v", err)
		}
	}

	// Instantiate ModelParams then populate it with each model from the config
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
		log.Fatalf("Failed to load model data to database: %v", err)
	}

	// Instantiate ImageModel then populate it with each model from the config
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
		log.Fatalf("Failed to load image model data to database: %v", err)
	}

	// Populate tools
	websearch := Tool{Name: "websearch", Enabled: false}
	imagegen := Tool{Name: "imagegen", Enabled: false}

	// Append tools to the list
	tools = append(tools, websearch, imagegen)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pterm.Info.Printf("Serving fronted on: %s:%s\n", config.ControlHost, config.ControlPort)
	pterm.Info.Println("Press Ctrl+C to stop")

	runFrontendServer(ctx, config, modelParams)

	pterm.Warning.Println("Shutdown signal received")

	os.Exit(0)
}

func runFrontendServer(ctx context.Context, config *AppConfig, modelParams []ModelParams) {

	// Create a http fs
	basePath := filepath.Join(config.DataPath, "web")
	baseFs := afero.NewBasePathFs(osFS, basePath)
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

	// CORS allow all origins for now while mvp dev mode
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

	// main route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("templates/index", fiber.Map{})
	})

	app.Get("/config", func(c *fiber.Ctx) error {
		// Return the app config as JSON
		return c.JSON(config)
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		pterm.Warning.Println("Uploads route hit")

		// Parse the multipart form
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		// Get the files from the form
		files := form.File["file"]

		// Loop through the files
		for _, file := range files {
			// Save the file to the datapath web/uploads directory
			filename := filepath.Join(config.DataPath, "web", "uploads", file.Filename)
			pterm.Warning.Printf("Uploading file: %s\n", filename)
			err := c.SaveFile(file, filename)
			if err != nil {
				return err
			}

			// Log the uploaded file
			log.Infof("Uploaded file: %s", file.Filename)
		}

		// Return a success response
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("%d files uploaded successfully", len(files)),
		})
	})

	// route to enable or disable a tool
	app.Post("/tool/:toolName", func(c *fiber.Ctx) error {
		toolName := c.Params("toolName")

		// Find the index of the tool and a flag indicating if it's found
		var index int
		found := false
		for i, t := range tools {
			if t.Name == toolName {
				index = i
				found = true
				break
			}
		}

		// If the tool is not found, return a 404 error
		if !found {
			return c.Status(404).SendString("Tool not found")
		}

		// Toggle the Enabled status of the tool
		tools[index].Enabled = !tools[index].Enabled

		// Return the updated tool as JSON
		return c.JSON(tools[index])
	})

	app.Get("/openai/models", func(c *fiber.Ctx) error {
		client := openai.NewClient(config.OAIKey)
		modelsResponse, err := openai.GetModels(client)

		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

		// Filter the models to include only those with IDs starting with 'gpt'
		// This needs to be changed to a different method later. Using the name
		// is a future bug waiting to happen.
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
	})

	app.Get("/modeldata/:modelName", func(c *fiber.Ctx) error {
		var model ModelParams

		modelName := c.Params("modelName")
		err := sqliteDB.First(modelName, &model)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("Model not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		return c.JSON(model)
	})

	app.Put("/modeldata/:modelName/downloaded", func(c *fiber.Ctx) error {
		modelName := c.Params("modelName")
		var payload struct {
			Downloaded bool `json:"downloaded"`
		}

		// Parse the JSON body to extract the 'downloaded' status
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		// Update the 'Downloaded' status of the model in the database using its name
		err := sqliteDB.UpdateDownloadedByName(modelName, payload.Downloaded)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update model: %v", err)})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Model 'Downloaded' status updated successfully",
		})
	})

	app.Post("/modelcards", func(c *fiber.Ctx) error {
		// Retrieve all models from the database
		err := sqliteDB.Find(&modelParams)

		if err != nil {
			log.Errorf("Database error: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		// Render the template with the models data
		return c.Render("templates/model", fiber.Map{"models": modelParams})
	})

	app.Post("/model/select", func(c *fiber.Ctx) error {
		var selection SelectedModels

		if err := c.BodyParser(&selection); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Bad request")
		}

		// Add or remove the model from the selection based on the action
		if selection.Action == "add" {
			if err := AddSelectedModel(sqliteDB.db, selection.ModelName); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		} else if selection.Action == "remove" {
			if err := RemoveSelectedModel(sqliteDB.db, selection.ModelName); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/models/selected", func(c *fiber.Ctx) error {
		selectedModels, err := GetSelectedModels(sqliteDB.db)

		if err != nil {
			log.Errorf("Error getting selected models: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		var selectedModelNames []string

		for _, model := range selectedModels {
			selectedModelNames = append(selectedModelNames, model.ModelName)
		}

		return c.JSON(selectedModelNames)
	})

	app.Post("/model/download", func(c *fiber.Ctx) error {
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
				err = sqliteDB.UpdateDownloadedByName(modelName, true)
				if err != nil {
					log.Errorf("Failed to update model downloaded state: %v", err)
				}
			}
		}()

		progressErr := fmt.Sprintf("<div class='w-100' id='progress-download-%s' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message' hx-trigger='load'></div>", modelName)

		return c.SendString(progressErr)
	})

	app.Post("/imgmodel/download", func(c *fiber.Ctx) error {
		modelName := c.Query("model")

		var downloadURL string
		for _, model := range config.ImageModels {
			if model.Name == modelName {
				downloadURL = model.Downloads[0]
			}
		}

		modelFileName := strings.Split(downloadURL, "/")[len(strings.Split(downloadURL, "/"))-1]

		if modelName == "" {
			log.Errorf("Missing parameters for download")
			return c.Status(fiber.StatusBadRequest).SendString("Missing parameters")
		}

		modelPath := fmt.Sprintf("%s/models/%s/%s", config.DataPath, modelName, modelFileName)

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
	})

	app.Post("/chattemplates", func(c *fiber.Ctx) error {
		modelsFile := fmt.Sprintf("%v/chat-templates.json", config)

		chatTemplates, err := os.ReadFile(modelsFile)
		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

		var chatTemplate []llm.ChatPromptTemplate
		err = json.Unmarshal(chatTemplates, &chatTemplate)

		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

		return c.Render("templates/chattemplates", fiber.Map{"templates": chatTemplate})
	})

	app.Post("/chatsubmit", func(c *fiber.Ctx) error {

		// userPrompt is the message displayed in the chat view
		userPrompt := c.FormValue("userprompt")

		var wsroute string

		selectedModels, err := GetSelectedModels(sqliteDB.db)
		if err != nil {
			log.Errorf("Error getting selected models: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		if len(selectedModels) > 0 {
			firstModelName := selectedModels[0].ModelName

			// Check if the first model name starts with "openai-"
			if strings.HasPrefix(firstModelName, "openai-") {
				wsroute = "/wsoai"
			} else if strings.HasPrefix(firstModelName, "google-") {
				wsroute = "/wsgoogle"
			} else if strings.HasPrefix(firstModelName, "anthropic-") {
				wsroute = "/wsanthropic"
			} else {
				wsroute = fmt.Sprintf("ws://%s:%s/ws", config.ServiceHosts["llm"]["llm_host_1"].Host, config.ServiceHosts["llm"]["llm_host_1"].Port)
			}

		} else {
			// return error
			return c.JSON(fiber.Map{"error": "No models selected"})
		}

		// Generate unique ID
		turnID := IncrementTurn()

		return c.Render("templates/chat", fiber.Map{
			"username":  config.CurrentUser,
			"message":   userPrompt, // This is the message that will be displayed in the chat
			"assistant": config.AssistantName,
			"model":     selectedModels[0].ModelName,
			"turnID":    turnID,
			"wsRoute":   wsroute,
			"hosts":     config.ServiceHosts["llm"],
		})
	})

	// Retrieve all chats from sqlite database
	app.Get("/chats", func(c *fiber.Ctx) error {
		chats, err := GetChats(sqliteDB.db)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chats"})
		}

		return c.Status(fiber.StatusOK).JSON(chats)
	})

	// Retrieve a single chat from sqlite database by id
	app.Get("/chats/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		chat, err := GetChatByID(sqliteDB.db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chat"})
		}

		return c.Status(fiber.StatusOK).JSON(chat)
	})

	// Update a single chat in sqlite database by id
	app.Put("/chats/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		chat := new(Chat)

		if err := c.BodyParser(chat); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		err = UpdateChat(sqliteDB.db, id, chat.Prompt, chat.Response, chat.ModelName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update chat"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// Delete a single chat in sqlite database by id
	app.Delete("/chats/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		err = DeleteChat(sqliteDB.db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete chat"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// Multi web page retrieval via local ChromeDP
	app.Get("/dpsearch", func(c *fiber.Ctx) error {
		urls := []string{}
		query := c.Query("q")
		res := web.SearchDDG(query)

		if len(res) == 0 {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving search results")
		}

		// Remove youtube results
		urls = append(urls, res...)

		// Send results as json object
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"urls": urls})
	})

	app.Get("/sseupdates", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for {
				// Get updated download progress
				progress := llm.GetDownloadProgress("sse-progress")

				// Format message for SSE
				msg := fmt.Sprintf("data: <div class='progress specific-h-25 m-4' role='progressbar' aria-label='download' aria-valuenow='%s' aria-valuemin='0' aria-valuemax='100'><div class='progress-bar progress-bar-striped progress-bar-animated' style='width: %s;'></div></div><div class='text-center fs-6'>Please refresh this page when the download completes.</br> Downloading...%s</div>\n\n", progress, progress, progress)

				// Write the message
				if _, err := w.WriteString(msg); err != nil {
					pterm.Printf("Error writing to stream: %v", err)
					break
				}
				if err := w.Flush(); err != nil {
					pterm.Printf("Error flushing writer: %v", err)
					break
				}

				time.Sleep(2 * time.Second) // Adjust the sleep time as necessary
			}
		}))

		return nil
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		handleWebSocket(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			// Process the message
			//cpt := llm.GetSystemTemplate(chatMessage)
			//fullPrompt := cpt.Messages[0].Content + "\n" + chatMessage

			// Get the details of the first model from database
			var model ModelParams
			err := sqliteDB.First(wsMessage.Model, &model)
			if err != nil {
				log.Errorf("Error getting model %s: %v", wsMessage.Model, err)
				return err
			}

			promptTemplate := model.Options.Prompt

			// Replace {user} with the chat Message
			fullPrompt := strings.ReplaceAll(promptTemplate, "{prompt}", chatMessage)

			// Replace {system} with the system message
			fullPrompt = strings.ReplaceAll(fullPrompt, "{system}", "You are a helpful AI assistant that responds in well structured markdown format. Do not repeat your instructions. Do not deviate from the topic. Begin all responses with 'Sure thing!' and end with 'Is there anything else I can help you with?'")

			modelOpts := &llm.GGUFOptions{
				NGPULayers:    config.ServiceHosts["llm"]["llm_host_1"].GgufGPULayers,
				Model:         model.Options.Model,
				Prompt:        fullPrompt,
				CtxSize:       model.Options.CtxSize,
				Temp:          0.1, // Prefer lower temperature for more controlled responses for now
				RepeatPenalty: 1.1,
				TopP:          1.0, // Prefer greedy decoding for now
				TopK:          1.0, // Prefer greedy decoding for now
			}

			// Search the search index for the chat message
			// searchResults, err := search.Search(searchIndex, chatMessage)
			// if err != nil {
			// 	log.Errorf("Error searching index: %v", err)
			// }

			// search for some text
			// query := bleve.NewMatchQuery(chatMessage)
			// search := bleve.NewSearchRequest(query)
			// searchResults, err := searchIndex.Search(search)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return err
			// }
			// pterm.Info.Println(searchResults)

			return llm.MakeCompletionWebSocket(*c, chatTurn, modelOpts, config.DataPath)
		})
	}))

	app.Get("/wsoai", websocket.New(func(c *websocket.Conn) {
		apiKey := config.OAIKey

		handleWebSocket(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			// Check if embeddings.db exists
			// if _, err := os.Stat(filepath.Join(config.DataPath, "embeddings.db")); os.IsNotExist(err) {
			// 	pterm.Warning.Println("embeddings.db does not exist. Generating embeddings...")
			// 	embeddings.GenerateEmbeddingChat(chatMessage, config.DataPath)
			// }

			cpt := llm.GetSystemTemplate(chatMessage)
			return openai.StreamCompletionToWebSocket(c, chatTurn, "gpt-4-turbo", cpt.Messages, 0.7, apiKey)
		})
	}))

	app.Get("/wsanthropic", websocket.New(func(c *websocket.Conn) {
		apiKey := config.AnthropicKey

		handleAnthropicWS(c, apiKey, chatTurn)
	}))

	app.Get("/wsgoogle", websocket.New(func(c *websocket.Conn) {
		apiKey := config.GoogleKey

		handleWebSocket(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			return google.StreamGeminiResponseToWebSocket(c, chatTurn, chatMessage, apiKey)
		})
	}))

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

func handleWebSocket(c *websocket.Conn, config *AppConfig, processMessage func(WebSocketMessage, string) error) {
	if c == nil {
		pterm.Error.Println("WebSocket connection is nil")
		return
	}
	defer c.Close()

	// Read the initial message
	_, message, err := c.ReadMessage()
	if err != nil {
		pterm.PrintOnError(err)
		return
	}

	// Unmarshal the JSON message
	var wsMessage WebSocketMessage
	err = json.Unmarshal(message, &wsMessage)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("Error unmarshalling JSON"))
		return
	}

	// Extract the chat_message value
	chatMessage := wsMessage.ChatMessage

	// Perform tool workflow and update chatMessage
	chatMessage = performToolWorkflow(c, config, chatMessage)

	// Process the message using the provided function
	res := processMessage(wsMessage, chatMessage)
	if res != nil {
		pterm.Warning.Println(res)

		if config.Tools.Memory.Enabled {
			err = storeChat(sqliteDB.db, config, chatMessage, res.Error(), wsMessage.Model)
			if err != nil {
				pterm.PrintOnError(err)
			}
		}

		// Increment the chat turn counter
		chatTurn = chatTurn + 1
		pterm.Warning.Println("Chat turn:", chatTurn)
		return
	}
}

func performToolWorkflow(c *websocket.Conn, config *AppConfig, chatMessage string) string {
	// Begin tool workflow. Tools will add context to the submitted message for
	// the model to use. Document is the abstraction that will hold that context.
	var document string

	if config.Tools.Memory.Enabled {
		topN := config.Tools.Memory.TopN // retrieve top N results. Adjust based on context size.
		topEmbeddings := embeddings.Search(config.DataPath, "embeddings.db", chatMessage, topN)

		var documents []string
		var documentString string
		if len(topEmbeddings) > 0 {
			for _, topEmbedding := range topEmbeddings {
				documents = append(documents, topEmbedding.Word)
			}
			documentString = strings.Join(documents, " ")

			pterm.Info.Println("Retrieving memory content...")
			document = fmt.Sprintf("%s\n%s", document, documentString)

			// Replace new lines with spaces
			document = strings.ReplaceAll(document, "\n\n", "\n")
		} else {
			pterm.Info.Println("No memory content found...")
		}
	}

	if config.Tools.WebGet.Enabled {
		url := web.ExtractURLs(chatMessage)
		if len(url) > 0 {
			pterm.Info.Println("Retrieving page content...")

			document, _ = web.WebGetHandler(url[0])
		}
	}

	if config.Tools.WebSearch.Enabled {

		topN := config.Tools.WebSearch.TopN // retrieve top N results. Adjust based on context size.

		pterm.Info.Println("Searching the web...")

		urls := web.SearchDDG(chatMessage)

		pterm.Warning.Printf("URLs to fetch: %v\n", urls)

		if len(urls[:topN]) > 0 {
			for _, url := range urls[:topN] {

				pterm.Info.Printf("Fetching URL: %s\n", url)

				page, err := web.WebGetHandler(url)
				if err != nil {
					pterm.PrintOnError(err)
				}

				document = fmt.Sprintf("%s\n%s", document, page)
			}
		}
	}

	// TODO: Legacy tool workflow, still need to expose this as a proper config item
	for _, tool := range tools {
		if tool.Name == "imagegen" && tool.Enabled {
			pterm.Info.Println("Generating image...")

			sdParams := &sd.SDParams{Prompt: chatMessage}

			// Call the sd tool
			sd.Text2Image(config.DataPath, sdParams)

			// Return the image to the client
			timestamp := time.Now().UnixNano() // Get the current timestamp in nanoseconds
			imgElement := fmt.Sprintf("<img class='rounded-2' src='public/img/sd_out.png?%d' />", timestamp)
			formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>", fmt.Sprint(chatTurn), imgElement)
			if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
				pterm.PrintOnError(err)
				return chatMessage
			}

			// Increment the chat turn counter
			chatTurn = chatTurn + 1

			return chatMessage
		}
	}

	//Remove http(s) links from the document so we do not retrieve them unintentionally
	document = web.RemoveUrls(document)

	chatMessage = fmt.Sprintf("%s Reference the previous information and respond to the following task or question:\n%s", document, chatMessage)

	pterm.Error.Println("Tool workflow complete")

	return chatMessage
}

func storeChat(db *gorm.DB, config *AppConfig, prompt, response, modelName string) error {
	// Generate embeddings
	pterm.Warning.Println("Generating embeddings for chat...")
	//turnMemoryText := prompt + "\n" + response

	// err := embeddings.GenerateEmbeddingForTask(searchIndex, "chat", prompt, "txt", 2048, 500, config.DataPath)
	// if err != nil {
	// 	pterm.Error.Println("Error generating embeddings:", err)
	// 	return err
	// }

	err := embeddings.GenerateEmbeddingForTask(searchIndex, "chat", response, "txt", 2048, 500, config.DataPath)
	if err != nil {
		pterm.Error.Println("Error generating embeddings:", err)
		return err
	}

	pterm.Warning.Print("Storing chat in database...")
	if _, err := CreateChat(db, prompt, response, modelName); err != nil {
		pterm.Error.Println("Error storing chat in database:", err)
		return err
	}

	return nil
}

func handleAnthropicWS(c *websocket.Conn, apiKey string, chatID int) {
	// Read the initial message
	_, message, err := c.ReadMessage()
	if err != nil {
		pterm.PrintOnError(err)
		return
	}

	// Unmarshal the JSON message
	var wsMessage WebSocketMessage
	err = json.Unmarshal(message, &wsMessage)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("Error unmarshalling JSON"))
		return
	}

	// Extract the chat_message value
	chatMessage := wsMessage.ChatMessage

	messages := []anthropic.Message{
		{Role: "user", Content: chatMessage},
	}

	res := anthropic.StreamCompletionToWebSocket(c, chatID, "claude-3-opus-20240229", messages, 0.5, apiKey)
	if res != nil {
		pterm.Error.Println("Error in anthropic completion:", res)
	}

	chatTurn = chatTurn + 1

	return // Return to close the connection
}
