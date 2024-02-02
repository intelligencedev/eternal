package openai

import "eternal/pkg/llm"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Model represents an AI model from the OpenAI API with its ID, name, and description.
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CompletionRequest represents the payload for the completion API.
type CompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []llm.Message `json:"messages"`
	Temperature float64       `json:"temperature"`
	Stream      bool          `json:"stream"`
}

// Choice represents a choice for the completion response.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Logprobs     *bool   `json:"logprobs"` // Pointer to a boolean or nil
	FinishReason string  `json:"finish_reason"`
}

// Usage contains information about token usage in the completion response.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CompletionResponse represents the response from the completion API.
type CompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

// AudioSpeechRequest represents the payload for the audio speech API.
type AudioSpeechRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

// ErrorData represents the structure of an error response from the OpenAI API.
type ErrorData struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
}

// ErrorResponse wraps the structure of an error when an API request fails.
type ErrorResponse struct {
	Error ErrorData `json:"error"`
}

// EmbedRequest encapsulates the request data for the OpenAI Embeddings API.
type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// EmbedResponse contains the response data from the Embeddings API.
type EmbedResponse struct {
	Object string       `json:"object"`
	Data   []EmbedData  `json:"data"`
	Model  string       `json:"model"`
	Usage  UsageMetrics `json:"usage"`
}

// EmbedData represents a single embedding and its associated data.
type EmbedData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// UsageMetrics details the token usage of the Embeddings API request.
type UsageMetrics struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
