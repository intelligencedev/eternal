package main

import (
	"eternal/pkg/web" // Replace with the actual path to the package
	"fmt"
	"os"
)

var pages string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <search-query>")
		os.Exit(1)
	}

	query := os.Args[1]

	// Perform a search on DuckDuckGo with the provided query
	results := web.SearchDuckDuckGo(query)

	// Check if the search returned any results
	if results == nil {
		fmt.Println("No results found or there was an error during the search")
		os.Exit(1)
	}

	// Print out the URLs obtained from the search
	for _, url := range results {
		fmt.Println(url)
	}

	// Extract content from the urls
	pages = web.GetSearchResults(results)

	fmt.Println(pages)
}
