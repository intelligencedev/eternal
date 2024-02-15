package main

import (
	"fmt"
	"os"

	embeddings "eternal/pkg/embeddings"
)

func main() {

	// Get prompt as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		return
	}

	embeddings.GenerateEmbeddingForTask("qa")

}
