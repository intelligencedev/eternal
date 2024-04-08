package embeddings

import (
	"context"
	"eternal/pkg/documents"
	"fmt"
	"os"
	"strings"

	estore "eternal/pkg/vecstore"

	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
	"github.com/pterm/pterm"
)

// var modelName = "BAAI/bge-large-en-v1.5"
var modelName = "avsolatorio/GIST-small-Embedding-v0"

const limit = 128

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

func GenerateEmbeddingForTask(task string, doctype string, chunkSize int, dataPath string) {

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

	var chunks []string
	var separators []string

	// NOTE: If using txt doctype, the size must equal the ChunkSize in RecursiveCharacterTextSplitter
	if doctype == "txt" {
		chunks = documents.SplitTextByCount(string(content), chunkSize)
	} else {
		// Convert doctype to uppercase
		doctype = strings.ToUpper(doctype)

		separators, _ = documents.GetSeparatorsForLanguage(documents.Language(doctype))

		// Updated the RecursiveCharacterTextSplitter to include OverlapSize and updated SplitText method
		splitter := documents.RecursiveCharacterTextSplitter{
			Separators:       separators,
			KeepSeparator:    true,
			IsSeparatorRegex: false,
			ChunkSize:        chunkSize,
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

	modelsDir := fmt.Sprintf("%s/data/models/HF/%s/", dataPath, modelName)

	model, err := tasks.Load[textencoding.Interface](&tasks.Config{ModelsDir: modelsDir, ModelName: modelName})
	if err != nil {
		pterm.Error.Println("Error loading model...")
		panic(err)
	}

	// 3. Embedding Generation
	pterm.Info.Println("Generating embeddings...")
	for _, chunk := range uniqueChunks {

		fmt.Print(instruction.Query)
		fmt.Println(chunk)

		var vec []float64

		encoder := func(text string) error {
			result, err := model.Encode(context.Background(), text, int(bert.MeanPooling))
			if err != nil {
				return err
			}

			vec = result.Vector.Data().F64()[:limit]
			fmt.Println(result.Vector.Data().F64()[:limit])
			return nil
		}

		err = encoder(chunk) // Actually invoke the encoder function with the chunk
		if err != nil {
			pterm.Error.Println("Error encoding text...")
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

// GenerateEmbedding generates an embedding from a prompt.
func GenerateEmbedding(dataPath string) {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <path_to_input_file>")
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

	modelsDir := fmt.Sprintf("%s/data/models/HF/%s/", dataPath, modelName)

	model, err := tasks.Load[textencoding.Interface](&tasks.Config{ModelsDir: modelsDir, ModelName: modelName})
	if err != nil {
		pterm.Error.Println("Error loading model...")
		panic(err)
	}

	// 3. Embedding Generation
	pterm.Info.Println("Generating embeddings...")
	for _, chunk := range chunks {

		fmt.Print("Encoding text...")
		fmt.Println(chunk)

		var vec []float64
		err := func(text string) error {
			result, err := model.Encode(context.Background(), text, int(bert.MeanPooling))
			if err != nil {
				return err
			}

			vec = result.Vector.Data().F64()[:limit]

			embedding := estore.Embedding{
				Word:       chunk,
				Vector:     vec,
				Similarity: 0.0,
			}

			db.AddEmbedding(embedding)

			fmt.Println(result.Vector.Data().F64()[:limit])

			return nil
		}(chunk) // Actually invoke the encoder function with the chunk

		if err != nil {
			pterm.Error.Println("Error encoding text...")
			panic(err)
		}
	}

	// Save the database to a file
	pterm.Info.Println("Saving embeddings...")
	db.SaveEmbeddings("./test.db")

	if len(chunks) > 0 {
		embedding, ok := db.RetrieveEmbedding(chunks[0])
		if ok {
			fmt.Printf("Embedding for the first chunk:\n%v\n", embedding)
		}
	}
}

// GenerateEmbeddingChat generates an embedding from a prompt for chatbot applications
func GenerateEmbeddingChat(prompt string, dataPath string) {

	db := estore.NewEmbeddingDB()

	modelsDir := fmt.Sprintf("%s/data/models/HF/", dataPath)

	model, err := tasks.Load[textencoding.Interface](&tasks.Config{ModelsDir: modelsDir, ModelName: modelName})
	if err != nil {
		pterm.Error.Println("Error loading model...")
		panic(err)
	}

	pterm.Info.Println("Generating embeddings...")

	var embeddings []estore.Embedding // Slice to store multiple embeddings

	encoder := func(text string) error {
		result, err := model.Encode(context.Background(), text, int(bert.MeanPooling))
		if err != nil {
			return err
		}

		vec := result.Vector.Data().F64()[:limit] // Ensure 'limit' is defined and valid
		fmt.Println(vec)

		embedding := estore.Embedding{
			Word:       text,
			Vector:     vec,
			Similarity: 0.0,
		}

		embeddings = append(embeddings, embedding) // Append to slice
		return nil
	}

	err = encoder(prompt)
	if err != nil {
		pterm.Error.Println("Error encoding text...")
		panic(err)
	}

	// Add generated embeddings to the database
	db.AddEmbeddings(embeddings)

	// Now save all embeddings
	pterm.Info.Println("Saving embeddings...")
	dbPath := fmt.Sprintf("%s/responses.db", dataPath)
	pterm.Info.Println(dbPath)

	if err := db.SaveEmbeddings(dbPath); err != nil {
		pterm.Error.Println("Error saving embeddings:", err)
		panic(err)
	}
}

func Search(dataPath string, dbName string, prompt string, topN int) []estore.Embedding {
	db := estore.NewEmbeddingDB()
	dbPath := fmt.Sprintf("%s/%s", dataPath, dbName)
	embeddings, err := db.LoadEmbeddings(dbPath)
	if err != nil {
		fmt.Println("Error loading embeddings:", err)
		return nil
	}

	embeddingsModelPath := fmt.Sprintf("%s/data/models/HF/", dataPath)

	model, err := tasks.Load[textencoding.Interface](&tasks.Config{
		ModelsDir:           embeddingsModelPath,
		ModelName:           modelName,
		DownloadPolicy:      tasks.DownloadMissing,
		ConversionPolicy:    tasks.ConvertMissing,
		ConversionPrecision: tasks.F32,
	})

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
