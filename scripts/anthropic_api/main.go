package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func main() {
	// Set your API key in an environment variable for security reasons
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Println("ANTHROPIC_API_KEY environment variable is not set.")
		os.Exit(1)
	}

	// Define the request body
	data := `{
		"model": "claude-3-opus-20240229",
		"max_tokens": 1024,
		"stream": true,
		"messages": [
			{"role": "user", "content": "Hello, world"}
		]
	}`

	// Create a new request
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Add headers to the request
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")
	req.Header.Add("content-type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Response status code: %d\n", resp.StatusCode)

	// Parse the response body into the Response struct
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error unmarshalling response body: %v\n", err)
		os.Exit(1)
	}

	// Print the raw response body
	fmt.Println("Raw Response:")
	fmt.Println(string(body))

	// Print the response in a more readable format
	fmt.Println("Response:")
	fmt.Printf("ID: %s\n", response.ID)
	fmt.Printf("Type: %s\n", response.Type)
	fmt.Printf("Role: %s\n", response.Role)
	for _, content := range response.Content {
		fmt.Printf("Content: %s\n", content.Text)
	}
	fmt.Printf("Model: %s\n", response.Model)
	fmt.Printf("Stop Reason: %s\n", response.StopReason)
	fmt.Printf("Usage - Input Tokens: %d, Output Tokens: %d\n", response.Usage.InputTokens, response.Usage.OutputTokens)
}
