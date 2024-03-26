package claude

import "eternal/pkg/llm"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Model represents an AI model from the Claude Opus API with its ID, name, and description.
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CompletionRequest represents the payload for the completion API.
type CompletionRequest struct {
	Model     string        `json:"model"`
	Messages  []llm.Message `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
	Stream    bool          `json:"stream"`
}

// Usage represents the token usage in the completion response.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ContentBlockDelta represents a delta update for a content block.
type ContentBlockDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// CompletionResponse represents the response from the completion API.
type CompletionResponse struct {
	Type       string `json:"type"`
	Message    Message
	StopReason string `json:"stop_reason"`
	Usage      Usage  `json:"usage"`
}

// ErrorData represents the structure of an error response from the Claude Opus API.
type ErrorData struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
}

// ErrorResponse wraps the structure of an error when an API request fails.
type ErrorResponse struct {
	Error ErrorData `json:"error"`
}
