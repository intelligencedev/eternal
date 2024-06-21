package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "No unclosed code blocks",
			input:    []byte("This is a test\n```\ncode\n```"),
			expected: []byte("This is a test\n```\ncode\n```"),
		},
		{
			name:     "One unclosed code block",
			input:    []byte("This is a test\n```\ncode"),
			expected: []byte("This is a test\n```\ncode\n```"),
		},
		{
			name:     "Multiple unclosed code blocks",
			input:    []byte("```\ncode1\n```\ntext\n```\ncode2"),
			expected: []byte("```\ncode1\n```\ntext\n```\ncode2\n```"),
		},
		{
			name:     "Empty input",
			input:    []byte(""),
			expected: []byte(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PreprocessMarkdown(tt.input)
			assert.Equal(t, tt.expected, result, "PreprocessMarkdown() result does not match expected output")
		})
	}
}

func TestMarkdownToHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Basic markdown",
			input:    []byte("# Hello\nThis is a test"),
			expected: "<h1>Hello</h1>\n\n<p>This is a test</p>\n",
		},
		{
			name:     "Markdown with code block",
			input:    []byte("```\ncode\n```"),
			expected: "<pre><code>code\n</code></pre>\n",
		},
		{
			name:     "Markdown with unclosed code block",
			input:    []byte("```\ncode"),
			expected: "<pre><code>code\n</code></pre>\n",
		},
		{
			name:     "Empty input",
			input:    []byte(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarkdownToHTML(tt.input)
			assert.Contains(t, string(result), tt.expected, "MarkdownToHTML() result does not contain expected output")
		})
	}
}
