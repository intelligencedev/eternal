package claude

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
	baseURL             = "https://api.anthropic.com/v1"
	completionsEndpoint = "/messages"
)

// SendRequest sends a request to the Claude Opus API and decodes the response.
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
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", "messages-2023-12-15")
	req.Header.Set("X-API-Key", apiKey)
	return http.DefaultClient.Do(req)
}

func StreamCompletionToWebSocket(c *websocket.Conn, chatID int, model string, messages []llm.Message, apiKey string) error {
	payload := &CompletionRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: 256,
		Stream:    true,
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
				Type         string            `json:"type"`
				Delta        ContentBlockDelta `json:"delta"`
				StopReason   string            `json:"stop_reason"`
				StopSequence interface{}       `json:"stop_sequence"`
				Usage        struct {
					OutputTokens int `json:"output_tokens"`
				} `json:"usage"`
			}
			if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
				return fmt.Errorf("%s", msgBuffer.String())
			}

			// Accumulate content from each delta in the buffer
			if data.Type == "content_block_delta" {
				msgBuffer.WriteString(data.Delta.Text)
			}

			// Process the accumulated content after streaming is complete
			if data.Type == "message_stop" {
				htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())
				turnIDStr := fmt.Sprint(chatID + llm.TurnCounter)
				formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg)
				if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
					pterm.Error.Println("WebSocket write error:", err)
					return err
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		pterm.Error.Println("Error reading stream:", err)
		return err
	}
	return nil
}
