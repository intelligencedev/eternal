package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	embeddings "eternal/pkg/embeddings"
)

var modelPath = "data/models/HF/"

// var modelName = "BAAI/bge-large-en-v1.5"
var modelName = "avsolatorio/GIST-small-Embedding-v0"

func main() {

	// Get user home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(usr)

	// Set the path to the model
	model := fmt.Sprintf("%s/%s/%s", usr, modelPath, modelName)

	// Get prompt as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <input file>")
		return
	}

	embeddings.GenerateEmbeddingForTask("qa", "txt", model)

}
