package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"eternal/pkg/llm"
	"eternal/pkg/web"

	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
}

type CompletionResponse struct {
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

func StreamCompletionToWebSocket(c *websocket.Conn, chatID int, model string, messages []Message, temperature float64, apiKey string) error {
	payload := &CompletionRequest{
		Model:       model,
		MaxTokens:   124000,
		Messages:    messages,
		Stream:      true,
		Temperature: temperature,
	}

	pterm.Info.Println(payload)

	resp, err := SendRequest(completionsEndpoint, payload, apiKey)
	if err != nil {
		pterm.Error.Println(err)
		return err
	}
	defer resp.Body.Close()

	pterm.Warning.Println("Response status:", resp.Status)
	pterm.Warning.Println(err)

	// Handle streaming response
	msgBuffer := new(bytes.Buffer)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			jsonStr := line[6:] // Strip the "data: " prefix
			var data CompletionResponse
			if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
				return fmt.Errorf("%s", msgBuffer.String())
			}

			// Accumulate content from each choice in the buffer
			for _, content := range data.Content {
				msgBuffer.WriteString(content.Text)
			}

			// Process the accumulated content after streaming is complete
			htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())
			turnIDStr := fmt.Sprint(chatID + llm.TurnCounter)
			formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg)
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
