// INSTRUCTIONS FOR THIS EMBEDDING MODEL AT:
// https://github.com/FlagOpen/FlagEmbedding/tree/master/FlagEmbedding/llm_embedder
//
// Reference examples at:
// https://github.com/FlagOpen/FlagEmbedding/tree/master/examples
//
// TODO: Implement Reranker example

package embeddings

import (
	"eternal/pkg/documents"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	estore "eternal/pkg/vecstore"

	"github.com/pterm/pterm"
)

var modelName = "sfr-embedding-mistral/ggml-sfr-embedding-mistral-q4_k_m.gguf"

//const limit = 10

var INSTRUCTIONS = map[string]struct {
	Query string
	Key   string
}{
	"qa": {
		Query: "Represent this query for retrieving relevant documents: ",
		Key:   "Represent this document for retrieval: ",
	},
	"icl": {
		Query: "Convert this example into vector to look for useful examples: ",
		Key:   "Convert this example into vector for retrieval: ",
	},
	"chat": {
		Query: "Embed this dialogue to find useful historical dialogues: ",
		Key:   "Embed this historical dialogue for retrieval: ",
	},
	"lrlm": {
		Query: "Embed this text chunk for finding useful historical chunks: ",
		Key:   "Embed this historical text chunk for retrieval: ",
	},
	"tool": {
		Query: "Transform this user request for fetching helpful tool descriptions: ",
		Key:   "Transform this tool description for retrieval: ",
	},
	"convsearch": {
		Query: "Encode this query and context for searching relevant passages: ",
		Key:   "Encode this passage for retrieval: ",
	},
}

// Embedding represents a word embedding.
type Embedding struct {
	Word       string
	Vector     []float64
	Similarity float64
}

func GenerateEmbeddingForTask(dataPath string, task string) {

	instruction, ok := INSTRUCTIONS[task]
	if !ok {
		fmt.Printf("Unknown task: %s\n", task)
		return
	}

	// 1. Initialization
	pterm.Info.Println("Initializing...")
	db := estore.NewEmbeddingDB()

	// 2. Code Splitting
	pterm.Info.Println("Splitting code...")
	inputFilePath := os.Args[1]
	content, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	separators, _ := documents.GetSeparatorsForLanguage(documents.JSON)
	splitter := documents.RecursiveCharacterTextSplitter{
		Separators:       separators,
		KeepSeparator:    true,
		IsSeparatorRegex: false,
		ChunkSize:        1000,
		LengthFunction:   func(s string) int { return len(s) },
	}
	chunks := splitter.SplitText(string(content))

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
		fmt.Print(instruction.Query)
		fmt.Println(chunk)

		strVec, err := Encoder(dataPath, chunk) // Invoke the encoder function with the chunk
		if err != nil {
			pterm.Error.Println("Error generating embedding for chunk:", err)
			panic(err)
		}

		// Convert the string vector to a float64 vector
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
				return
			}
			vec = append(vec, val)
		}

		fmt.Println("Vector:")
		fmt.Println(vec)

		embedding := estore.Embedding{
			Word:       chunk,
			Vector:     vec,
			Similarity: 0.0,
		}

		fmt.Println("Embedding:", embedding)

		db.AddEmbedding(embedding)
	}

	// Save the database to a file
	pterm.Info.Println("Saving embeddings...")
	db.SaveEmbeddings("./embeddings.db")
}

// # "/Users/arturoaquino/.eternal-v1/gguf/embedding"
func Encoder(dataPath string, text string) (string, error) {
	cmdPath := fmt.Sprintf("%s/embedding", dataPath)
	modelPath := fmt.Sprintf("%s/models/%s", dataPath, modelName)
	cmd := exec.Command(cmdPath, "-c", "4096", "-m", modelPath, "--log-disable", "-p", text)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdoutPipe.Close() // Close the pipe after starting the command

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var outputBuilder strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := stdoutPipe.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if n > 0 {
			outputBuilder.Write(buf[:n])
		}
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return outputBuilder.String(), nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func Search(dataPath string, dbName string, prompt string, topN int) []estore.Embedding {
	// Initialize DB and load embeddings
	db := estore.NewEmbeddingDB()
	dbPath := fmt.Sprintf("%s/%s", dataPath, dbName)
	embeddings, err := db.LoadEmbeddings(dbPath)
	if err != nil {
		fmt.Println("Error loading embeddings:", err)
		return nil
	}

	// Generate embedding for the prompt
	cmdPath := fmt.Sprintf("%s/gguf", dataPath)
	strVec, err := Encoder(cmdPath, prompt)
	if err != nil {
		fmt.Println("Error generating embedding for prompt:", err)
		return nil
	}

	// Convert the string vector to a float64 vector
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
			return nil
		}
		vec = append(vec, val)
	}

	embeddingForPrompt := estore.Embedding{
		Word:       prompt,
		Vector:     vec,
		Similarity: 0.0,
	}

	// Retrieve the top N similar embeddings
	return estore.FindTopNSimilarEmbeddings(embeddingForPrompt, embeddings, topN)
}
