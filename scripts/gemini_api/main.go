package main

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const MODEL = "models/gemini-1.5-pro-latest"

func main() {
	// Get prompt as cli parameter
	if len(os.Args) < 2 {
		log.Fatal("Please provide a prompt")
	}

	prompt := os.Args[1]

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel(MODEL)

	iter := model.GenerateContentStream(ctx, genai.Text(prompt))
	p := message.NewPrinter(language.English)
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		p.Print(resp.Candidates[0].Content.Parts[0])
	}
}
