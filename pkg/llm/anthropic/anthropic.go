package anthropic

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	baseURL             = "https://api.anthropic.com/v1"
	completionsEndpoint = "/messages"
)

// SendRequest sends a request to the Anthropic API and decodes the response.
func SendRequest(endpoint string, payload interface{}, apiKey string) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	return http.DefaultClient.Do(req)
}
