package documents

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

// GetPdfContents extracts the text content from the given PDF file and returns it as Markdown
func GetPdfContents(filePath string) (string, error) {
	// Open the PDF file
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()

	// Iterate through the pages
	totalPage := reader.NumPage()
	var sb strings.Builder
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		page := reader.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		// Extract text from the page
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", pageIndex, err)
		}

		// Format the text content of the page as Markdown
		markdownText := formatAsMarkdown(text)
		sb.WriteString(markdownText)
	}

	return sb.String(), nil
}

// formatAsMarkdown formats the given text as Markdown
func formatAsMarkdown(text string) string {
	var sb strings.Builder

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		// Detect headers (lines in all caps)
		if isHeader(trimmedLine) {
			sb.WriteString(fmt.Sprintf("## %s\n\n", trimmedLine))
		} else if isListItem(trimmedLine) {
			// Detect list items
			sb.WriteString(fmt.Sprintf("%s\n", trimmedLine))
		} else {
			sb.WriteString(fmt.Sprintf("%s\n\n", trimmedLine))
		}
	}

	return sb.String()
}

// isHeader checks if the given line is a header (i.e., all uppercase)
func isHeader(line string) bool {
	return strings.ToUpper(line) == line
}

// isListItem checks if the given line is a list item
var listItemPattern = regexp.MustCompile(`^\s*[\*\-\+] `)

func isListItem(line string) bool {
	return listItemPattern.MatchString(line)
}
