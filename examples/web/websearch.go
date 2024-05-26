package main

import (
	"eternal/pkg/web"
	"fmt"
	"log"

	"github.com/blevesearch/bleve/v2"
	index "github.com/blevesearch/bleve_index_api"
)

var (
	searchIndex bleve.Index
	webDocument string
)

// WebDocument is a struct that represents a web document
type WebDocument struct {
	URL     string
	Content string
}

func main() {
	// URL as cli input argument
	url := "https://blevesearch.com/docs/Character-Filters/"

	searchDB := "search.bleve"

	// Search the index, if no results are found, retrieve the document from the web
	searchIndex, err := bleve.Open(searchDB)
	if err != nil {
		log.Fatalf("Failed to open search index: %v", err)
	}

	// Search the index
	query := bleve.NewMatchQuery("character")
	search := bleve.NewSearchRequest(query)
	searchResults, err := searchIndex.Search(search)
	if err != nil {
		log.Fatal(err)
	}

	if searchResults.Total == 0 {
		// If the document is not in the index, retrieve it from the web
		webDocument, _ = web.WebGetHandler(url)
		searchIndexDocument(url, webDocument)
	} else {
		// Print the search results
		for _, hit := range searchResults.Hits {
			doc, err := searchIndex.Document(hit.ID)
			if err != nil {
				fmt.Errorf("Error retrieving document: %v", err)
				continue
			}
			doc.VisitFields(func(field index.Field) {
				fmt.Printf("%s: %s\n", field.Name(), field.Value())

				// Append the response field to the document
				if field.Name() == "content" {
					// Print the document content
					fmt.Println(field.Value())
				}
			})
		}
	}
}

func searchIndexDocument(url string, content string) {
	// Index the document
	err := searchIndex.Index(url, WebDocument{URL: url, Content: content})
	if err != nil {
		log.Fatal(err)
	}
}
