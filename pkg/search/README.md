# Search Package

NOTE: This is NOT implemented yet, just an example of how the search package could work.

The `search` package implements full-text search using the Bleve search package.

## Features

- **Initialize an Index**: Easily create or open a Bleve index.
- **Index Data**: Add data to the Bleve index for future searches.
- **Search Data**: Perform full-text search queries on the indexed data.

## Installation

Before using the `search` package, ensure you have Bleve installed:

```bash
go get github.com/blevesearch/bleve/v2
```

## Usage

### Initializing the Index

To initialize a new index or open an existing one:

```go
index, err := search.InitIndex("path/to/index")
if err != nil {
    // Handle error
}
```

### Indexing Data

To add data to the index:

```go
err := search.IndexData(index, "uniqueID", dataObject)
if err != nil {
    // Handle error
}
```

`dataObject` can be any struct or map that you want to index. Ensure it's structured in a way that Bleve can understand for indexing.

### Searching Data

To search the indexed data:

```go
searchResults, err := search.Search(index, "search query")
if err != nil {
    // Handle error
}
// Process searchResults
```

The search function supports a variety of query formats and options. This basic example uses a match query, which finds documents that match a specified text.

### Customization

The search package uses Bleve's default mapping for indexes, but you can customize this mapping according to your specific needs. Refer to Bleve's documentation for more details on custom mappings.

### Dependencies

This package is built on top of Bleve v2, which is a powerful full-text search and indexing library for Go.
For more details on Bleve and its functionalities the following resources can be referenced:
- [Bleve Github](https://github.com/blevesearch/bleve)
- [Bleve Documentation](https://blevesearch.com/docs/Home/)