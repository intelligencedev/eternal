package embeddings

import (
	"context"
	"eternal/pkg/documents"
	"fmt"
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

func GenerateEmbeddingForTask(task string, content string, doctype string, chunkSize int, overlapSize int, dataPath string) error {

	_, ok := INSTRUCTIONS[task]
	if !ok {
		fmt.Printf("Unknown task: %s\n", task)
		return fmt.Errorf("unknown task: %s", task)
	}

	db := estore.NewEmbeddingDB()

	var chunks []string
	var separators []string

	if doctype == "txt" {
		// convert to lower case
		content = strings.ToLower(content)
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

	modelsDir := fmt.Sprintf("%s/data/models/HF/%s/", dataPath, modelName)

	tasksConfig := &tasks.Config{
		ModelsDir:        modelsDir,
		ModelName:        modelName,
		DownloadPolicy:   tasks.DownloadMissing,
		ConversionPolicy: tasks.ConvertMissing,
	}

	model, err := tasks.Load[textencoding.Interface](tasksConfig)
	if err != nil {
		pterm.Error.Println("Error loading model...")
		return err
	}

	// 3. Embedding Generation
	pterm.Info.Println("Generating embeddings...")
	for _, chunk := range uniqueChunks {
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
			return err
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

	dbPath := fmt.Sprintf("%s/embeddings.db", dataPath)

	db.SaveEmbeddings(dbPath)

	return nil
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
