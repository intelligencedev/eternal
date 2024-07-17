package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"eternal/pkg/documents"
	"eternal/pkg/embeddings"
	"eternal/pkg/llm"
	"eternal/pkg/llm/anthropic"
	"eternal/pkg/llm/google"
	"eternal/pkg/llm/openai"
	"eternal/pkg/vecstore"
	"eternal/pkg/web"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	index "github.com/blevesearch/bleve_index_api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/semaphore"
)

const (
// maxConcurrentJobs = 4 // Adjust this based on your system's capabilities
)

var maxConcurrentJobs = runtime.NumCPU()

type ChatTurnMessage struct {
	ID       string `json:"id"`
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
	Model    string `json:"model"`
}

// handleChatSubmit handles the submission of chat messages.
func handleChatSubmit(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPrompt := c.FormValue("userprompt")
		var wsroute string
		var model ModelParams

		if config.Tools.Team.Enabled {
			pterm.Info.Println("Team: ", currentProject.Team)

			// Retrieve the model parameters from the database.
			err := sqliteDB.First(currentProject.Team.Assistants[0].Name, &model)
			if err != nil {
				log.Errorf("Error getting model %s: %v", currentProject.Team.Assistants[0].Name, err)
				return c.Status(500).SendString("No model selected")
			}

		} else {
			selectedModels, err := GetSelectedModels(sqliteDB.db)
			if err != nil {
				log.Errorf("Error getting selected models: %v", err)
				return c.Status(500).SendString("Server Error")
			}

			if len(selectedModels) == 0 {
				log.Errorf("No models selected")
				errHtml := "<div class='alert alert-danger alert-dismissible fade show' role='alert'>No model selected <button type='button' class='btn-close' data-bs-dismiss='alert' aria-label='Close'></button></div>"
				return c.SendString(errHtml)
			}

			// Retrieve the model parameters from the database.
			err = sqliteDB.First(selectedModels[0].ModelName, &model)
			if err != nil {
				log.Errorf("Error getting model %s: %v", selectedModels[0].ModelName, err)
				return err
			}
		}

		wsroute = fmt.Sprintf("ws://%s:%s/ws", config.ServiceHosts["llm"]["llm_host_1"].Host, config.ServiceHosts["llm"]["llm_host_1"].Port)

		turnID := IncrementTurn()
		chatTurn = int(turnID)

		return c.Render("templates/chat", fiber.Map{
			"username":  config.CurrentUser,
			"message":   userPrompt,
			"assistant": config.AssistantName,
			"model":     model.Name,
			"turnID":    turnID,
			"wsRoute":   wsroute,
			"hosts":     config.ServiceHosts["llm"],
		})
	}
}

// TEST FUNCTIONS TO STOP STREAMING CHAT RESPONSE
func handleStopStreaming() fiber.Handler {
	return func(c *fiber.Ctx) error {
		turnID := c.Params("turnID")
		pterm.Info.Println("Stopping streaming for turnID:", turnID)
		// Logic to stop the streaming process for the given turnID
		// This could involve stopping the CLI process or closing the WebSocket connection
		return c.SendStatus(fiber.StatusOK)
	}
}

func stopCLIProcess(turnID string) error {
	// Logic to find and stop the CLI process associated with the given turnID
	// This could involve maintaining a map of turnIDs to process IDs
	return nil
}

func closeWebSocketConnection(turnID string) error {
	// Logic to find and close the WebSocket connection associated with the given turnID
	// This could involve maintaining a map of turnIDs to WebSocket connections
	return nil
}

// END TEST FUNCTIONS

// handleGetChats retrieves and returns all chat records.
func handleGetChats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		chats, err := GetChats(sqliteDB.db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chats"})
		}
		return c.Status(fiber.StatusOK).JSON(chats)
	}
}

// handleGetChatByID retrieves and returns a chat record by its ID.
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

// handleUpdateChat updates a chat record by its ID.
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

// handleDeleteChat handles the deletion of a chat by its ID.
func handleDeleteChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the chat ID from the request parameters.
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			// Return a bad request status if the ID is invalid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		// Attempt to delete the chat from the database.
		err = DeleteChat(sqliteDB.db, id)
		if err != nil {
			// Return an internal server error status if the deletion fails.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete chat"})
		}
		// Return a no content status if the deletion is successful.
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// handleListProjects retrieves and returns a list of projects from the database.
func handleListProjects() fiber.Handler {
	return func(c *fiber.Ctx) error {
		projects, err := sqliteDB.ListProjects()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get projects"})
		}
		return c.Status(fiber.StatusOK).JSON(projects)
	}
}

// handleUpload handles file uploads and saves them to the specified directory.
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
			log.Infof("Uploaded file %s to %s", file.Filename, filename)

			// If the file is a pdf, extract the text content and print it as Markdown.
			if strings.HasSuffix(file.Filename, ".pdf") {
				pdfDoc, err := documents.GetPdfContents(filename)
				if err != nil {
					pterm.Error.Println(err)
				}

				err = searchIndex.Index(file.Filename, pdfDoc)
				if err != nil {
					log.Errorf("Error storing chat message in Bleve: %v", err)
				}

				return c.JSON(fiber.Map{"file": file.Filename, "content": pdfDoc})
			}
		}

		// return the file path of all the documents uploaded
		return c.JSON(fiber.Map{"files": files})
	}
}

// handleDPSearch handles search requests using DuckDuckGo.
func handleDPSearch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the search query from the request.
		query := c.Query("q")
		res := web.SearchDDG(query)

		// Return an internal server error status if no results are found.
		if len(res) == 0 {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving search results")
		}

		// Return the search results as a JSON response.
		urls := res
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"urls": urls})
	}
}

// handleSSEUpdates handles Server-Sent Events (SSE) for updates.
func handleSSEUpdates() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set the necessary headers for SSE.
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		// Write updates to the response stream.
		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for {
				// Get the current download progress.
				progress := llm.GetDownloadProgress("sse-progress")
				msg := fmt.Sprintf("data: <div class='progress specific-h-25 m-4' role='progressbar' aria-label='download' aria-valuenow='%s' aria-valuemin='0' aria-valuemax='100'><div class='progress-bar progress-bar-striped progress-bar-animated' style='width: %s;'></div></div><div class='text-center fs-6'>Please refresh this page when the download completes.</br> Downloading...%s</div>\n\n", progress, progress, progress)

				// Write the progress message to the stream.
				if _, err := w.WriteString(msg); err != nil {
					pterm.Printf("Error writing to stream: %v", err)
					break
				}
				// Flush the writer to ensure the message is sent.
				if err := w.Flush(); err != nil {
					pterm.Printf("Error flushing writer: %v", err)
					break
				}

				// Sleep for 2 seconds before sending the next update.
				time.Sleep(2 * time.Second)
			}
		}))

		return nil
	}
}

// handleWebSocket handles WebSocket connections for general use.
func handleWebSocket(config *AppConfig) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		handleWebSocketConnection(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			return nil
		})
	}
}

func handleWebSocketConnection(c *websocket.Conn, config *AppConfig, processMessage func(WebSocketMessage, string) error) {
	var responseBuffer bytes.Buffer
	var wsMessage WebSocketMessage
	var err error

	// Generate uid for the chat turn

	// Read and unmarshal the initial WebSocket message
	wsMessage, err = readAndUnmarshalMessage(c)
	if err != nil {
		log.Errorf("Error reading or unmarshalling message: %v", err)
		return
	}

	chatMessage := wsMessage.ChatMessage

	// Only perform the tool workflow if any of the tools are enabled
	if config.Tools.ImgGen.Enabled || config.Tools.Memory.Enabled || config.Tools.WebGet.Enabled || config.Tools.WebSearch.Enabled {
		_, endTurn := performToolWorkflow(c, config, chatMessage)
		if endTurn {
			return
		}
	}

	if config.Tools.Team.Enabled {
		// Loop through each assistant in the group and process the chat message
		for _, assistant := range currentProject.Team.Assistants {
			err = handleAssistantTurn(c, config, wsMessage, chatMessage, &responseBuffer, assistant)
			if err != nil {
				log.Errorf("Error processing model: %v", err)
				// return
			}

			pterm.Info.Println("Assistant finished turn: ", assistant.Name)
		}
	} else {
		// Set chat as the default role if no role is selected
		if config.CurrentRoleName == "" {
			config.CurrentRoleName = "chat"

			// Get the role from the config that matches the name of the assistant role
			for _, r := range config.AssistantRoles {
				if r.Name == config.CurrentRoleName {
					config.CurrentRoleInstructions = r.Instructions
				}
			}
		}

		role := Role{
			Name:         config.CurrentRoleName,
			Instructions: config.CurrentRoleInstructions,
		}

		// Get the model params from the database
		model := ModelParams{}
		err = sqliteDB.First(wsMessage.Model, &model)
		if err != nil {
			log.Errorf("Error getting model %s: %v", wsMessage.Model, err)
			return
		}

		assistant := Assistant{
			Name: wsMessage.Model,
			Role: role,
		}
		err = handleAssistantTurn(c, config, wsMessage, chatMessage, &responseBuffer, assistant)
		if err != nil {
			if err.Error() == "EOF" {
				log.Info("Response finished: %v", err)
			} else {
				log.Errorf("Error processing model: %v", err)
			}
		}
	}

	// Handle the completed chat turn
	handleChatTurnFinished(config, wsMessage, fmt.Errorf("%s", responseBuffer.String()))
}

func handleAssistantTurn(c *websocket.Conn, config *AppConfig, wsMessage WebSocketMessage, chatMessage string, responseBuffer *bytes.Buffer, assistant Assistant) error {
	var model ModelParams
	err := sqliteDB.First(wsMessage.Model, &model)
	if err != nil {
		return fmt.Errorf("error getting model %s: %v", assistant.Name, err)
	}

	role := assistant.Role.Name

	// get the role from the config that matches the name of the assistant role
	for _, r := range config.AssistantRoles {
		if r.Name == role {
			config.CurrentRoleInstructions = r.Instructions
		}
	}

	promptTemplate := model.Options.Prompt
	// fullInstructions := fmt.Sprintf("%s\n\n%s", config.CurrentRoleInstructions, chatMessage)
	fullInstructions := fmt.Sprintf("Query: %s\n\nPrevious Response: %s\n\n%s", chatMessage, responseBuffer.String(), config.CurrentRoleInstructions)
	fullPrompt := strings.ReplaceAll(promptTemplate, "{prompt}", fullInstructions)
	fullPrompt = strings.ReplaceAll(fullPrompt, "{system}", "You are a helpful AI assistant.")

	modelOpts := &llm.GGUFOptions{
		NGPULayers:    config.ServiceHosts["llm"]["llm_host_1"].GgufGPULayers,
		Model:         model.Options.Model,
		Prompt:        fullPrompt,
		CtxSize:       model.Options.CtxSize,
		Temp:          model.Options.Temp,
		RepeatPenalty: model.Options.RepeatPenalty,
		TopP:          model.Options.TopP,
		TopK:          model.Options.TopK,
	}

	// Insert an alert with the name and role of the assistant into the response buffer
	responseBuffer.WriteString(fmt.Sprintf("\n### %s - %s\n", assistant.Name, assistant.Role.Name))

	// invoke the correct handler based on the model name
	if strings.HasPrefix(model.Name, "openai-") {
		// Get the system template for the chat message.
		cpt := llm.GetSystemTemplate(fullPrompt)
		return openai.StreamCompletionToWebSocket(*c, chatTurn, "gpt-4o", cpt.Messages, 0.3, config.OAIKey, responseBuffer)
	} else if strings.HasPrefix(model.Name, "google-") {
		apiKey := config.GoogleKey
		return google.StreamGeminiResponseToWebSocket(*c, chatTurn, fullPrompt, apiKey, responseBuffer)
	} else if strings.HasPrefix(model.Name, "anthropic-") {
		apiKey := config.AnthropicKey

		// Prepare the messages for the completion request.
		messages := []anthropic.Message{
			{Role: "user", Content: fullPrompt},
		}

		// Stream the completion response from Anthropic to the WebSocket.
		return anthropic.StreamCompletionToWebSocket(*c, chatTurn, "claude-3-5-sonnet-20240620", messages, 0.3, apiKey, responseBuffer)
	} else {
		return llm.MakeCompletionWebSocket(*c, chatTurn, modelOpts, config.DataPath, responseBuffer)
	}
}

// readAndUnmarshalMessage reads and unmarshals a WebSocket message.
func readAndUnmarshalMessage(c *websocket.Conn) (WebSocketMessage, error) {
	// Read the message from the WebSocket.
	_, messageBytes, err := c.ReadMessage()
	if err != nil {
		return WebSocketMessage{}, err
	}

	// Unmarshal the JSON message.
	var wsMessage WebSocketMessage
	err = json.Unmarshal(messageBytes, &wsMessage)
	if err != nil {
		return WebSocketMessage{}, err
	}

	return wsMessage, nil
}

// handleError handles errors that occur during message processing.
func handleChatTurnFinished(config *AppConfig, message WebSocketMessage, err error) {
	// Store the chat turn in the sqlite db.
	if _, err := CreateChat(sqliteDB.db, message.ChatMessage, err.Error(), message.Model); err != nil {
		pterm.Error.Println("Error storing chat in database:", err)
		return
	}

	if config.Tools.Memory.Enabled {

		// Get the timestamp for the chat message in human-readable format.
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		memHeader := fmt.Sprintf("Previous chat - %s", timestamp)

		// Two examples of how to store chat messages in the Bleve index.
		// 1. Split the text and store each chunk in the index.
		// 2. Store the entire chat message in the index.
		// Split the chat message into chunks 500 characters long with a 200 character overlap.
		chunks := documents.SplitTextByCount(err.Error(), 500)

		// Prepend the header to all chunks.
		for i, chunk := range chunks {
			chunks[i] = fmt.Sprintf("%s\n%s", memHeader, chunk)
		}

		// 1. Store the chunk in Bleve.
		for _, chunk := range chunks {
			chatMessage := ChatTurnMessage{
				ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
				Prompt:   message.ChatMessage,
				Response: chunk,
				Model:    message.Model,
			}

			err = searchIndex.Index(chatMessage.ID, chatMessage)
			if err != nil {
				log.Errorf("Error storing chat message in Bleve: %v", err)
			}
		}
		// 2. Store the entire chat message in Bleve.
		// chatMessage := ChatTurnMessage{
		// 	ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		// 	Prompt:   message.ChatMessage,
		// 	Response: err.Error(),
		// 	Model:    message.Model,
		// }

		// err = searchIndex.Index(chatMessage.ID, chatMessage)
		// if err != nil {
		// 	log.Errorf("Error storing chat message in Bleve: %v", err)
		// }

		log.Errorf("Chat turn finished: %v", err)
	}
}

// handleChatMemory retrieves and returns chat memory.
func handleChatMemory(config *AppConfig, chatMessage string, docType string) (string, error) {
	topN := config.Tools.Memory.TopN

	chatMessage = strings.TrimSpace(chatMessage)
	chatMessage = strings.ReplaceAll(chatMessage, "\n", "")

	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		return "", err
	}
	chatMessage = reg.ReplaceAllString(chatMessage, " ")

	query := bleve.NewQueryStringQuery(chatMessage)
	searchRequest := bleve.NewSearchRequestOptions(query, topN, 0, false)
	searchResults, err := searchIndex.Search(searchRequest)
	if err != nil {
		log.Errorf("Error searching index: %v", err)
	}

	modelPath := fmt.Sprintf("%s/models", config.DataPath)

	pterm.Info.Println("Splitting prompt into chunks...")
	chunks := documents.SplitTextByCount(chatMessage, 500)

	pterm.Info.Println("Calculating total weight...")
	totalWeight := 0.0
	for i := 0; i < len(chunks); i++ {
		if i == 0 || i == len(chunks)-1 {
			totalWeight += 1
		} else {
			totalWeight += 0.1
		}
	}

	pterm.Info.Println("Calculating weighted average...")
	sem := semaphore.NewWeighted(int64(maxConcurrentJobs))
	ctx := context.Background()

	promptVec, err := calculateWeightedAverage(chunks, modelPath, config.EmbeddingModel, totalWeight, sem, ctx)
	if err != nil {
		return "", err
	}

	pterm.Info.Println("Prompt vector:", promptVec)

	resDocument := processSearchResults(searchResults, promptVec, modelPath, config.EmbeddingModel, sem, ctx)

	return resDocument, nil
}

func calculateWeightedAverage(chunks []string, modelPath, embeddingModel string, totalWeight float64, sem *semaphore.Weighted, ctx context.Context) ([]float64, error) {
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks to process")
	}

	firstVec, err := embeddings.GenerateEmbeddingsForChat(chunks[0], modelPath, embeddingModel)
	if err != nil {
		return nil, fmt.Errorf("error generating embedding for first chunk: %v", err)
	}

	promptVec := make([]float64, len(firstVec))
	var wg sync.WaitGroup
	chunkChan := make(chan struct {
		index int
		vec   []float64
	}, len(chunks))

	for i, chunk := range chunks {
		wg.Add(1)
		go func(i int, chunk string) {
			defer wg.Done()
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Errorf("Failed to acquire semaphore: %v", err)
				return
			}
			defer sem.Release(1)

			vec, err := embeddings.GenerateEmbeddingsForChat(chunk, modelPath, embeddingModel)
			if err != nil {
				log.Errorf("Error generating embedding for chat: %v", err)
				return
			}
			chunkChan <- struct {
				index int
				vec   []float64
			}{i, vec}
		}(i, chunk)
	}

	go func() {
		wg.Wait()
		close(chunkChan)
	}()

	for chunk := range chunkChan {
		weight := 0.1
		if chunk.index == 0 || chunk.index == len(chunks)-1 {
			weight = 1
		}
		for j := 0; j < len(chunk.vec); j++ {
			promptVec[j] += chunk.vec[j] * weight / totalWeight
		}
	}

	return promptVec, nil
}

func processSearchResults(searchResults *bleve.SearchResult, promptVec []float64, modelPath, embeddingModel string, sem *semaphore.Weighted, ctx context.Context) string {
	var resDocument string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, hit := range searchResults.Hits {
		wg.Add(1)
		go func(hit *search.DocumentMatch) {
			defer wg.Done()
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Errorf("Failed to acquire semaphore: %v", err)
				return
			}
			defer sem.Release(1)

			doc, err := searchIndex.Document(hit.ID)
			if err != nil {
				log.Errorf("Error retrieving document: %v", err)
				return
			}

			doc.VisitFields(func(field index.Field) {
				if field.Name() == "prompt" || field.Name() == "response" {
					chunk := fmt.Sprintf("%s", field.Value())
					pterm.Info.Println("Chunk:", chunk)

					vec, err := embeddings.GenerateEmbeddingsForChat(chunk, modelPath, embeddingModel)
					if err != nil {
						log.Errorf("Error generating embedding for chat: %v", err)
						return
					}

					similarity := vecstore.CosineSimilarity(promptVec, vec)
					pterm.Info.Println("Similarity:", similarity)

					threshold := 0.5
					if field.Name() == "prompt" {
						threshold = 0.6
					}

					if similarity > threshold {
						mu.Lock()
						resDocument = fmt.Sprintf("%s\n%s", resDocument, chunk)
						mu.Unlock()
					}
				}
			})
		}(hit)
	}

	wg.Wait()
	return resDocument
}

// func handleChatMemory(config *AppConfig, chatMessage string, docType string) (string, error) {
// 	// var document string

// 	topN := config.Tools.Memory.TopN

// 	// Remove any leading or trailing whitespace from the chat message.
// 	chatMessage = strings.TrimSpace(chatMessage)

// 	// Remove line breaks from the chat message.
// 	chatMessage = strings.ReplaceAll(chatMessage, "\n", "")

// 	// Create a regular expression to remove special characters (keeping alphanumeric and spaces)
// 	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
// 	if err != nil {
// 		return "", err
// 	}

// 	// Apply the regex to sanitize the chat message
// 	chatMessage = reg.ReplaceAllString(chatMessage, " ")

// 	// Print the chat message
// 	// pterm.Info.Println("Chat message:")
// 	// pterm.Info.Println(chatMessage)

// 	// Create a search query
// 	query := bleve.NewQueryStringQuery(chatMessage)

// 	// Create a search request with the query and limit the results
// 	searchRequest := bleve.NewSearchRequestOptions(query, topN, 0, false)

// 	// Execute the search
// 	searchResults, err := searchIndex.Search(searchRequest)
// 	if err != nil {
// 		log.Errorf("Error searching index: %v", err)
// 	}

// 	modelPath := fmt.Sprintf("%s/models", config.DataPath)

// 	// Chunk the chat message and generate embeddings
// 	pterm.Info.Println("Splitting prompt into chunks...")
// 	chunks := documents.SplitTextByCount(chatMessage, 500)

// 	// Use Weighted averaging: Assign weights to chunks based on importance (e.g., position in document) and compute a weighted average.
// 	// The first and last chunks are assigned a weight of 1, and the rest are assigned a weight of 0.1.

// 	// Calculate the total weight
// 	pterm.Info.Println("Calculating total weight...")
// 	totalWeight := 0.0
// 	for i := 0; i < len(chunks); i++ {
// 		if i == 0 || i == len(chunks)-1 {
// 			totalWeight += 1
// 		} else {
// 			totalWeight += 0.1
// 		}
// 	}

// 	// Calculate the weighted average
// 	pterm.Info.Println("Calculating weighted average...")
// 	// Initialize promptVec with the correct length
// 	var promptVec []float64
// 	if len(chunks) > 0 {
// 		// Generate embeddings for the first chunk to get the vector length
// 		firstVec, err := embeddings.GenerateEmbeddingsForChat(chunks[0], modelPath, config.EmbeddingModel)
// 		if err != nil {
// 			log.Errorf("Error generating embedding for first chunk: %v", err)
// 			return "", err
// 		}
// 		promptVec = make([]float64, len(firstVec))
// 	}

// 	for i, chunk := range chunks {
// 		// Generate embeddings for the chunk
// 		vec, err := embeddings.GenerateEmbeddingsForChat(chunk, modelPath, config.EmbeddingModel)
// 		if err != nil {
// 			log.Errorf("Error generating embedding for chat: %v", err)
// 			continue
// 		}

// 		// Calculate the weight for the chunk
// 		var weight float64
// 		if i == 0 || i == len(chunks)-1 {
// 			weight = 1
// 		} else {
// 			weight = 0.1
// 		}

// 		// Calculate the weighted average
// 		for j := 0; j < len(vec); j++ {
// 			promptVec[j] += vec[j] * weight / totalWeight
// 		}
// 	}

// 	// Print the prompt vector
// 	pterm.Info.Println("Prompt vector:")
// 	pterm.Info.Println(promptVec)

// 	// Will contain the chunks with the highest similarity
// 	var resDocument string

// 	// Print the search results
// 	for _, hit := range searchResults.Hits {
// 		doc, err := searchIndex.Document(hit.ID)
// 		if err != nil {
// 			log.Errorf("Error retrieving document: %v", err)
// 			continue
// 		}

// 		doc.VisitFields(func(field index.Field) {
// 			// fmt.Printf("%s: %s\n", field.Name(), field.Value())

// 			// Append the response field to the document and store it for later use
// 			if field.Name() == "prompt" {
// 				chunk := fmt.Sprintf("%s", field.Value())

// 				// Print the chunk
// 				pterm.Info.Println("Chunk:")

// 				// Generate embeddings and calculate similarity
// 				vec, err := embeddings.GenerateEmbeddingsForChat(chunk, modelPath, config.EmbeddingModel)
// 				if err != nil {
// 					log.Errorf("Error generating embedding for chat: %v", err)
// 				}

// 				// Compute similarity
// 				similarity := vecstore.CosineSimilarity(promptVec, vec)

// 				// Print the similarity score
// 				pterm.Info.Println("Similarity:", similarity)

// 				// If the similarity is above a certain threshold, add to our document
// 				if similarity > 0.6 {
// 					resDocument = fmt.Sprintf("%s\n%s", resDocument, chunk)
// 				}
// 			}
// 			if field.Name() == "response" {
// 				chunk := fmt.Sprintf("%s", field.Value())

// 				// Print the chunk
// 				pterm.Info.Println("Chunk:")

// 				// Generate embeddings and calculate similarity
// 				vec, err := embeddings.GenerateEmbeddingsForChat(chunk, modelPath, config.EmbeddingModel)
// 				if err != nil {
// 					log.Errorf("Error generating embedding for chat: %v", err)
// 				}

// 				// Compute similarity
// 				similarity := vecstore.CosineSimilarity(promptVec, vec)

// 				// Print the similarity score
// 				pterm.Info.Println("Similarity:", similarity)

// 				// If the similarity is above a certain threshold, add to our document
// 				if similarity > 0.5 {
// 					resDocument = fmt.Sprintf("%s\n%s", resDocument, chunk)
// 				}
// 			}
// 		})
// 	}

// 	return resDocument, nil
// }

// storePromptInBleve stores the user's prompt in the Bleve index.
// func storeChat(text string, doctype string) error {
// 	// Generate a unique ID for the document
// 	docID := fmt.Sprintf("prompt_%d", time.Now().UnixNano())

// 	// Create a document to store in the index
// 	doc := struct {
// 		ID      string `json:"id"`
// 		Content string `json:"content"`
// 		Type    string `json:"type"`
// 	}{
// 		ID:      docID,
// 		Content: text,
// 		Type:    doctype,
// 	}

// 	// Index the document
// 	err := searchIndex.Index(docID, doc)
// 	if err != nil {
// 		return fmt.Errorf("error storing prompt in Bleve: %v", err)
// 	}

// 	return nil
// }

// handleTextSplitAndIndex handles the splitting and indexing of text.
func handleTextSplitAndIndex(inputTags string, inputText string, chunkSize int, modelName string) error {
	// Split the input text into chunks.
	chunks := documents.SplitTextByCount(inputText, chunkSize)

	// Prepend the input tags to each chunk.
	for i, chunk := range chunks {
		chunks[i] = fmt.Sprintf("TAGS: [%s]\n%s", inputTags, chunk)
	}

	var wg sync.WaitGroup

	for _, chunk := range chunks {
		wg.Add(1)

		go func(c string) {
			defer wg.Done()

			docID := fmt.Sprintf("%d", time.Now().UnixNano())
			doc := ChatTurnMessage{
				ID:       docID,
				Prompt:   inputText,
				Response: c,
				Model:    modelName,
			}

			if err := searchIndex.Index(docID, doc); err != nil {
				log.Errorf("Error indexing chunk in Bleve: %v", err)
			}
		}(chunk)
	}

	wg.Wait()

	return nil
}

// We need to recode the function so it does not store embeddings at all
// We can compute the cosine similarity for each chunk Bleve returns to the Prompt
// Then we can return the top N chunks with the highest similarity

// searchSimilarEmbeddings searches for similar embeddings in the database.
// func searchSimilarEmbeddings(config *AppConfig, prompt string, topN int) []vecstore.Embedding {
// 	// db := vecstore.NewEmbeddingDB()
// 	// dbPath := fmt.Sprintf("%s/embeddings.json", config.DataPath)
// 	// embeddings, err := db.LoadEmbeddings(dbPath)
// 	// if err != nil {
// 	// 	fmt.Println("Error loading embeddings:", err)
// 	// 	return nil
// 	// }

// 	modelsDir := fmt.Sprintf("%s/models", config.DataPath)

// 	tasksConfig := &tasks.Config{
// 		ModelsDir:        modelsDir,
// 		ModelName:        config.EmbeddingModel,
// 		DownloadPolicy:   tasks.DownloadMissing,
// 		ConversionPolicy: tasks.ConvertMissing,
// 	}

// 	model, err := tasks.Load[textencoding.Interface](tasksConfig)
// 	if err != nil {

// 		pterm.Error.Println("Error loading model while searching most similar chunks: ", err)
// 		return nil
// 	}

// 	var vec []float64
// 	result, err := model.Encode(context.Background(), prompt, int(bert.MeanPooling))
// 	if err != nil {
// 		fmt.Println("Error encoding text:", err)
// 		return nil
// 	}
// 	vec = result.Vector.Data().F64()[:128]

// 	// Compute similarity
// 	similarity := vecstore.CosineSimilarity(vec, vec)

// 	embeddingForPrompt := vecstore.Embedding{
// 		Word:       prompt,
// 		Vector:     vec,
// 		Similarity: 0.0,
// 	}

// 	// Retrieve the top N similar embeddings
// 	topEmbeddings := vecstore.FindTopNSimilarEmbeddings(embeddingForPrompt, embeddings, topN)
// 	if len(topEmbeddings) == 0 {
// 		fmt.Println("Error finding similar embeddings.")
// 		return nil
// 	}

// 	return topEmbeddings
// }

// ToolState represents the state of a tool.
type ToolState struct {
	Tool    string `json:"tool"`
	Enabled bool   `json:"enabled"`
}
