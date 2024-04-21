package web

import (
	"bytes"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// PreprocessMarkdown scans the markdown content for unclosed code blocks and closes them.
func PreprocessMarkdown(content []byte) []byte {
	// Simple heuristic: if the count of ``` is odd, append one at the end
	if bytes.Count(content, []byte("```"))%2 != 0 {
		content = append(content, []byte("\n```")...) // Append closing code block
	}
	return content
}

// MarkdownToHTML converts preprocessed markdown content to HTML.
func MarkdownToHTML(mdContent []byte) []byte {
	// Preprocess to ensure code blocks are properly closed
	preprocessedContent := PreprocessMarkdown(mdContent)

	// Setup parser and renderer
	extensions := parser.CommonExtensions
	parser := parser.NewWithExtensions(extensions)
	htmlFlags := html.CommonFlags

	// Custom CSS to apply 0px bottom margin to <code> elements
	customCSS := "code { margin-bottom: 0px; }"

	renderer := html.NewRenderer(html.RendererOptions{
		Flags: htmlFlags,
		CSS:   customCSS, // Adding custom CSS here
	})

	// Convert markdown to HTML
	htmlContent := markdown.ToHTML(preprocessedContent, parser, renderer)

	return htmlContent
}
