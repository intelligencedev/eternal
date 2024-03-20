package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"eternal/pkg/llm/embeddings"
	store "eternal/pkg/vecstore"
)

func main() {
	// Check if the user provided the prompt argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <prompt>")
		return
	}

	// Extract the prompt from command line arguments
	prompt := os.Args[1]

	// Generate embedding for the prompt
	strVec, err := embeddings.Encoder(prompt)
	if err != nil {
		fmt.Println("Error generating embedding for prompt:", err)
		return
	}

	// Convert the string vector to a float64 vector
	vec := parseVector(strVec)

	embedding := store.Embedding{
		Word:   prompt,
		Vector: vec,
	}

	topN := 13 // For instance, retrieve top 3 results. Change as needed.
	topEmbeddings := Search("/Users/arturoaquino/.eternal-v1/models/dolphin-phi2/dolphin-2_6-phi-2.Q8_0.gguf", "dolphin-2_6-phi-2.Q8_0.gguf", embedding, topN)

	fmt.Println("Top", topN, "similar words or chunks for the given prompt are:")
	for _, emb := range topEmbeddings {
		fmt.Print(emb.Word)
	}
}

func Search(modelPath string, modelName string, embedding store.Embedding, topN int) []store.Embedding {
	db := store.NewEmbeddingDB()
	embeddings, err := db.LoadEmbeddings("/Users/arturoaquino/Documents/eternal/examples/embeddings/embeddings.db")
	if err != nil {
		fmt.Println("Error loading embeddings:", err)
		return nil
	}

	// Retrieve the top N similar embeddings
	topEmbeddings := store.FindTopNSimilarEmbeddings(embedding, embeddings, topN)
	if len(topEmbeddings) == 0 {
		fmt.Println("No similar embeddings found")
		return nil
	}

	return topEmbeddings
}

func parseVector(strVec string) []float64 {
	strVec = strings.TrimSpace(strVec)
	parts := strings.Split(strVec, " ")
	var vec []float64

	for _, part := range parts {
		if part == "" || !isNumeric(part) {
			continue
		}

		val, err := strconv.ParseFloat(part, 64)
		if err != nil {
			fmt.Printf("Error parsing float for part '%s': %v\n", part, err)
			continue
		}
		vec = append(vec, val)
	}

	return vec
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
