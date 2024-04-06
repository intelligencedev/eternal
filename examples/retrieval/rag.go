package main

import (
	"context"
	store "eternal/pkg/vecstore"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
	"github.com/pterm/pterm"
)

const modelName = "BAAI/bge-large-en-v1.5"
const limit = 128

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		pterm.Error.Printf("Error getting user home directory: %s\n", err)
		os.Exit(1)
	}

	dataPath := filepath.Join(homeDir, ".eternal-v1/")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		err = os.Mkdir(dataPath, 0755)
		if err != nil {
			pterm.Error.Printf("Error creating data directory: %s\n", err)
		}
	}

	modelPath := filepath.Join(dataPath, "data/models/HF/")

	// Check if the user provided the prompt argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <prompt>")
		return
	}

	// Extract the prompt from command line arguments
	prompt := os.Args[1]

	// Use the Search function to find the top N most similar words or chunks
	topN := 5 // For instance, retrieve top 5 results. Change as needed.
	topEmbeddings := Search(modelPath, modelName, prompt, topN)

	fmt.Println("Top", topN, "similar words or chunks for the given prompt are:")
	for _, embedding := range topEmbeddings {
		fmt.Println(embedding.Word, "-", embedding.Similarity)
	}

}

func Search(modelPath string, modelName string, prompt string, topN int) []store.Embedding {
	db := store.NewEmbeddingDB()
	embeddings, err := db.LoadEmbeddings("./embeddings.db")
	if err != nil {
		fmt.Println("Error loading embeddings:", err)
		return nil
	}

	model, err := tasks.Load[textencoding.Interface](&tasks.Config{ModelsDir: modelPath, ModelName: modelName})
	if err != nil {
		fmt.Println("Error loading model:", err)
		return nil
	}

	var vec []float64
	result, err := model.Encode(context.Background(), prompt, int(bert.MeanPooling))
	if err != nil {
		fmt.Println("Error encoding text:", err)
		return nil
	}
	vec = result.Vector.Data().F64()[:limit]

	embeddingForPrompt := store.Embedding{
		Word:       prompt,
		Vector:     vec,
		Similarity: 0.0,
	}

	// Retrieve the top N similar embeddings
	topEmbeddings := store.FindTopNSimilarEmbeddings(embeddingForPrompt, embeddings, topN)
	if len(topEmbeddings) == 0 {
		fmt.Println("Error finding similar embeddings.")
		return nil
	}

	return topEmbeddings
}
