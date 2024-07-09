package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"eternal/pkg/web"
	"fmt"
	"image"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/pterm/pterm"

	socket "github.com/gorilla/websocket"
)

var currentChatUid string

// handleGetTools returns a handler function that retrieves all tools
func handleRenderTools(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("templates/tools", fiber.Map{
			"Tools": config.Tools,
		})
	}
}

// performToolWorkflow performs the tool workflow on a chat message.
func performToolWorkflow(c *websocket.Conn, config *AppConfig, chatMessage string) string {

	// Begin tool workflow. Tools will add context to the submitted message for the model to use.
	var document string

	if config.Tools.ImgGen.Enabled {
		pterm.Info.Println("Generating image...")
		chatId := uuid.New().String()
		res := performImageGen(chatId, config, chatMessage)
		c.WriteMessage(socket.TextMessage, []byte(res))
		chatTurn = chatTurn + 1
		return chatMessage
	}

	if config.Tools.Memory.Enabled {
		document, _ = handleChatMemory(config, chatMessage)
	}

	if config.Tools.WebGet.Enabled {
		url := web.ExtractURLs(chatMessage)
		if len(url) > 0 {
			pterm.Info.Println("Retrieving page content...")

			document, _ = web.WebGetHandler(url[0])

			// Add the page content to the chat message.

		}
	}

	if config.Tools.WebSearch.Enabled {
		topN := config.Tools.WebSearch.TopN

		pterm.Info.Println("Searching the web...")

		var urls []string
		switch config.Tools.WebSearch.Name {
		case "ddg":
			urls = web.SearchDDG(chatMessage)
		case "sxng":
			urls = web.GetSearXNGResults(config.Tools.WebSearch.Endpoint, chatMessage)
		}

		//pterm.Warning.Printf("URLs to fetch: %v\n", urls)

		ignoredURLs, err := sqliteDB.ListURLTrackings()
		if err != nil {
			log.Errorf("Error listing URL trackings: %v", err)
		}

		// match the ignored URLs with the fetched URLs and remove them from the list
		for _, ignoredURL := range ignoredURLs {
			for i, url := range urls {
				if strings.Contains(url, ignoredURL.URL) {
					urls = append(urls[:i], urls[i+1:]...)

					pterm.Warning.Printf("Ignoring URL: %s\n", ignoredURL.URL)
				}
			}
		}

		var wg sync.WaitGroup
		urlsChan := make(chan string, len(urls))
		failedURLsChan := make(chan []string)
		pagesChan := make(chan string, topN)
		done := make(chan struct{})

		// Fetch URLs concurrently
		for _, url := range urls {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				select {
				case <-done:
					return
				default:
					pterm.Info.Printf("Fetching URL: %s\n", u)
					page, err := web.WebGetHandler(u)
					if err != nil {
						if errors.Is(err, context.DeadlineExceeded) {
							pterm.Warning.Printf("Timeout exceeded for URL: %s\n", u)

							// Add the URL to the channel to be processed later
							failedURLsChan <- []string{u}
						} else {
							log.Errorf("Error fetching URL: %v", err)

							failedURLsChan <- []string{u}
						}
						return
					}

					// Prepent the URL to the page content
					page = fmt.Sprintf("%s\n%s", u, page)

					urlsChan <- page
				}
			}(url)
		}

		// Close urlsChan when all fetches are done
		go func() {
			wg.Wait()
			close(urlsChan)
			close(failedURLsChan)
		}()

		// Collect topN pages
		go func() {
			var pagesRetrieved int
			for page := range urlsChan {
				if pagesRetrieved >= topN {
					close(done)
					break
				}
				pagesChan <- page
				pagesRetrieved++
			}
			close(pagesChan)
		}()

		// Process failed URLs
		var failedURLs []string
		for url := range failedURLsChan {
			failedURLs = append(failedURLs, url...)

			// Insert the failed URLs back into the URLTracking table
			for _, failedURL := range failedURLs {
				// Parse the top-level domain from the URL by splitting the URL by slashes and getting the second element.
				tld := strings.Split(failedURL, "/")[2]

				err := sqliteDB.CreateURLTracking(tld)
				if err != nil {
					log.Errorf("Error inserting failed URL into database: %v", err)
				}
			}
		}

		// Retreve the failed URLs from the URLTracking table
		trackedURLs, err := sqliteDB.ListURLTrackings()
		if err != nil {
			log.Errorf("Error listing URL trackings: %v", err)
		}

		// Print the failed URLs
		for _, trackedURL := range trackedURLs {
			pterm.Warning.Printf("New failed URL: %s\n", trackedURL.URL)
		}

		// Process pages
		var document string
		for page := range pagesChan {
			// Parse the first line of the page to get the URL
			pageURL := strings.Split(page, "\n")[0]
			documentTags := fmt.Sprintf("web, %s", pageURL)

			// Remove any '403 Forbidden' text from the documentTags
			documentTags = strings.ReplaceAll(documentTags, "403 Forbidden", "")

			err := handleTextSplitAndIndex(documentTags, page, 1024, "avsolatorio/GIST-small-Embedding-v0")
			if err != nil {
				log.Errorf("Error handling text split and index: %v", err)
			}
			document = fmt.Sprintf("%s\n%s", document, page)
		}

		pterm.Error.Printf("Fetching web search chunks from memory...")
		document, _ = handleChatMemory(config, chatMessage)
		//pterm.Error.Printf("Web Search Document: %s\n", document)
		chatMessage = fmt.Sprintf("%s Reference the previous information if it is relevant to the next query only. Do not provide any additional information other than what is necessary to answer the next question or respond to the query. Be concise. Do not deviate from the topic of the query.\nQUERY:\n%s", document, chatMessage)

		pterm.Info.Println("Tool workflow complete")

		return chatMessage
	}

	chatMessage = fmt.Sprintf("REFERENCE DOCUMENT:\n%s\n\nQUERY:\n%s", document, chatMessage)

	pterm.Info.Println("Tool workflow complete")

	return chatMessage
}

// handleToolToggle toggles the state of various tools based on the provided tool name.
func handleToolToggle(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		toolName := c.Params("toolName")
		enabled := c.Params("enabled")
		topN := c.Params("topN")

		pterm.Info.Println(enabled)

		// Convert the enabled parameter to a boolean.
		enabledBool, err := strconv.ParseBool(enabled)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid enabled parameter")
		}

		// Convert the topN parameter to an integer.
		topNInt, err := strconv.Atoi(topN)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid topN parameter")
		}

		// Print the params to the console.
		pterm.Info.Println("Params:")
		pterm.Info.Println(toolName)

		switch toolName {
		case "memory":
			pterm.Warning.Sprintf("Memory tool toggled: %t\n", config.Tools.Memory.Enabled)
			config.Tools.Memory.Enabled = enabledBool
			config.Tools.Memory.TopN = topNInt
		case "webget":
			pterm.Warning.Sprintf("WebGet tool toggled: %t\n", config.Tools.WebGet.Enabled)
			config.Tools.WebGet.Enabled = !config.Tools.WebGet.Enabled
		case "websearch":
			pterm.Warning.Sprintf("WebSearch tool toggled: %t\n", config.Tools.WebSearch.Enabled)
			config.Tools.WebSearch.Enabled = enabledBool
			config.Tools.WebSearch.TopN = topNInt
		case "imggen":
			config.Tools.ImgGen.Enabled = enabledBool
		default:
			return c.Status(fiber.StatusNotFound).SendString("Tool not found")
		}

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Tool %s toggled", toolName)})
	}
}

// handleToolList retrieves and returns a list of tools from the configuration with all parameters.
func handleToolList(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(config.Tools)
	}
}

// ComfyUI Service Handlers
const serverAddress = "192.168.0.148:8188"

type Image struct {
	Filename  string `json:"filename"`
	Subfolder string `json:"subfolder"`
	Type      string `json:"type"`
}

type Outputs struct {
	Images []Image `json:"images"`
}

type Meta struct {
	Title string `json:"title"`
}

type Inputs struct {
	Cfg            int           `json:"cfg,omitempty"`
	Denoise        int           `json:"denoise,omitempty"`
	LatentImage    []interface{} `json:"latent_image,omitempty"`
	Model          []interface{} `json:"model,omitempty"`
	Negative       []interface{} `json:"negative,omitempty"`
	Positive       []interface{} `json:"positive,omitempty"`
	SamplerName    string        `json:"sampler_name,omitempty"`
	Scheduler      string        `json:"scheduler,omitempty"`
	Seed           int           `json:"seed,omitempty"`
	Steps          int           `json:"steps,omitempty"`
	CkptName       string        `json:"ckpt_name,omitempty"`
	BatchSize      int           `json:"batch_size,omitempty"`
	Height         int           `json:"height,omitempty"`
	Width          int           `json:"width,omitempty"`
	Clip           []interface{} `json:"clip,omitempty"`
	Text           string        `json:"text,omitempty"`
	Samples        []interface{} `json:"samples,omitempty"`
	Vae            []interface{} `json:"vae,omitempty"`
	Images         []interface{} `json:"images,omitempty"`
	FilenamePrefix string        `json:"filename_prefix,omitempty"`
}

type PromptData struct {
	Meta      Meta   `json:"_meta"`
	ClassType string `json:"class_type"`
	Inputs    Inputs `json:"inputs"`
}

type Status struct {
	Completed bool            `json:"completed"`
	Messages  [][]interface{} `json:"messages"`
	StatusStr string          `json:"status_str"`
}

type HistoryEntry struct {
	Outputs  Outputs                `json:"outputs"`
	Prompt   map[string]interface{} `json:"prompt"`
	Status   Status                 `json:"status"`
	ClientID string                 `json:"client_id,omitempty"`
}

type History struct {
	Entries HistoryEntry
}

type Prompt struct {
	Prompt   map[string]interface{} `json:"prompt"`
	ClientID string                 `json:"client_id"`
}

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Load the prompt text from json file
func loadPromptText(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func SaveBytesAsImage(data []byte, filename string) error {
	fmt.Println("Saving image to file:", filename)

	reader := bytes.NewReader(data)

	img, _, err := image.Decode(reader)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return err
	}

	out, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return err
	}

	fmt.Println("Image saved successfully to file:", filename)
	return nil
}

func queuePrompt(prompt map[string]interface{}, clientID string) (map[string]interface{}, error) {
	p := Prompt{
		Prompt:   prompt,
		ClientID: clientID,
	}
	jsonData, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/prompt", serverAddress), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("Result:", result)

	return result, err
}

func getHistory() (History, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/history", serverAddress))
	if err != nil {
		return History{}, err
	}
	defer resp.Body.Close()

	// Print the body
	body, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(body))

	var result History
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshalling history:", err)
	}

	return result, err
}

func getImage(filename, subfolder, folderType string) ([]byte, error) {
	data := map[string]string{
		"filename":  filename,
		"subfolder": subfolder,
		"type":      folderType,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/view", serverAddress), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []byte
	err = json.NewDecoder(resp.Body).Decode(&result)

	imgTurn := strconv.Itoa(chatTurn)
	imgPath := fmt.Sprintf("public/uploads/%s_%s", imgTurn, "sd_out.png")

	SaveBytesAsImage(result, imgPath)

	return result, err
}

func getImages(prompt map[string]interface{}) (map[string][][]byte, error) {

	ws, _, err := socket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", serverAddress), nil)
	if err != nil {
		fmt.Println("Error dialing websocket:", err)
		return nil, err
	}
	defer ws.Close()

	clientID := uuid.New().String()
	result, err := queuePrompt(prompt, clientID)
	if err != nil {
		return nil, err
	}

	promptID := result["prompt_id"].(string)
	outputImages := make(map[string][][]byte)
	currentNode := ""

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return nil, err
		}

		if msg[0] == '{' {
			var message Message
			err = json.Unmarshal(msg, &message)
			if err != nil {
				fmt.Println("Error unmarshalling message:", err)
				return nil, err
			}

			if message.Type == "executing" {
				data := message.Data
				if data["prompt_id"] == promptID {
					if data["node"] == nil {
						break
					} else {
						currentNode = data["node"].(string)
					}
				}
			}
		} else {
			if currentNode == "save_image_websocket_node" {
				outputImages[currentNode] = append(outputImages[currentNode], msg[8:])
			}
		}
	}

	return outputImages, nil
}

// handleRoleSelection handles the selection of assistant roles.
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
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("Role set to %s", foundRole.Name),
			})
		}

		if foundRole != nil {
			config.CurrentRoleInstructions = foundRole.Instructions
			pterm.Info.Printf("Role set to: %s\n", foundRole.Name)
			pterm.Info.Println(foundRole.Instructions)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("Role set to %s", foundRole.Name),
			})
		}

		return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
	}
}

func fetchURL(url string) ([]byte, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: %s, status code: %d", url, resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func pollURL(url string, timeout time.Duration) ([]byte, error) {
	startTime := time.Now()
	for {
		if time.Since(startTime) > timeout {
			return nil, fmt.Errorf("timed out after %v seconds", timeout.Seconds())
		}

		data, err := fetchURL(url)
		if err == nil {
			return data, nil
		}

		fmt.Println("Error fetching URL, retrying:", err)
		time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
	}
}

func performImageGen(chatId string, config *AppConfig, chatMessage string) string {

	currentChatUid = chatId
	var prompt map[string]interface{}
	workflowPath := fmt.Sprintf("%s/web/pixart.json", config.DataPath)
	promptText, err := loadPromptText(workflowPath)
	if err != nil {
		fmt.Println("Error loading prompt text:", err)
		return ""
	}

	err = json.Unmarshal([]byte(promptText), &prompt)
	if err != nil {
		fmt.Println("Error unmarshaling prompt:", err)
		return ""
	}

	// Generate a random seed between 1 and 999999999999999
	seed := rand.Intn(999999999999999)

	pterm.Info.Println("Generating image using seed:", seed)

	prompt["374"].(map[string]interface{})["inputs"].(map[string]interface{})["filename_prefix"] = chatId
	prompt["196"].(map[string]interface{})["inputs"].(map[string]interface{})["text"] = chatMessage
	prompt["206"].(map[string]interface{})["inputs"].(map[string]interface{})["text"] = chatMessage
	prompt["286"].(map[string]interface{})["inputs"].(map[string]interface{})["noise_seed"] = seed

	_, err = queuePrompt(prompt, chatId)
	if err != nil {
		log.Errorf("Error queuing prompt: %v", err)
		formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>", strconv.Itoa(chatTurn), err)
		return formattedContent
	}

	imgFileName := fmt.Sprintf("%s_00001_.png", chatId)
	imgPath := fmt.Sprintf("%s/web/uploads/%s", config.DataPath, imgFileName)

	// We want to find the image using the UID and regex since we do not know the exact turn number
	imageUrl := fmt.Sprintf("http://%s/view?filename=%s&subfolder=&type=output", serverAddress, imgFileName)
	image, err := pollURL(imageUrl, 240*time.Second) // Timeout after 240 seconds since the workflow is complex, make this configurable in future commit
	if err != nil {
		fmt.Println("Error retrieving generated image:", err)
	}

	// Save the image to a file.
	_ = SaveBytesAsImage(image, imgPath)

	imgElement := fmt.Sprintf("<img class='rounded-2 object-fit-scale' width='512' height='512' src='public/uploads/%s_00001_.png' />", chatId)
	formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>", strconv.Itoa(chatTurn), imgElement)

	chatTurn = chatTurn + 1
	return formattedContent
}
