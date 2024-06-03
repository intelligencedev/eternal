package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"eternal/pkg/embeddings"
	"eternal/pkg/hfutils"
	"eternal/pkg/llm"
	"eternal/pkg/llm/anthropic"
	"eternal/pkg/llm/google"
	"eternal/pkg/llm/openai"
	"eternal/pkg/sd"
	"eternal/pkg/web"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	index "github.com/blevesearch/bleve_index_api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

var assistantRole = "You are a helpful AI assistant that responds in well-structured markdown format. Do not repeat your instructions. Do not deviate from the topic."

type ChatTurnMessage struct {
	ID       string `json:"id"`
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
	Model    string `json:"model"`
}

func handleUpload(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pterm.Warning.Println("Uploads route hit")

		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		files := form.File["file"]
		for _, file := range files {
			filename := filepath.Join(config.DataPath, "web", "uploads", file.Filename)
			pterm.Warning.Printf("Uploading file: %s\n", filename)
			if err := c.SaveFile(file, filename); err != nil {
				return err
			}
			log.Infof("Uploaded file: %s", file.Filename)
		}

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("%d files uploaded successfully", len(files)),
		})
	}
}

func handleToolToggle(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		toolName := c.Params("toolName")

		switch toolName {
		case "websearch":
			config.Tools.WebSearch.Enabled = !config.Tools.WebSearch.Enabled
		case "webget":
			config.Tools.WebGet.Enabled = !config.Tools.WebGet.Enabled
		case "imggen":
			config.Tools.ImgGen.Enabled = true
		default:
			return c.Status(fiber.StatusNotFound).SendString("Tool not found")
		}

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Tool %s is now %t", toolName, config.Tools.ImgGen.Enabled)})
	}
}

func handleOpenAIModels(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		client := openai.NewClient(config.OAIKey)
		modelsResponse, err := openai.GetModels(client)

		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

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

func handleModelData() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

func handleModelDownloadUpdate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		modelName := c.Params("modelName")
		var payload struct {
			Downloaded bool `json:"downloaded"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		err := sqliteDB.UpdateDownloadedByName(modelName, payload.Downloaded)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update model: %v", err)})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Model 'Downloaded' status updated successfully",
		})
	}
}

func handleModelCards(modelParams []ModelParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := sqliteDB.Find(&modelParams)

		if err != nil {
			log.Errorf("Database error: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		return c.Render("templates/model", fiber.Map{"models": modelParams})
	}
}

func handleModelSelect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var selection SelectedModels

		if err := c.BodyParser(&selection); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Bad request")
		}

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
	}
}

func handleSelectedModels() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

func handleModelDownload(config *AppConfig) fiber.Handler {
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

		var partialDownload bool
		if info, err := os.Stat(modelPath); err == nil {
			if info.Size() > 0 {
				expectedSize, err := llm.GetExpectedFileSize(downloadURL)
				if err != nil {
					log.Errorf("Error getting expected file size: %v", err)
				}
				partialDownload = info.Size() < expectedSize
			}
		}

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
				err = sqliteDB.UpdateDownloadedByName(modelName, true)
				if err != nil {
					log.Errorf("Failed to update model downloaded state: %v", err)
				}
			}
		}()

		progressErr := fmt.Sprintf("<div class='w-100' id='progress-download-%s' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message' hx-trigger='load'></div>", modelName)

		return c.SendString(progressErr)
	}
}

func handleImgModelDownload(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		config.Tools.ImgGen.Enabled = true

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

		modelRoot := fmt.Sprintf("%s/models/%s", config.DataPath, modelName)
		modelPath := fmt.Sprintf("%s/models/%s/%s", config.DataPath, modelName, modelFileName)
		tmpPath := fmt.Sprintf("%s/tmp", config.DataPath)

		if _, err := os.Stat(modelRoot); os.IsNotExist(err) {
			if err := os.MkdirAll(modelRoot, 0755); err != nil {
				log.Errorf("Error creating model directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
			if err := os.MkdirAll(tmpPath, 0755); err != nil {
				log.Errorf("Error creating tmp directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		if _, err := os.Stat(modelPath); err != nil {
			dm := hfutils.ConcurrentDownloadManager{
				FileName:    modelFileName,
				URL:         downloadURL,
				Destination: modelPath,
				NumParts:    1,
				TempDir:     tmpPath,
			}

			go dm.PrintProgress()

			if err := dm.Download(); err != nil {
				fmt.Println("Download failed:", err)
			} else {
				fmt.Println("Download successful!")
			}
		}

		vaeName := "sdxl_vae.safetensors"
		vaeURL := "https://huggingface.co/madebyollin/sdxl-vae-fp16-fix/blob/main/sdxl_vae.safetensors"
		vaePath := fmt.Sprintf("%s/models/%s/%s", config.DataPath, modelName, vaeName)

		if _, err := os.Stat(modelRoot); os.IsNotExist(err) {
			if err := os.MkdirAll(modelRoot, 0755); err != nil {
				log.Errorf("Error creating model directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		if _, err := os.Stat(vaePath); os.IsNotExist(err) {
			go func() {
				response, err := http.Get(vaeURL)
				if err != nil {
					pterm.Error.Printf("Failed to download file: %v", err)
					return
				}
				defer response.Body.Close()

				file, err := os.Create(vaePath)
				if err != nil {
					pterm.Error.Printf("Failed to create file: %v", err)
					return
				}
				defer file.Close()

				_, err = io.Copy(file, response.Body)
				if err != nil {
					pterm.Error.Printf("Failed to write to file: %v", err)
					return
				}

				pterm.Info.Printf("Downloaded file: %s", vaeName)
			}()
		}

		progressErr := "<div name='sse-messages' class='w-100' id='sse-messages' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message'></div>"

		return c.SendString(progressErr)
	}
}

func handleRoleSelection(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName := c.Params("name")
		var foundRole *struct {
			Name         string `yaml:"name"`
			Instructions string `yaml:"instructions"`
		}

		for i := range config.AssistantRoles {
			if config.AssistantRoles[i].Name == roleName {
				foundRole = &config.AssistantRoles[i]
				break
			}
		}

		if foundRole == nil {
			pterm.Warning.Printf("Role %s not found. Defaulting to 'chat'.\n", roleName)
			for i := range config.AssistantRoles {
				if config.AssistantRoles[i].Name == "chat" {
					foundRole = &config.AssistantRoles[i]
					break
				}
			}
		}

		if foundRole == nil && len(config.AssistantRoles) > 0 {
			foundRole = &config.AssistantRoles[0]
		}

		if foundRole != nil {
			assistantRole = foundRole.Instructions
			pterm.Info.Printf("Role set to: %s\n", foundRole.Name)
			pterm.Info.Println(foundRole.Instructions)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("Role set to %s", foundRole.Name),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No roles configured",
		})
	}
}

func handleChatSubmit(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPrompt := c.FormValue("userprompt")
		var wsroute string

		selectedModels, err := GetSelectedModels(sqliteDB.db)
		if err != nil {
			log.Errorf("Error getting selected models: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		if len(selectedModels) > 0 {
			firstModelName := selectedModels[0].ModelName

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
			return c.JSON(fiber.Map{"error": "No models selected"})
		}

		turnID := IncrementTurn()

		return c.Render("templates/chat", fiber.Map{
			"username":  config.CurrentUser,
			"message":   userPrompt,
			"assistant": config.AssistantName,
			"model":     selectedModels[0].ModelName,
			"turnID":    turnID,
			"wsRoute":   wsroute,
			"hosts":     config.ServiceHosts["llm"],
		})
	}
}

func handleGetChats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		chats, err := GetChats(sqliteDB.db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chats"})
		}
		return c.Status(fiber.StatusOK).JSON(chats)
	}
}

func handleGetChatByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		chat, err := GetChatByID(sqliteDB.db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chat"})
		}
		return c.Status(fiber.StatusOK).JSON(chat)
	}
}

func handleUpdateChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

func handleDeleteChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		err = DeleteChat(sqliteDB.db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete chat"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func handleDPSearch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("q")
		res := web.SearchDDG(query)

		if len(res) == 0 {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving search results")
		}

		urls := res
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"urls": urls})
	}
}

func handleSSEUpdates() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
}

func handleWebSocket(config *AppConfig) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		handleWebSocketConnection(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			var model ModelParams
			err := sqliteDB.First(wsMessage.Model, &model)
			if err != nil {
				log.Errorf("Error getting model %s: %v", wsMessage.Model, err)
				return err
			}

			promptTemplate := model.Options.Prompt
			fullPrompt := strings.ReplaceAll(promptTemplate, "{prompt}", chatMessage)
			fullPrompt = strings.ReplaceAll(fullPrompt, "{system}", assistantRole)

			modelOpts := &llm.GGUFOptions{
				NGPULayers:    config.ServiceHosts["llm"]["llm_host_1"].GgufGPULayers,
				Model:         model.Options.Model,
				Prompt:        fullPrompt,
				CtxSize:       model.Options.CtxSize,
				Temp:          0.1,
				RepeatPenalty: 1.1,
				TopP:          1.0,
				TopK:          1.0,
			}

			return llm.MakeCompletionWebSocket(*c, chatTurn, modelOpts, config.DataPath)
		})
	}
}

func handleOpenAIWebSocket(config *AppConfig) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		handleWebSocketConnection(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			cpt := llm.GetSystemTemplate(chatMessage)
			return openai.StreamCompletionToWebSocket(c, chatTurn, "gpt-4o", cpt.Messages, 0.3, config.OAIKey)
		})
	}
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
}

func handleAnthropicWebSocket(config *AppConfig) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		apiKey := config.AnthropicKey
		handleAnthropicWS(c, apiKey, chatTurn)
	}
}

func handleGoogleWebSocket(config *AppConfig) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		apiKey := config.GoogleKey

		handleWebSocketConnection(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			return google.StreamGeminiResponseToWebSocket(c, chatTurn, chatMessage, apiKey)
		})
	}
}

func handleWebSocketConnection(c *websocket.Conn, config *AppConfig, processMessage func(WebSocketMessage, string) error) {
	for {
		wsMessage, err := readAndUnmarshalMessage(c)
		if err != nil {
			log.Errorf("Error reading or unmarshalling message: %v", err)
			return
		}

		log.Infof("Received WebSocket message: %+v", wsMessage)

		chatMessage := performToolWorkflow(c, config, wsMessage.ChatMessage)
		log.Infof("Processed chat message: %s", chatMessage)

		err = processMessage(wsMessage, chatMessage)
		if err != nil {
			handleError(c, wsMessage, err)
			return
		}

		log.Info("Message processed successfully")
	}
}

func readAndUnmarshalMessage(c *websocket.Conn) (WebSocketMessage, error) {
	_, messageBytes, err := c.ReadMessage()
	if err != nil {
		return WebSocketMessage{}, err
	}

	var wsMessage WebSocketMessage
	err = json.Unmarshal(messageBytes, &wsMessage)
	if err != nil {
		return WebSocketMessage{}, err
	}

	return wsMessage, nil
}

func handleError(message WebSocketMessage, err error) {
	log.Errorf("Error processing message: %v", err)

	// Store chat message in Bleve
	chatMessage := ChatTurnMessage{
		ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		Prompt:   message.ChatMessage,
		Response: err.Error(),
		Model:    message.Model,
	}

	err = searchIndex.Index(chatMessage.ID, chatMessage)
	if err != nil {
		log.Errorf("Error storing chat message in Bleve: %v", err)
	}

	// Increment the chat turn counter
	chatTurn++
}

func performToolWorkflow(c *websocket.Conn, config *AppConfig, chatMessage string) string {
	// Begin tool workflow. Tools will add context to the submitted message for
	// the model to use. Document is the abstraction that will hold that context.
	var document string

	if config.Tools.ImgGen.Enabled {
		pterm.Info.Println("Generating image...")
		sdParams := &sd.SDParams{Prompt: chatMessage}

		// Call the sd tool
		sd.Text2Image(config.DataPath, sdParams)

		// Return the image to the client
		timestamp := time.Now().UnixNano() // Get the current timestamp in nanoseconds
		imgElement := fmt.Sprintf("<img class='rounded-2 object-fit-scale' width='512' height='512' src='public/img/sd_out.png?%d' />", timestamp)
		formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>", fmt.Sprint(chatTurn), imgElement)
		if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
			pterm.PrintOnError(err)
			return chatMessage
		}

		// Increment the chat turn counter
		chatTurn = chatTurn + 1

		// End the tool workflow
		return chatMessage
	}

	if config.Tools.Memory.Enabled {
		topN := config.Tools.Memory.TopN // retrieve top N results. Adjust based on context size.

		// Create a search query
		query := bleve.NewQueryStringQuery(chatMessage)

		// Create a search request with the query and limit the results
		searchRequest := bleve.NewSearchRequestOptions(query, topN, 0, false)

		// Execute the search
		searchResults, err := searchIndex.Search(searchRequest)
		if err != nil {
			log.Errorf("Error searching index: %v", err)
			return chatMessage
		}

		// Print the search results
		for _, hit := range searchResults.Hits {
			doc, err := searchIndex.Document(hit.ID)
			if err != nil {
				log.Errorf("Error retrieving document: %v", err)
				continue
			}
			doc.VisitFields(func(field index.Field) {
				fmt.Printf("%s: %s\n", field.Name(), field.Value())

				// Append the response field to the document
				if field.Name() == "response" {
					document = fmt.Sprintf("%s\n%s", document, field.Value())
				}
			})
		}

		pterm.Info.Println(searchResults)
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

		var urls []string
		if config.Tools.WebSearch.Name == "ddg" {
			urls = web.SearchDDG(chatMessage)
		} else if config.Tools.WebSearch.Name == "sxng" {
			urls = web.GetSearXNGResults(config.Tools.WebSearch.Endpoint, chatMessage)
		}

		pterm.Warning.Printf("URLs to fetch: %v\n", urls)

		pagesRetrieved := 0
		for _, url := range urls {
			if pagesRetrieved >= topN {
				break
			}
			pterm.Info.Printf("Fetching URL: %s\n", url)

			page, err := web.WebGetHandler(url)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					pterm.Warning.Printf("Timeout exceeded for URL: %s\n", url)
					continue
				}
				pterm.PrintOnError(err)
			} else {
				document = fmt.Sprintf("%s\n%s", document, page)
				pagesRetrieved++
			}
		}
	}

	chatMessage = fmt.Sprintf("%s Reference the previous information if it is relevant to the next query only. Do not provide any additional information other than what is necessary to answer the next question or respond to the query. Be concise. Do not deviate from the topic of the query.\nQUERY:\n%s", document, chatMessage)

	pterm.Info.Println("Tool workflow complete")

	//pterm.Warning.Println("Chat message:\n", chatMessage)

	return chatMessage
}

func storeChat(db *gorm.DB, config *AppConfig, prompt, response, modelName string) error {
	// Generate embeddings
	pterm.Warning.Println("Generating embeddings for chat...")

	chatText := fmt.Sprintf("QUESTION: %s\n RESPONSE: %s", prompt, response)
	err := embeddings.GenerateEmbeddingForTask("chat", chatText, "txt", 500, 100, config.DataPath)
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
