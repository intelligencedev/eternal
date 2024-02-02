# Hugging Face API Client

NOTE: This is not complete. It is intended as an example of how the hf package could be extended to support direct model loading/usage.

The `hf` package interacts with the Hugging Face API to retrieve model information.

## Features

- Fetches model data from the Hugging Face API.
- Structured representation of model data (`ModelData`).

## Usage

### Fetching Model Data

To fetch model data from the Hugging Face API:

```go
models, err := hf.GetHFModels("search_query")
if err != nil {
// handle error
}
// use models
```

You can pass a specific search query to the GetHFModels function, such as the name of a user or organization.

### ModelData Structure

The ModelData struct is designed to capture key details about each model:

```go
type ModelData struct {
ID          string
ModelID     string
Likes       int
IsPrivate   bool
Downloads   int
Tags        []string
PipelineTag string
LibraryName string
CreatedAt   time.Time
}
```

Each field corresponds to a piece of data returned from the Hugging Face API.
