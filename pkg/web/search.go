package web

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pterm/pterm"
	"golang.org/x/net/html"
)

// postRequest sends a POST request to the given endpoint with a named parameter 'q' and returns the response body as a string.
func postRequest(endpoint string, queryParam string) (string, error) {
	// Create the form data
	formData := url.Values{}
	formData.Set("q", queryParam)

	// Convert form data to a byte buffer
	data := bytes.NewBufferString(formData.Encode())

	// Create a new POST request
	req, err := http.NewRequest("POST", endpoint, data)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

// extractURLs parses the HTML response and extracts the URLs from the search results.
func extractURLs(htmlContent string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var urls []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.Contains(attr.Val, "http") {
					urls = append(urls, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return urls, nil
}

func GetSearXNGResults(endpoint string, query string) []string {
	htmlContent, err := postRequest(endpoint, query)
	if err != nil {
		pterm.Error.Printf("Error: %v\n", err)
		return nil
	}

	urls, err := extractURLs(htmlContent)
	if err != nil {
		pterm.Error.Printf("Error extracting URLs: %v\n", err)
		return nil
	}

	// Remove unwanted URLs
	urls = RemoveUnwantedURLs(urls)

	return urls
}
