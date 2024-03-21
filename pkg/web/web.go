package web

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/pterm/pterm"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	unwantedURLs = []string{
		"youtube.com",
		"wired.com",
		"techcrunch.com",
		"wsj.com",
		"cnn.com",
		"nytimes.com",
		"forbes.com",
		"businessinsider.com",
		"theverge.com",
		"thehill.com",
		"theatlantic.com",
		"foxnews.com",
		"theguardian.com",
		"nbcnews.com",
		"msn.com",
		"sciencedaily.com",
		// Add more URLs to block from search results
	}
)

func WebGetHandler(url string) (string, error) {
	// Set up chromedp
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		// other options if needed...
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Retrieve and sanitize the page
	var docs string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &docs),
	)
	if err != nil {
		log.Println("Error retrieving page:", err)
		return "", err
	}

	// Clean up consecutive newlines
	docs = strings.ReplaceAll(docs, "\n\n", "\n")

	// Create a strings.Reader from docs to provide an io.Reader to ParseReaderView
	reader := strings.NewReader(docs)

	// Invoke web.ParseReaderView to get the reader view of the HTML
	cleanedHTML, err := ParseReaderView(reader)
	if err != nil {
		log.Println("Error parsing reader view:", err)
		return "", err
	}

	pterm.Info.Println("Document:", cleanedHTML)

	return cleanedHTML, nil
}

var (
	resultURLs []string
)

// ExtractURLs extracts and cleans URLs from the input string.
func ExtractURLs(input string) []string {
	// Regular expression to match URLs and port numbers
	urlRegex := `http.*?://[^\s<>{}|\\^` + "`" + `"]+`
	re := regexp.MustCompile(urlRegex)

	// Find all URLs in the input string
	matches := re.FindAllString(input, -1)

	var cleanedURLs []string
	for _, match := range matches {
		cleanedURL := cleanURL(match)
		cleanedURLs = append(cleanedURLs, cleanedURL)
	}

	return cleanedURLs
}

// cleanURL removes illegal trailing characters from the URL.
func cleanURL(url string) string {
	// Define illegal trailing characters.
	illegalTrailingChars := []rune{'.', ',', ';', '!', '?', ')'}

	for _, char := range illegalTrailingChars {
		if url[len(url)-1] == byte(char) {
			url = url[:len(url)-1]
		}
	}

	return url
}

func ParseReaderView(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	var readerView bytes.Buffer
	var f func(*html.Node)
	f = func(n *html.Node) {
		// Check if the node is an element node
		if n.Type == html.ElementNode {
			// Check if the node is a block element that typically contains content
			if isContentElement(n) {
				renderNode(&readerView, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return readerView.String(), nil
}

// isContentElement checks if the node is an element of interest for the reader view.
func isContentElement(n *html.Node) bool {
	switch n.DataAtom {
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.Article, atom.P, atom.Li, atom.Code, atom.Span, atom.Br:
		return true
	case atom.A, atom.Strong, atom.Em, atom.B, atom.I:
		return true
	}
	return false
}

// renderNode writes the content of the node to the buffer, including inline tags.
func renderNode(buf *bytes.Buffer, n *html.Node) {
	// Render the opening tag if it's not a text node
	if n.Type == html.ElementNode {
		buf.WriteString("<" + n.Data + ">")
	}

	// Render the contents
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			buf.WriteString(strings.TrimSpace(c.Data))
		} else if c.Type == html.ElementNode && isInlineElement(c) {
			// Render inline elements and their children
			renderNode(buf, c)
		}
	}

	// Render the closing tag if it's not a text node
	if n.Type == html.ElementNode {
		buf.WriteString("</" + n.Data + ">")
	}
}

// isInlineElement checks if the node is an inline element that should be included in the output.
func isInlineElement(n *html.Node) bool {
	switch n.DataAtom {
	case atom.A, atom.Strong, atom.Em, atom.B, atom.I, atom.Span, atom.Br, atom.Code:
		return true
	}
	return false
}

// SearchDuckDuckGo performs a search on DuckDuckGo and retrieves the HTML of the first page of results.
func SearchDuckDuckGo(query string) []string {
	// Initialize headless Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var nodes []*cdp.Node
	var resultURLs []string

	// Perform the search on DuckDuckGo
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://duckduckgo.com/`),
		chromedp.WaitVisible(`input[name="q"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="q"]`, query+kb.Enter, chromedp.ByQuery),
		chromedp.Sleep(5*time.Second), // Wait for JavaScript to load the search results
		chromedp.WaitVisible(`button[id="more-results"]`, chromedp.ByQuery),
		chromedp.Nodes(`a`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Printf("Error during search: %v", err)
		return nil
	}

	// Process the search results
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(c context.Context) error {
			re, err := regexp.Compile(`^http[s]?://`)
			if err != nil {
				return err
			}

			uniqueUrls := make(map[string]bool)
			for _, n := range nodes {
				for _, attr := range n.Attributes {
					if re.MatchString(attr) && !strings.Contains(attr, "duckduckgo") {
						uniqueUrls[attr] = true
					}
				}
			}

			for url := range uniqueUrls {
				resultURLs = append(resultURLs, url)
			}

			return nil
		}),
	)

	if err != nil {
		log.Printf("Error processing results: %v", err)
		return nil
	}

	// Remove unwanted URLs
	resultURLs = RemoveUnwantedURLs(resultURLs)

	// Only return the top n results
	resultURLs = resultURLs[:1]

	return resultURLs
}

// GetSearchResults loops over a list of URLs and retrieves the HTML of each page.
func GetSearchResults(urls []string) string {
	var resultHTML string

	for _, url := range urls {
		res, err := WebGetHandler(url)
		if err != nil {
			pterm.Error.Printf("Error getting search result: %v", err)
			continue
		}

		if res != "" {
			resultHTML += res
		}

		// time.Sleep(5 * time.Second)
	}

	return resultHTML
}

// RemoveUnwantedURLs removes unwanted URLs from the list of URLs.
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RemoveUnwantedURLs(urls []string) []string {
	var resultURLs []string
	for _, url := range urls {
		if !Contains(unwantedURLs, url) {
			resultURLs = append(resultURLs, url)
		}
	}

	return resultURLs
}

// GetPageScreenshot navigates to the given address and takes a full page screenshot.
func GetPageScreen(chromeUrl string, pageAddress string) string {

	instanceUrl := chromeUrl

	// Create allocator context for using existing Chrome instance
	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), instanceUrl)
	defer cancel()

	// Create context with logging for actions performed by chromedp
	ctx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Set a timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Run tasks
	var buf []byte // Buffer to store screenshot data
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageAddress),
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the URL to get the domain name
	u, err := url.Parse(pageAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Get the date and time
	t := time.Now()

	// Create a filename
	filename := u.Hostname() + "-" + t.Format("20060102150405") + ".png"

	// Save the screenshot to a file
	err = os.WriteFile(filename, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}
