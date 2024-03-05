package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"eternal/pkg/llm"
	"eternal/pkg/web"

	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
)

const (
	baseURL             = "https://api.openai.com/v1"
	completionsEndpoint = "/chat/completions"
)

// SendRequest sends a request to the OpenAI API and decodes the response.
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
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return http.DefaultClient.Do(req)
}

func StreamCompletionToWebSocket(c *websocket.Conn, chatID int, model string, messages []llm.Message, temperature float64, apiKey string) error {
	payload := &CompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
		Stream:      true,
	}

	resp, err := SendRequest(completionsEndpoint, payload, apiKey)
	if err != nil {
		pterm.Error.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Handle streaming response
	msgBuffer := new(bytes.Buffer)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			jsonStr := line[6:] // Strip the "data: " prefix
			var data struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
				FinishReason string `json:"finish_reason"`
			}

			if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
				return fmt.Errorf("%s", msgBuffer.String())
			}

			// Accumulate content from each choice in the buffer
			for _, choice := range data.Choices {
				msgBuffer.WriteString(choice.Delta.Content)
			}

			// Process the accumulated content after streaming is complete
			htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())

			turnIDStr := fmt.Sprint(chatID + llm.TurnCounter)

			// TODO: Abstract this into a function that all backends use.
			formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet sandbox='python' editor='external'></codapi-snippet>", turnIDStr, htmlMsg)

			if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
				pterm.Error.Println("WebSocket write error:", err)
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		pterm.Error.Println("Error reading stream:", err)
		return err
	}

	return nil
}
