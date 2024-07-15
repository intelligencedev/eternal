package embeddings

import (
	"bytes"
	"encoding/json"
	"eternal/pkg/documents"
	"eternal/pkg/llm/openai"
	"eternal/pkg/vecstore"
	"fmt"
	"net/http"
	"os"

	"github.com/pterm/pterm"
)

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

// GetEmbeddings interacts with the OpenAI Embeddings API to retrieve embeddings based on the provided request.
// It returns an EmbedResponse pointer and any error encountered during the API call.
func GetEmbeddings(req EmbedRequest) (*EmbedResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	client := openai.NewClient(apiKey)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}

		// Additional logic to handle Code as a number or string
		var codeStr string
		switch code := errorResponse.Error.Code.(type) {
		case float64:
			codeStr = fmt.Sprintf("%.0f", code) // Convert number to string
		case string:
			codeStr = code
		default:
			return nil, fmt.Errorf("unexpected type for error code")
		}

		return nil, fmt.Errorf("API error: %s (Code: %s)", errorResponse.Error.Message, codeStr)
	}

	var embedResponse EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode successful response: %w", err)
	}

	return &embedResponse, nil
}

func GenerateEmbeddingOAI() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <path_to_input_file>")
		return
	}

	// 1. Initialization
	pterm.Info.Println("Initializing...")

	// Create a new OpenAI client
	//client := client()

	db := vecstore.NewEmbeddingDB()

	// 2. Code Splitting
	pterm.Info.Println("Splitting code...")
	inputFilePath := os.Args[1]
	content, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	separators, _ := documents.GetSeparatorsForLanguage(documents.JSON)
	// Updated the RecursiveCharacterTextSplitter to include OverlapSize and updated SplitText method
	splitter := documents.RecursiveCharacterTextSplitter{
		Separators:       separators,
		KeepSeparator:    true,
		IsSeparatorRegex: false,
		ChunkSize:        1000,
		LengthFunction:   func(s string) int { return len(s) },
	}
	chunks := splitter.SplitText(string(content))

	// 3. Embedding Generation
	pterm.Info.Println("Generating embeddings...")
	for _, chunk := range chunks {
		req := EmbedRequest{
			Model: "text-embedding-ada-002",
			Input: chunk,
		}
		resp, err := GetEmbeddings(req)
		if err != nil {
			fmt.Printf("Error getting embeddings: %v\n", err)
			panic(err)
		}

		response := EmbedData{
			Object:    resp.Object,
			Embedding: resp.Data[0].Embedding,
			Index:     resp.Data[0].Index,
		}

		embedding := vecstore.Embedding{
			Word:       chunk,
			Vector:     response.Embedding,
			Similarity: 0.0,
		}

		db.AddEmbedding(embedding)

	}

	// Save the database to a file
	pterm.Info.Println("Saving embeddings...")
	db.SaveEmbeddings("./db/embeddings.json")

	if len(chunks) > 0 {
		embedding, ok := db.RetrieveEmbedding(chunks[0])
		if ok {
			fmt.Printf("Embedding for the first chunk:\n%v\n", embedding)
		}
	}
}
