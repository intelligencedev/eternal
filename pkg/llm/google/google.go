package google

import (
	"bytes"
	"context"
	"fmt"

	"eternal/pkg/llm"
	"eternal/pkg/web"

	"github.com/gofiber/websocket/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/pterm/pterm"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	model = "models/gemini-1.5-pro-latest"
)

// StreamGeminiResponseToWebSocket streams the response from the Gemini API to a WebSocket connection.
func StreamGeminiResponseToWebSocket(c websocket.Conn, chatID int, prompt string, apiKey string) error {
	pterm.Warning.Printfln("Using model: %s", model)
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		pterm.Error.Println(err)
		return err
	}
	defer client.Close()

	pterm.Warning.Printfln("Sending prompt to api...")
	generativeModel := client.GenerativeModel(model)

	// Configure model parameters by invoking Set* methods on the model.
	generativeModel.SetTemperature(0.1)
	generativeModel.SetTopK(1)
	generativeModel.SetTopP(1)

	pterm.Warning.Printfln("Generating content stream...")
	iter := generativeModel.GenerateContentStream(ctx, genai.Text(prompt))

	msgBuffer := new(bytes.Buffer)
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			return err
		}
		if err != nil {
			pterm.Error.Println(err)
			return err
		}

		// Access the Content field of the genai.Part type
		p := message.NewPrinter(language.English)
		content := p.Sprintf("%s", resp.Candidates[0].Content.Parts[0])

		msgBuffer.WriteString(content)

		htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())

		turnIDStr := fmt.Sprint(chatID + llm.TurnCounter)

		formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg)

		if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
			pterm.Error.Println("WebSocket write error:", err)
			return err
		}
	}
}
