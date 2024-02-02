package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pterm/pterm"
)

// Client represents the HTTP client for interacting with the LLM API.
type Client struct {
	APIKey string
	HTTP   *http.Client
}

// OAIModel represents a single model in the JSON response.
type OAIModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelsResponse represents the top-level structure of the JSON response.
type ModelsResponse struct {
	Object string     `json:"object"`
	Data   []OAIModel `json:"data"`
}

// NewClient creates and initializes a new instance of an LLM API client using the provided API key.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTP:   &http.Client{},
	}
}

// Initialize sets up the LLM API client using the API key from environment variables.
func (c *Client) Initialize() error {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		pterm.Error.Println("Please set the LLM_API_KEY environment variable.")
		return fmt.Errorf("LLM_API_KEY is not set")
	}

	c.APIKey = apiKey
	return nil
}

func (c *Client) Connect(endpoint string) (*http.Response, error) {
	// Include the base URL and the protocol scheme
	fullURL := "https://api.openai.com/v1" + endpoint // Example base URL

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", fullURL, err)
	}

	return resp, nil
}

// GetModels retrieves the available models from the OpenAI API.
// It returns a ModelsResponse and any error encountered.
func GetModels(client *Client) (ModelsResponse, error) {
	resp, err := client.Connect("/models")
	if err != nil {
		return ModelsResponse{}, fmt.Errorf("failed to connect to OpenAI API: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return ModelsResponse{}, fmt.Errorf("failed to fetch models: non-OK status received - %s", resp.Status)
	}

	var models ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return ModelsResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return models, nil
}
