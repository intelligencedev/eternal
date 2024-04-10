package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	embeddings "eternal/pkg/embeddings"
	store "eternal/pkg/vecstore"

	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
)

// Chunk size should be less than the max tokens for the model used: https://huggingface.co/spaces/mteb/leaderboard

var (
	modelPathFlag = flag.String("model-path", ".eternal/models/HF/", "The path to the model directory")
	modelNameFlag = flag.String("model-name", "avsolatorio/GIST-small-Embedding-v0", "The name of the model")
	limitFlag     = flag.Int("limit", 128, "The limit for the number of dimensions in the embedding vector")

	generateCommand = flag.NewFlagSet("generate", flag.ExitOnError)
	inputFileFlag   = generateCommand.String("input-file", "", "The input file to generate embeddings for")
	chunkSize       = generateCommand.Int("chunk-size", 500, "The size of the chunk to generate embeddings for. Lower the size if an error is returned.")
	overlapSize     = generateCommand.Int("overlap-size", 200, "The size of the overlap between chunks")

	retrieveCommand = flag.NewFlagSet("retrieve", flag.ExitOnError)
	promptFlag      = retrieveCommand.String("prompt", "", "The prompt to retrieve similar words or chunks for")
	topNFlag        = retrieveCommand.Int("top-n", 5, "The number of top similar words or chunks to retrieve")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: main.go <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  generate - Generate embeddings for the specified input file")
		fmt.Println("  retrieve - Retrieve top N similar words or chunks for the given prompt")
		return
	}

	command := flag.Arg(0)

	// Get user home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Set the path to the model
	model := fmt.Sprintf("%s/%s/%s", usr.HomeDir, *modelPathFlag, *modelNameFlag)

	switch command {
	case "generate":
		generateCommand.Parse(flag.Args()[1:])
		if generateCommand.Parsed() {
			if *inputFileFlag == "" {
				fmt.Println("Usage: main.go generate --input-file <input file>")
				return
			}

			fmt.Println("Generating embeddings for the input file:", *inputFileFlag)

			// Open the file and read the contents into document var
			file, err := os.Open(*inputFileFlag)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			// Convert file contents to string
			document := ""
			buf := make([]byte, 1024)
			for {
				n, err := file.Read(buf)
				if n == 0 || err != nil {
					break
				}
				document += string(buf[:n])
			}

			embeddings.GenerateEmbeddingForTask("qa", document, "json", *chunkSize, *overlapSize, model)
		}
	case "retrieve":
		retrieveCommand.Parse(flag.Args()[1:])
		if retrieveCommand.Parsed() {
			if *promptFlag == "" {
				fmt.Println("Usage: main.go retrieve --prompt <prompt>")
				return
			}
			topEmbeddings := Search(model, *promptFlag, *topNFlag)
			fmt.Println("Top", *topNFlag, "similar words or chunks for the given prompt are:")
			for _, embedding := range topEmbeddings {
				fmt.Println(embedding.Word, "-", embedding.Similarity)
			}
		}
	default:
		fmt.Println("Invalid command. Available commands: generate, retrieve")
	}
}

func Search(modelPath string, prompt string, topN int) []store.Embedding {
	db := store.NewEmbeddingDB()
	embeddings, err := db.LoadEmbeddings("./embeddings.db")
	if err != nil {
		fmt.Println("Error loading embeddings:", err)
		return nil
	}

	model, err := tasks.Load[textencoding.Interface](&tasks.Config{ModelsDir: modelPath, ModelName: *modelNameFlag})
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
	vec = result.Vector.Data().F64()[:*limitFlag]

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
