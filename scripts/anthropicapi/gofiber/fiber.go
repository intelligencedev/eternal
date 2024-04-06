package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string      `json:"model"`
	StopReason   string      `json:"stop_reason"`
	StopSequence interface{} `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func SendMessage(c *fiber.Ctx) error {
	// Set your API key in an environment variable for security reasons
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return c.SendString("ANTHROPIC_API_KEY environment variable is not set.")
	}

	// Define the request body
	data := `{
		"model": "claude-3-opus-20240229",
		"max_tokens": 1024,
		"messages": [
			{"role": "user", "content": "Hello, world"}
		]
	}`

	// Create a new request
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error creating request: %v", err))
	}

	// Add headers to the request
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")
	req.Header.Add("content-type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error executing request: %v", err))
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error reading response body: %v", err))
	}

	// Parse the response body into the Response struct
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error unmarshalling response body: %v", err))
	}

	// Respond with the parsed response in a more readable format
	return c.JSON(response)
}

func main() {
	// Initialize a new Fiber instance
	app := fiber.New()

	// Define the route
	app.Post("/send-message", SendMessage)

	// Start the server on port 8884
	err := app.Listen(":8884")
	if err != nil {
		panic(err)
	}
}
