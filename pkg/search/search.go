package search

import (
	"github.com/blevesearch/bleve/v2"
)

// InitIndex initializes or opens a Bleve index at the given path.
func InitIndex(indexPath string) (*bleve.Index, error) {
	mapping := bleve.NewIndexMapping() // Default mapping; can be customized as needed
	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// IndexData indexes the given data with the specified ID.
func IndexData(index *bleve.Index, id string, data interface{}) error {
	return (*index).Index(id, data)
}

// Search performs a search query on the given index.
func Search(index *bleve.Index, query string) (*bleve.SearchResult, error) {
	searchQuery := bleve.NewMatchQuery(query)
	search := bleve.NewSearchRequest(searchQuery)
	return (*index).Search(search)
}
