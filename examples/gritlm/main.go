package main

import (
	"eternal/pkg/documents"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	estore "eternal/pkg/vecstore"

	"github.com/pterm/pterm"
)

var (
	dataPath      = flag.String("data-path", ".", "The path to save the embeddings file to")
	modelPathFlag = flag.String("model-path", ".eternal/models/", "The path to the model directory")
	modelNameFlag = flag.String("model-name", "gritlm-7b/ggml-gritlm-7b-q8_0.gguf", "The name of the model")

	generateCommand = flag.NewFlagSet("generate", flag.ExitOnError)
	inputFileFlag   = generateCommand.String("input-file", "", "The input file to generate embeddings for")
	chunkSize       = generateCommand.Int("chunk-size", 500, "The size of the chunk to generate embeddings for. Lower the size if an error is returned.")
	overlapSize     = generateCommand.Int("overlap-size", 200, "The size of the overlap between chunks")

	promptFlag      = retrieveCommand.String("prompt", "", "The prompt to retrieve similar words or chunks for")
	topNFlag        = retrieveCommand.Int("top-n", 5, "The number of top similar words or chunks to retrieve")
	retrieveCommand = flag.NewFlagSet("retrieve", flag.ExitOnError)
)

// Embedding represents a word embedding.
type Embedding struct {
	Word   string
	Vector []float64
}

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
	// usr, err := user.Current()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Set the path to the model
	// model := fmt.Sprintf("%s/%s/%s", usr.HomeDir, *modelPathFlag, *modelNameFlag)
	model := "/Users/arturoaquino/.eternal-v1/models/gritlm-7b/ggml-gritlm-7b-q8_0.gguf"

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

			GenerateEmbeddingForTask("qa", document, "txt", *chunkSize, *overlapSize, model, *dataPath)
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

func GenerateEmbeddingForTask(task string, content string, doctype string, chunkSize int, overlapSize int, model string, dataPath string) error {

	db := estore.NewEmbeddingDB()

	var chunks []string
	var separators []string

	if doctype == "txt" {
		chunks = documents.SplitTextByCount(string(content), chunkSize)
	} else {
		doctype = strings.ToUpper(doctype)
		separators, _ = documents.GetSeparatorsForLanguage(documents.Language(doctype))

		overlapSize := chunkSize / 2 // Set the overlap size to half of the chunk size

		splitter := documents.RecursiveCharacterTextSplitter{
			Separators:       separators,
			KeepSeparator:    true,
			IsSeparatorRegex: false,
			ChunkSize:        chunkSize,
			OverlapSize:      overlapSize, // Add the OverlapSize field
			LengthFunction:   func(s string) int { return len(s) },
		}
		chunks = splitter.SplitText(string(content))
	}

	// Remove duplicate chunks
	seen := make(map[string]bool)
	var uniqueChunks []string
	for _, chunk := range chunks {

		if _, ok := seen[chunk]; !ok {
			uniqueChunks = append(uniqueChunks, chunk)
			seen[chunk] = true
		}
	}

	// 3. Embedding Generation
	pterm.Info.Println("Generating embeddings...")
	for _, chunk := range uniqueChunks {
		var vec []float64

		result, err := getGritEmbedding(chunk, model)
		if err != nil {
			log.Printf("Error generating embedding for chunk: %v", err)
			continue
		}

		// Split the result by spaces into the vec variable
		vecComponents := strings.Split(result, " ")

		for _, strVal := range vecComponents {
			if strVal != "" {
				floatVal, err := strconv.ParseFloat(strVal, 64) // Convert string to float64

				if err != nil {
					log.Printf("Error converting string to float64: %v", err)
					continue // Skip this value and move to the next
				}

				vec = append(vec, floatVal) // Append the float64 value to the vec slice
			}
		}

		vec = normalizeVector(vec)

		embedding := estore.Embedding{
			Word:   chunk,
			Vector: vec,
		}

		fmt.Print(embedding)

		db.AddEmbedding(embedding)
	}

	// Save the database to a file
	pterm.Info.Println("Saving embeddings...")

	dbPath := fmt.Sprintf("%s/embeddings.db", dataPath)

	db.SaveEmbeddings(dbPath)

	return nil
}

func getGritEmbedding(inputText string, model string) (string, error) {
	// Construct the command
	// cmd := exec.Command(
	// 	"/Users/arturoaquino/.eternal-v1/gguf/embedding",
	// 	"-m", model,
	// 	"-p", inputText,
	// 	"--embedding",
	// 	"--log-disable",
	// 	"--no-display-prompt",
	// )

	// Specify the Unicode character (e.g., Unicode Replacement Character U+FFFD)
	unicodeChar := "\uFFFD"

	// Replace newlines with the Unicode character
	inputText = strings.ReplaceAll(inputText, "\n", unicodeChar)

	cmd := exec.Command(
		"/Users/arturoaquino/.eternal-v1/gguf/embedding",
		"-m", model,
		"-p", inputText,
		"--embedding",
		"--log-disable",
		"--no-display-prompt",
	)

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	outputStr := string(output)

	vecs := strings.Split(outputStr, "embedding 0: ")
	//fmt.Println(vecs[1])

	return vecs[1], nil
}

func Search(modelPath string, prompt string, topN int) []estore.Embedding {
	db := estore.NewEmbeddingDB()
	embeddings, err := db.LoadEmbeddings("./embeddings.db")
	if err != nil {
		fmt.Println("Error loading embeddings:", err)
		return nil
	}

	var vec []float64

	result, err := getGritEmbedding(prompt, "/Users/arturoaquino/.eternal-v1/models/gritlm-7b/ggml-gritlm-7b-q8_0.gguf")
	if err != nil {
		log.Printf("Error generating embedding for chunk: %v", err)
		panic(err)
	}

	// Split the result by spaces into the vec variable
	vecComponents := strings.Fields(result) // Split the result into components

	for _, strVal := range vecComponents {
		floatVal, err := strconv.ParseFloat(strVal, 64) // Convert string to float64
		if err != nil {
			log.Printf("Error converting string to float64: %v", err)
			continue // Skip this value and move to the next
		}
		vec = append(vec, floatVal) // Append the float64 value to the vec slice
	}

	promptEmbeddings := estore.Embedding{
		Word:       prompt,
		Vector:     vec,
		Similarity: 0,
	}

	// Retrieve the top N similar embeddings
	topEmbeddings := estore.FindTopNSimilarEmbeddings(promptEmbeddings, embeddings, topN)
	if len(topEmbeddings) == 0 {
		fmt.Println("Error finding similar embeddings.")
		return nil
	}

	return topEmbeddings
}

func normalizeVector(vec []float64) []float64 {
	var norm float64
	for _, val := range vec {
		norm += val * val
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		return vec // Return the original vector if norm is 0 to avoid division by zero
	}
	for i, val := range vec {
		vec[i] = val / norm
	}
	return vec
}
