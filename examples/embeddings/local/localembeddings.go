package main

import (
	emb "eternal/pkg/llm/embeddings"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: main.go <data path> <prompt>")
		return
	}

	dataPath := os.Args[1]
	prompt := os.Args[2]
	dbName := "embeddings.db"
	topN := 3

	topEmbeddings := emb.Search(dataPath, dbName, prompt, topN)

	// Display the results
	for i, embedding := range topEmbeddings {
		fmt.Printf("Rank %d: Word: %s, Similarity: %f\n", i+1, embedding.Word, embedding.Similarity)
	}
}

// Encoder is a placeholder for a real encoder function that would generate an embedding for the given text.
func Encoder(dataPath string, text string) (string, error) {
	// Implement the actual encoding logic here
	return "", nil
}

// stringToVector converts a space-separated string of numbers into a slice of float64.
func stringToVector(strVec string) ([]float64, error) {
	parts := strings.Fields(strVec)
	vec := make([]float64, len(parts))
	var err error
	for i, part := range parts {
		vec[i], err = strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting string to float64: %v", err)
		}
	}
	return vec, nil
}
