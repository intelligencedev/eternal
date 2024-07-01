package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"eternal/pkg/llm"
	"eternal/pkg/web"

	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
)

const (
	baseURL             = "https://api.openai.com/v1"
	completionsEndpoint = "/chat/completions"
	ttsEndpoint         = "/audio/speech"
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

func StreamCompletionToWebSocket(c websocket.Conn, chatID int, model string, messages []llm.Message, temperature float64, apiKey string, responseBuffer *bytes.Buffer) error {
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

			//the stream terminated by a data: [DONE] message.

			// If the stream is done, break out of the loop
			if strings.Contains(jsonStr, "[DONE]") {
				pterm.Error.Println("OpenAI stream completed")
				return fmt.Errorf("OpenAI stream completed")
			}

			if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
				return fmt.Errorf("%s", responseBuffer.String())
			}

			// Accumulate content from each choice in the buffer
			for _, choice := range data.Choices {
				responseBuffer.WriteString(choice.Delta.Content)
			}

			// Process the accumulated content after streaming is complete
			htmlMsg := web.MarkdownToHTML(responseBuffer.Bytes())

			turnIDStr := fmt.Sprint(chatID + llm.TurnCounter)

			// TODO: Abstract this into a function that all backends use.
			//formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet url='http://localhost:1313/v1/exec' sandbox='go' editor='external'></codapi-snippet>", turnIDStr, htmlMsg)
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

	return err
}

// StreamTTSToFile streams TTS response to a file.
func StreamTTSToFile(inputText, voice, apiKey, outputFilePath string) error {
	payload := &AudioSpeechRequest{
		Model: "tts-1",
		Input: inputText,
		Voice: voice,
	}

	resp, err := SendRequest(ttsEndpoint, payload, apiKey)
	if err != nil {
		pterm.Error.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Create the output file
	file, err := os.Create(outputFilePath)
	if err != nil {
		pterm.Error.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Stream the response to the file
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if _, err := file.Write(line); err != nil {
			pterm.Error.Println("Error writing to file:", err)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		pterm.Error.Println("Error reading stream:", err)
		return err
	}

	return nil
}
