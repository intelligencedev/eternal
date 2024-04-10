package documents

import (
	"errors"
	"regexp"
)

// RecursiveCharacterTextSplitter is a struct that represents a text splitter
// that splits text based on recursive character separators.
type RecursiveCharacterTextSplitter struct {
	Separators       []string
	KeepSeparator    bool
	IsSeparatorRegex bool
	ChunkSize        int
	OverlapSize      int
	LengthFunction   func(string) int
}

// Language is a type that represents a programming language.
type Language string

const (
	PYTHON   Language = "PYTHON"
	GO       Language = "GO"
	HTML     Language = "HTML"
	JS       Language = "JS"
	TS       Language = "TS"
	MARKDOWN Language = "MARKDOWN"
	JSON     Language = "JSON"
)

// escapeString is a helper function that escapes special characters in a string.
func escapeString(s string) string {
	return regexp.QuoteMeta(s)
}

// splitTextWithRegex is a helper function that splits text using a regular expression separator.
func splitTextWithRegex(text string, separator string, keepSeparator bool) []string {
	sepPattern := regexp.MustCompile(separator)
	splits := sepPattern.Split(text, -1)
	if keepSeparator {
		matches := sepPattern.FindAllString(text, -1)
		result := make([]string, 0, len(splits)+len(matches))
		for i, split := range splits {
			result = append(result, split)
			if i < len(matches) {
				result = append(result, matches[i])
			}
		}
		return result
	}
	return splits
}

// SplitTextByCount splits the given text into chunks of the given size.
func SplitTextByCount(text string, size int) []string {
	// slice the string into chunks of size
	var chunks []string
	for i := 0; i < len(text); i += size {
		end := i + size
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}
	return chunks
}

// SplitText splits the given text using the configured separators.
func (r *RecursiveCharacterTextSplitter) SplitText(text string) []string {
	chunks := r.splitTextHelper(text, r.Separators)

	// Apply chunk overlap
	if r.OverlapSize > 0 {
		overlappedChunks := make([]string, 0)
		for i := 0; i < len(chunks)-1; i++ {
			currentChunk := chunks[i]
			nextChunk := chunks[i+1]

			nextChunkOverlap := nextChunk[:min(len(nextChunk), r.OverlapSize)]

			overlappedChunk := currentChunk + nextChunkOverlap
			overlappedChunks = append(overlappedChunks, overlappedChunk)
		}
		overlappedChunks = append(overlappedChunks, chunks[len(chunks)-1])

		chunks = overlappedChunks
	}

	return chunks
}

// splitTextHelper is a recursive helper function that splits text using the given separators.
func (r *RecursiveCharacterTextSplitter) splitTextHelper(text string, separators []string) []string {
	finalChunks := make([]string, 0)

	if len(separators) == 0 {
		return []string{text}
	}

	// Determine the separator
	separator := separators[len(separators)-1]
	newSeparators := make([]string, 0)
	for i, sep := range separators {
		sepPattern := sep
		if !r.IsSeparatorRegex {
			sepPattern = escapeString(sep)
		}
		if regexp.MustCompile(sepPattern).MatchString(text) {
			separator = sep
			newSeparators = separators[i+1:]
			break
		}
	}

	// Split the text using the determined separator
	splits := splitTextWithRegex(text, separator, r.KeepSeparator)

	// Check each split
	for _, s := range splits {
		if r.LengthFunction(s) < r.ChunkSize {
			finalChunks = append(finalChunks, s)
		} else if len(newSeparators) > 0 {
			// If the split is too large, try to split it further using remaining separators
			recursiveSplits := r.splitTextHelper(s, newSeparators)
			finalChunks = append(finalChunks, recursiveSplits...)
		} else {
			// If no more separators left, add the large chunk as it is
			finalChunks = append(finalChunks, s)
		}
	}

	return finalChunks
}

// FromLanguage creates a RecursiveCharacterTextSplitter based on the given language.
func FromLanguage(language Language) (*RecursiveCharacterTextSplitter, error) {
	separators, err := GetSeparatorsForLanguage(language)
	if err != nil {
		return nil, err
	}
	return &RecursiveCharacterTextSplitter{
		Separators:       separators,
		IsSeparatorRegex: true,
	}, nil
}

// GetSeparatorsForLanguage returns the separators for the given language.
func GetSeparatorsForLanguage(language Language) ([]string, error) {
	switch language {
	case PYTHON:
		return []string{
			// Split along class definitions
			"\nclass ",
			"\ndef ",
			"\n\tdef ",
			// Split by the normal type of lines
			"\n\n",
			"\n",
			" ",
			"",
		}, nil
	case GO:
		return []string{
			// Split along function definitions
			"\nfunc ",
			"\nvar ",
			"\nconst ",
			"\ntype ",
			// Split along control flow statements
			"\nif ",
			"\nfor ",
			"\nswitch ",
			"\ncase ",
			// Split by the normal type of lines
			"\n\n",
			"\n",
			" ",
			"",
		}, nil
	case HTML:
		return []string{
			// Split along HTML tags
			"<body",
			"<div",
			"<p",
			"<br",
			"<li",
			"<h1",
			"<h2",
			"<h3",
			"<h4",
			"<h5",
			"<h6",
			"<span",
			"<table",
			"<tr",
			"<td",
			"<th",
			"<ul",
			"<ol",
			"<header",
			"<footer",
			"<nav",
			// Head
			"<head",
			"<style",
			"<script",
			"<meta",
			"<title",
			"",
			"\n</",
		}, nil
	case JS:
		return []string{
			// Split along function definitions
			"\nfunction ",
			"\nconst ",
			"\nlet ",
			"\nvar ",
			"\nclass ",
			// Split along control flow statements
			"\nif ",
			"\nfor ",
			"\nwhile ",
			"\nswitch ",
			"\ncase ",
			"\ndefault ",
			// Split by the normal type of lines
			"\n\n",
			"\n",
			" ",
			"",
		}, nil
	case TS:
		return []string{
			"\nenum ",
			"\ninterface ",
			"\nnamespace ",
			"\ntype ",
			// Split along class definitions
			"\nclass ",
			// Split along function definitions
			"\nfunction ",
			"\nconst ",
			"\nlet ",
			"\nvar ",
			// Split along control flow statements
			"\nif ",
			"\nfor ",
			"\nwhile ",
			"\nswitch ",
			"\ncase ",
			"\ndefault ",
			// Split by the normal type of lines
			"\n\n",
			"\n",
			" ",
			"",
		}, nil
	case MARKDOWN:
		return []string{
			// First, try to split along Markdown headings (starting with level 2)
			"\n#{1,6} ",
			// Note the alternative syntax for headings (below) is not handled here
			// Heading level 2
			// ---------------
			// End of code block
			"```\n",
			// Horizontal lines
			"\n\\*\\*\\*+\n",
			"\n---+\n",
			"\n___+\n",
			// Note that this splitter doesn't handle horizontal lines defined
			// by *three or more* of ***, ---, or ___, but this is not handled
			"\n\n",
			"\n",
			" ",
			"",
		}, nil
	case JSON:
		return []string{
			"}\n",
		}, nil
	default:
		return nil, errors.New("unsupported language")
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
