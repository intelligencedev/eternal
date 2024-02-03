package main

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"eternal/pkg/llm"
	"eternal/pkg/llm/openai"
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
	//go:embed public/* pkg/llm/local/bin/*
	embedfs embed.FS

	osFS  afero.Fs = afero.NewOsFs()
	memFS afero.Fs = afero.NewMemMapFs()

	chatTurn = 1
	sqliteDB *SQLiteDB
)

type WebSocketMessage struct {
	ChatMessage string                 `json:"chat_message"`
	Model       string                 `json:"model"`
	Headers     map[string]interface{} `json:"HEADERS"`
}

func main() {
	_ = pterm.DefaultBigText.WithLetters(putils.LettersFromString("ETERNAL")).Render()

	// LOG SETTINGS
	//log.SetOutput(io.Discard)

	// Log configuration
	//log.SetLevel(log.LevelDebug)

	//zerolog.SetGlobalLevel(zerolog.InfoLevel)

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

	err = sqliteDB.AutoMigrate(&ModelParams{}, &SelectedModels{}, &Chat{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pterm.Info.Printf("Serving fronted on: %s:%s\n", config.ControlHost, config.ControlPort)
	pterm.Info.Println("Press Ctrl+C to stop")

	runFrontendServer(ctx, config, modelParams)

	pterm.Warning.Println("Shutdown signal received")

	os.Exit(0)
}

func runFrontendServer(ctx context.Context, config *AppConfig, modelParams []ModelParams) {

	engine := html.NewFileSystem(http.FS(embedfs), ".html")

	httpFs := afero.NewHttpFs(memFS)

	app := fiber.New(fiber.Config{
		AppName:               "Eternal v0.1.0",
		BodyLimit:             100 * 1024 * 1024, // 100MB, to allow for larger file uploads
		DisableStartupMessage: true,
		ServerHeader:          "Eternal",
		PassLocalsToViews:     true,
		Views:                 engine,
		StrictRouting:         true,
	})

	// CORS allow all origins for now while mvp dev mode
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	app.Use("/public", filesystem.New(filesystem.Config{
		Root:  httpFs.Dir("./public"),
		Index: "index.html",
	}))

	app.Static("/", "./public")

	// main route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("public/templates/index", fiber.Map{})
	})

	app.Get("/config", func(c *fiber.Ctx) error {
		// Return the app config as JSON
		return c.JSON(config)
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

	app.Post("/modelcards", func(c *fiber.Ctx) error {
		// Retrieve all models from the database
		err := sqliteDB.Find(&modelParams)
		if err != nil {
			log.Errorf("Database error: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		// Render the template with the models data
		return c.Render("public/templates/model", fiber.Map{"models": modelParams})
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

		var downloadURL string
		for _, model := range config.LanguageModels {
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

		// Start the download in a goroutine
		go func() {
			if err := llm.Download(downloadURL, modelPath); err != nil {
				log.Errorf("Error in download: %v", err)
			}
		}()

		return c.SendStatus(fiber.StatusOK)
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

		return c.Render("public/templates/chattemplates", fiber.Map{"templates": chatTemplate})
	})

	app.Post("/chatsubmit", func(c *fiber.Ctx) error {

		// Parse URL from c.FormValue("userprompt")

		// userPrompt is the message displayed in the chat view
		userPrompt := c.FormValue("userprompt")

		// fullPromt will be the message sent to LLM after it is modified by workflows
		var fullPrompt string
		url := web.ExtractURLs(userPrompt)

		var document string

		if len(url) > 0 {
			document, _ = web.WebGetHandler(url[0])
			document = fmt.Sprintf("%s\nUse the previous unformation as reference for the following:\n", document)
			fullPrompt = fmt.Sprintf("%s%s", document, userPrompt)
		} else {
			fullPrompt = userPrompt
		}

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
			} else {
				wsroute = "/ws"
			}
		} else {
			// return error
			return c.JSON(fiber.Map{"error": "No models selected"})
		}

		//fmt.Println(wsroute)
		//pterm.Info.Println(userPrompt)
		//pterm.Info.Println(fullPrompt)

		// Generate unique ID
		turnID := IncrementTurn()

		return c.Render("public/templates/chat", fiber.Map{
			"username":  config.CurrentUser,
			"prompt":    userPrompt, // This is the message that will be displayed in the chat
			"message":   fullPrompt, // This is the message that will be sent to the AI
			"assistant": config.AssistantName,
			"model":     selectedModels[0].ModelName,
			"turnID":    turnID,
			"wsRoute":   wsroute,
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

	// SSE endpoint
	app.Get("/sseupdates", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			fmt.Println("WRITER")
			var i int
			for {
				i++

				// Get current time
				t := time.Now().Format("2006-01-02 15:04:05")

				// Write message
				msg := fmt.Sprintf("data: <div>%s</div>", t)
				fmt.Fprintf(w, "%s\n\n", msg)
				fmt.Println(msg)

				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				time.Sleep(2 * time.Second)
			}
		}))

		return nil
	})

	// Multi web page retrieval via serpapi
	// app.Get("/search", func(c *fiber.Ctx) error {
	// 	return web.SearchHandler(c)
	// })

	// Multi web page retrieval via local ChromeDP
	// app.Get("/dpsearch", func(c *fiber.Ctx) error {
	// 	return web.SearchChromeDPHandler(c)
	// })

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
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

		// // Process the message
		cpt := llm.GetSystemTemplate(chatMessage)
		fullPrompt := cpt.Messages[0].Content + "\n" + chatMessage

		// // Get the details of the first model from database
		var model ModelParams
		err = sqliteDB.First(wsMessage.Model, &model)
		if err != nil {
			log.Errorf("Error getting model %s: %v", wsMessage.Model, err)
			return
		}

		modelOpts := new(llm.GGUFOptions)

		modelOpts.Model = model.Options.Model
		modelOpts.Prompt = fullPrompt
		modelOpts.ResponseDelimiter = "###RESPONSE"

		modelOpts.Temp = 0.7
		modelOpts.RepeatPenalty = 1.1

		if err := llm.MakeCompletionWebSocket(*c, chatTurn, modelOpts, config.DataPath); err != nil {
			pterm.PrintOnError(err)
			// Store the chat in the database
			chat := new(Chat)
			chat.Prompt = cpt.Messages[0].Content

			chat.Response = err.Error()
			chat.ModelName = wsMessage.Model

			pterm.Warning.Print("Storing chat in database...")
			if _, err := CreateChat(sqliteDB.db, fullPrompt, chat.Response, chat.ModelName); err != nil {
				pterm.Error.Println("Error storing chat in database:", err)
				return
			}

			// Increment the chat turn counter
			chatTurn = chatTurn + 1

			return
		}
	}))

	app.Get("/wsoai", websocket.New(func(c *websocket.Conn) {
		apiKey := config.OAIKey

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
			pterm.PrintOnError(err)
			return
		}

		// Extract the chat_message value
		chatMessage := wsMessage.ChatMessage

		// Process the message (existing logic)
		cpt := llm.GetSystemTemplate(chatMessage)

		// Sends the prompt to the AI assistant for a response
		if err := openai.StreamCompletionToWebSocket(c, chatTurn, "gpt-4-1106-preview", cpt.Messages, 0.7, apiKey); err != nil {
			pterm.PrintOnError(err)
			// Store the chat in the database
			chat := new(Chat)
			chat.Prompt = cpt.Messages[0].Content
			chat.Response = err.Error()
			chat.ModelName = wsMessage.Model

			pterm.Warning.Print("Storing chat in database...")
			if _, err := CreateChat(sqliteDB.db, chatMessage, chat.Response, chat.ModelName); err != nil {
				pterm.Error.Println("Error storing chat in database:", err)
				return
			}

			// Increment the chat turn counter
			chatTurn = chatTurn + 1

			return
		}
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
