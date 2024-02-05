
# Documents Package

Documents is a package for processing and analyzing files and documents. It supports extraction, transformation, and organization of textual data for use in machine learning and AI workflows.

## Features

### Text Splitter
Provides a tool for recursively splitting text into smaller chunks based on a set of separators. It supports customization for different programming languages by adjusting the separators according to the syntax and structure of the language. This package can be especially useful for processing large text files, source code, or documents where specific splitting logic is required.

- **Recursive Splitting**: Splits text into smaller chunks recursively based on a set of separators.
- **Custom Separators**: Allows for custom separators and supports regular expressions for advanced splitting logic.
- **Language Support**: Predefined separators for various programming languages and document formats, including Python, Go, HTML, JavaScript (JS), TypeScript (TS), Markdown, and JSON.
- **Keep Separator**: Option to keep the separator as part of the returned chunks.
- **Custom Length Function**: Allows for a custom function to determine the chunk size, providing flexibility in how text is split.

### Git Repository Loader
Provides a tool for loading documents from a Git repository, including functionality for cloning repositories, checking out branches, and filtering files based on custom criteria. It is designed to integrate easily into Go projects requiring automatic fetching and processing of files from Git repositories.

- **Clone Git Repositories**: Clone a remote Git repository to a local path if it doesn't already exist.
- **Open Existing Repositories**: Open and use an existing local repository.
- **Branch Checkout**: Optionally checkout a specific branch of the repository.
- **Custom File Filtering**: Include or exclude files based on custom logic provided via a filter function.
- **SSH Authentication**: Authenticate to remote repositories using SSH private keys.
- **Insecure Host Key Verification Skip**: Option to skip SSH host key verification (use with caution).

## Example

```
package main

import (
    "fmt"
    "path/to/gitloader"
)

func main() {
    loader := gitloader.NewGitLoader("/path/to/local/repo", "git@github.com:user/repo.git", "main", "/path/to/private/key", nil, false)
    documents, err := loader.Load()
    if err != nil {
        fmt.Println("Error loading documents:", err)
        return
    }

    for _, doc := range documents {
        fmt.Println(doc.Metadata["file_name"], doc.PageContent)
    }
}

```