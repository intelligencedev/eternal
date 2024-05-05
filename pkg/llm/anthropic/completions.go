package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"eternal/pkg/web"
	"fmt"
	"strconv"
	"strings"

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

type ContentBlockDelta struct {
	Type  string    `json:"type"`
	Index int       `json:"index"`
	Delta TextDelta `json:"delta"`
}

type TextDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func StreamCompletionToWebSocket(c *websocket.Conn, chatID int, model string, messages []Message, temperature float64, apiKey string) error {

	payload := &CompletionRequest{
		Model:       model,
		MaxTokens:   4096,
		Stream:      true,
		Messages:    messages,
		Temperature: temperature,
	}

	resp, err := SendRequest(completionsEndpoint, payload, apiKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle streaming response
	msgBuffer := new(bytes.Buffer)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")

			// Unmarshal the JSON response
			var data ContentBlockDelta
			if err := json.Unmarshal([]byte(line), &data); err != nil {
				return err
			}

			msgBuffer.WriteString(data.Delta.Text)

			htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())

			turnIDStr := strconv.Itoa(chatID)

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
