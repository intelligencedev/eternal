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
	"os"
	"os/exec"
	"strconv"
	"strings"

	estore "eternal/pkg/vecstore"

	"github.com/pterm/pterm"
)

// var modelPath = "./data/models/HF/"
// var modelName = "BAAI/bge-large-en-v1.5"
var modelName = "dolphin-2_6-phi-2.Q8_0.gguf"

//var modelName = "BAAI/llm-embedder"

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

func GenerateEmbeddingForTask(task string) {

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
	// Updated the RecursiveCharacterTextSplitter to include OverlapSize and updated SplitText method
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

	// modelsDir := fmt.Sprintf("%s/data/models/HF/BAAI/bge-large-en-v1.5/", dataPath)
	//modelsDir := fmt.Sprintf("%s/models/dolphin-phi2/", dataPath)

	// 3. Embedding Generation
	pterm.Info.Println("Generating embeddings...")
	for _, chunk := range uniqueChunks {

		fmt.Print(instruction.Query)
		fmt.Println(chunk)

		var vec []float64

		encoder := func(text string) error {

			// Run command and capture output
			//./embedding -m ../models/dolphin-phi2/dolphin-2_6-phi-2.Q8_0.gguf --log-disable -p

			// Run command using go exec
			cmd := exec.Command("/Users/arturoaquino/.eternal-v1/gguf/embedding", "-c", "4096", "-m", "/Users/arturoaquino/.eternal-v1/models/dolphin-phi2/dolphin-2_6-phi-2.Q8_0.gguf", "--log-disable", "-p", chunk)
			cmd.Stdin = strings.NewReader(text)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error running command:", err)
				return err
			}

			// Parse output
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "Vector:") {
					parts := strings.Split(line, " ")
					vec = make([]float64, len(parts)-1)
					for i, part := range parts[1:] {
						val, err := strconv.ParseFloat(part, 64)
						if err != nil {
							fmt.Println("Error parsing float:", err)
							return err
						}
						vec[i] = val
					}
				}
			}

			return nil
		}

		err = encoder(chunk) // Actually invoke the encoder function with the chunk
		if err != nil {
			pterm.Error.Println(err)
			panic(err)
		}

		embedding := estore.Embedding{
			Word:       chunk,
			Vector:     vec,
			Similarity: 0.0,
		}

		db.AddEmbedding(embedding)
	}

	// Save the database to a file
	pterm.Info.Println("Saving embeddings...")
	db.SaveEmbeddings("./embeddings.db")
}

func Search(dataPath string, dbName string, prompt string, topN int) []estore.Embedding {
	db := estore.NewEmbeddingDB()
	dbPath := fmt.Sprintf("%s/%s", dataPath, dbName)
	embeddings, err := db.LoadEmbeddings(dbPath)
	if err != nil {
		fmt.Println("Error loading embeddings:", err)
		return nil
	}

	// Run command and capture output
	//./embedding -m ../models/dolphin-phi2/dolphin-2_6-phi-2.Q8_0.gguf --log-disable -p

	// Run command using go exec
	cmd := exec.Command("/Users/arturoaquino/.eternal-v1/gguf/embedding", "-c", "4096", "-m", "/Users/arturoaquino/.eternal-v1/models/dolphin-phi2/dolphin-2_6-phi-2.Q8_0.gguf", "--log-disable", "-p", prompt)
	cmd.Stdin = strings.NewReader(prompt)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running command:", err)
		return nil
	}

	// Parse output
	lines := strings.Split(string(out), "\n")
	var vec []float64
	for _, line := range lines {
		if strings.HasPrefix(line, "Vector:") {
			parts := strings.Split(line, " ")
			vec = make([]float64, len(parts)-1)
			for i, part := range parts[1:] {
				val, err := strconv.ParseFloat(part, 64)
				if err != nil {
					fmt.Println("Error parsing float:", err)
					return nil
				}
				vec[i] = val
			}
		}
	}

	embeddingForPrompt := estore.Embedding{
		Word:       prompt,
		Vector:     vec,
		Similarity: 0.0,
	}

	// Retrieve the top N similar embeddings
	topEmbeddings := estore.FindTopNSimilarEmbeddings(embeddingForPrompt, embeddings, topN)
	if len(topEmbeddings) == 0 {
		fmt.Println("Error finding similar embeddings.")
		return nil
	}

	return topEmbeddings
}
