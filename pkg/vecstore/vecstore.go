package store

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// Vector represents a vector of floats.
type Vector []float64

// Node is a struct that represents a node in a k-d tree. It has three fields:
// Domain: A slice of float64 values representing the domain or feature space of the data point associated with this node.
// Value: A float64 value representing the pivot value used to partition the data points into two subsets.
// Left: A pointer to the left child node in the tree, or nil if there is no left child.
// Right: A pointer to the right child node in the tree, or nil if there is no right child.
type Node struct {
	Domain []float64
	Value  float64
	Left   *Node
	Right  *Node
}

// Embedding represents a word embedding.
type Embedding struct {
	Word       string
	Vector     []float64
	Similarity float64 // Similarity field to store the cosine similarity
}

// EmbeddingDB represents a database of Embeddings.
type EmbeddingDB struct {
	Embeddings map[string]Embedding
}

// Document represents a document to be ranked.
type Document struct {
	ID     string
	Score  float64
	Length int
}

// NewEmbeddingDB creates a new embedding database.
func NewEmbeddingDB() *EmbeddingDB {
	return &EmbeddingDB{
		Embeddings: make(map[string]Embedding),
	}
}

// AddEmbedding adds a new embedding to the database.
func (db *EmbeddingDB) AddEmbedding(embedding Embedding) {
	db.Embeddings[embedding.Word] = embedding
}

// AddEmbeddings adds a slice of embeddings to the database.
func (db *EmbeddingDB) AddEmbeddings(embeddings []Embedding) {
	for _, embedding := range embeddings {
		db.AddEmbedding(embedding)
	}
}

// SaveEmbeddings saves the Embeddings to a file, appending new ones to existing data.
func (db *EmbeddingDB) SaveEmbeddings(path string) error {
	// Read the existing content from the file
	var existingEmbeddings map[string]Embedding
	content, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error reading file: %v", err)
		}
		existingEmbeddings = make(map[string]Embedding)
	} else {
		err = json.Unmarshal(content, &existingEmbeddings)
		if err != nil {
			return fmt.Errorf("error unmarshaling existing embeddings: %v", err)
		}
	}

	// Merge new embeddings with existing ones
	for key, embedding := range db.Embeddings {
		existingEmbeddings[key] = embedding
	}

	// Marshal the combined embeddings to JSON
	jsonData, err := json.Marshal(existingEmbeddings)
	if err != nil {
		return fmt.Errorf("error marshaling embeddings: %v", err)
	}

	// Open the file in write mode (this will overwrite the existing file)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Write the combined JSON to the file
	if _, err := f.Write(jsonData); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

// LoadEmbeddings loads the Embeddings from a file.
func (db *EmbeddingDB) LoadEmbeddings(path string) (map[string]Embedding, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var embeddings map[string]Embedding
	err = json.Unmarshal(content, &embeddings)
	if err != nil {
		return nil, err
	}

	return embeddings, nil
}

// RetrieveEmbedding retrieves an embedding from the database.
func (db *EmbeddingDB) RetrieveEmbedding(word string) ([]float64, bool) {
	embedding, exists := db.Embeddings[word]
	if !exists {
		return nil, false
	}

	return embedding.Vector, true
}

// RecreateDocument recreates a document from a slice of embeddings.
func (db *EmbeddingDB) RecreateDocument(embeddings []Embedding) string {
	var document []string
	for _, embedding := range embeddings {
		document = append(document, embedding.Word)
	}

	return strings.Join(document, " ")
}

// CosineSimilarity calculates the cosine similarity between two vectors.
func CosineSimilarity(a, b []float64) float64 {
	var dotProduct, magnitudeA, magnitudeB float64
	var wg sync.WaitGroup

	// Adjust the number of partitions based on the number of CPU cores.
	partitions := runtime.NumCPU()
	partSize := len(a) / partitions

	results := make([]struct {
		dotProduct, magnitudeA, magnitudeB float64
	}, partitions)

	for i := 0; i < partitions; i++ {
		wg.Add(1)
		go func(partition int) {
			defer wg.Done()
			start := partition * partSize
			end := start + partSize
			if partition == partitions-1 {
				end = len(a)
			}
			for j := start; j < end; j++ {
				results[partition].dotProduct += a[j] * b[j]
				results[partition].magnitudeA += a[j] * a[j]
				results[partition].magnitudeB += b[j] * b[j]
			}
		}(i)
	}

	wg.Wait()

	for _, result := range results {
		dotProduct += result.dotProduct
		magnitudeA += result.magnitudeA
		magnitudeB += result.magnitudeB
	}

	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
}

// MostSimilarWord returns the word with the highest similarity value.
func (db *EmbeddingDB) MostSimilarWord(embeddings map[string]Embedding, targetWord string) (string, float64) {
	// Check for an exact match
	if _, exists := embeddings[targetWord]; exists {
		return targetWord, 1.0
	}
	var targetVector []float64

	// If the target word exists in the embeddings, use its vector.
	if embedding, exists := embeddings[targetWord]; exists {
		targetVector = embedding.Vector
	} else {
		// If target word doesn't exist, print the error and try to find the most similar word using the embeddings available.
		// (In a more robust implementation, you might want to obtain a vector for the targetWord from another source.)
		fmt.Printf("Error: Word '%s' not found in embeddings database.\n", targetWord)
		return "", -1.0 // We'll use this -1.0 later to identify that the target word wasn't in the database.
	}

	mostSimilarWord := ""
	highestSimilarity := -2.0 // Starting with -2 to ensure that any cosine similarity will be higher.

	for word, embedding := range embeddings {
		// If the word is the same as target, skip.
		if word == targetWord {
			continue
		}

		// Compute similarity.
		similarity := CosineSimilarity(targetVector, embedding.Vector)

		// If this word's similarity is greater than the highest similarity seen so far, update.
		if similarity > highestSimilarity {
			mostSimilarWord = word
			highestSimilarity = similarity
		}
	}

	if highestSimilarity == -1.0 {
		return "No similar words found in database", highestSimilarity
	}

	return mostSimilarWord, highestSimilarity
}

// TODO: Implement an Efficient Search Mechanism using KD-Trees or Ball Trees.
// Consider integrating with a library or service that offers efficient nearest-neighbor search capabilities.
// Placeholder function for this:
func EfficientSearch(targetWord string) (string, float64) {
	// Implement efficient search here
	return "", 0.0
}

// dataPoints: A slice of slices of float64 values representing the data points to be inserted into the k-d tree. Each inner slice should have the same length, which represents the number of dimensions or features of the data point.
// dimensions: An integer representing the number of dimensions or features of the data points.
func BuildKdTree(dataPoints [][]float64, dimensions int) *Node {
	if len(dataPoints) == 0 {
		return nil
	}

	// Use a stack to store the nodes in the tree.
	stack := make([]*Node, 0)

	// Select the median value of a random dimension as the pivot.
	pivotIndex := rand.Intn(len(dataPoints))
	pivotValue := dataPoints[pivotIndex][rand.Intn(dimensions)]

	// Partition the data points into two subsets based on the pivot value.
	left := make([][]float64, 0)
	right := make([][]float64, 0)
	for _, point := range dataPoints {
		if point[rand.Intn(dimensions)] < pivotValue {
			left = append(left, point)
		} else {
			right = append(right, point)
		}
	}

	// Build the left and right subtrees iteratively using a stack.
	node := &Node{
		Domain: dataPoints[pivotIndex],
		Value:  pivotValue,
	}
	stack = append(stack, node)
	for len(left) > 0 || len(right) > 0 {
		if len(left) > 0 {
			// Build the left subtree.
			pivotIndex := rand.Intn(len(left))
			pivotValue := left[pivotIndex][rand.Intn(dimensions)]
			var partition [][]float64
			for _, point := range left {
				if point[rand.Intn(dimensions)] < pivotValue {
					partition = append(partition, point)
				}
			}
			node := &Node{
				Domain: left[pivotIndex],
				Value:  pivotValue,
			}
			stack[len(stack)-1].Left = node
			if len(partition) > 0 {
				stack = append(stack, node)
				left = partition
			} else {
				left = nil
			}
		} else {
			// Build the right subtree.
			pivotIndex := rand.Intn(len(right))
			pivotValue := right[pivotIndex][rand.Intn(dimensions)]
			var partition [][]float64
			for _, point := range right {
				if point[rand.Intn(dimensions)] < pivotValue {
					partition = append(partition, point)
				}
			}
			node := &Node{
				Domain: right[pivotIndex],
				Value:  pivotValue,
			}
			stack[len(stack)-1].Right = node
			if len(partition) > 0 {
				stack = append(stack, node)
				right = partition
			} else {
				right = nil
			}
		}
	}

	// Return the root node.
	return stack[0]
}

// FindMostSimilarEmbedding finds the most similar embeddings in the database.
func FindMostSimilarEmbedding(targetEmbedding Embedding, embeddings map[string]Embedding) (Embedding, bool) {
	var mostSimilarEmbedding Embedding
	var highestSimilarity float64

	for _, embedding := range embeddings {
		// If the word is the same as target, skip.
		if embedding.Word == targetEmbedding.Word {
			continue
		}

		// Compute similarity.
		similarity := CosineSimilarity(targetEmbedding.Vector, embedding.Vector)

		// If this word's similarity is greater than the highest similarity seen so far, update.
		if similarity > highestSimilarity {
			mostSimilarEmbedding = embedding
			highestSimilarity = similarity
		}
	}

	if highestSimilarity == -1.0 {
		return Embedding{}, false
	}

	return mostSimilarEmbedding, true
}

// NormalizeL2 normalizes a vector using L2 normalization.
func NormalizeL2(vec []float64) []float64 {
	var sumSquares float64
	for _, value := range vec {
		sumSquares += value * value
	}
	norm := math.Sqrt(sumSquares)
	for i, value := range vec {
		vec[i] = value / norm
	}
	return vec
}

// ComputeSimilarityMatrix computes the cosine similarity matrix between two slices of embeddings.
func ComputeSimilarityMatrix(queryEmbeddings, keyEmbeddings []Embedding) [][]float64 {
	matrix := make([][]float64, len(queryEmbeddings))
	for i, query := range queryEmbeddings {
		matrix[i] = make([]float64, len(keyEmbeddings))
		for j, key := range keyEmbeddings {
			matrix[i][j] = CosineSimilarity(query.Vector, key.Vector)
		}
	}
	return matrix
}

// SimilarityWithKey is a type that holds both the similarity value and the corresponding word key.
type SimilarityWithKey struct {
	Similarity float64
	Key        string
}

// FindTopNSimilarEmbeddings finds the top N most similar embeddings in the database.
func FindTopNSimilarEmbeddings(targetEmbedding Embedding, embeddings map[string]Embedding, topN int) []Embedding {
	var topEmbeddings []Embedding
	var similarityList []SimilarityWithKey

	// Compute the cosine similarity for each embedding in the database and store it with its key.
	for key, embedding := range embeddings {
		similarity := CosineSimilarity(targetEmbedding.Vector, embedding.Vector)
		similarityList = append(similarityList, SimilarityWithKey{similarity, key})
	}

	// Sort the similarityList in descending order of similarity.
	sort.SliceStable(similarityList, func(i, j int) bool {
		return similarityList[i].Similarity > similarityList[j].Similarity
	})

	// Retrieve the top N most similar embeddings.
	for i := 0; i < topN && i < len(similarityList); i++ {
		topEmbeddings = append(topEmbeddings, embeddings[similarityList[i].Key])
	}

	return topEmbeddings
}

// Reranker function reranks documents based on a weighted combination of score and length.
func Reranker(documents []Document, weightScore float64, weightLength float64) []Document {
	// Validate weights
	if weightScore < 0 || weightLength < 0 || (weightScore+weightLength) == 0 {
		// Handle invalid weights
		return documents
	}

	rerankedDocuments := make([]Document, len(documents))
	copy(rerankedDocuments, documents)

	sort.SliceStable(rerankedDocuments, func(i, j int) bool {
		scoreDiffI := rerankedDocuments[i].Score * weightScore
		lengthDiffI := float64(rerankedDocuments[i].Length) * weightLength
		combinedScoreI := scoreDiffI + lengthDiffI

		scoreDiffJ := rerankedDocuments[j].Score * weightScore
		lengthDiffJ := float64(rerankedDocuments[j].Length) * weightLength
		combinedScoreJ := scoreDiffJ + lengthDiffJ

		if combinedScoreI == combinedScoreJ {
			// Handle tie-breaking here if needed
		}

		return combinedScoreI > combinedScoreJ
	})

	return rerankedDocuments
}
