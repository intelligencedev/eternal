package web

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

// Run headless Chrome instance with remote debugging enabled:
// Linux: $ google-chrome --headless --remote-debugging-port=9222 --disable-gpu --no-sandbox --headless --disable-dev-shm-usage
// MacOS: $ /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --remote-debugging-port=9222 --headless
//
// EXAMPLE USAGE -- Get search results from DuckDuckGo
// Dump pages into JSON document
//
// func main() {
// 	searchQuery := flag.String("query", "Top AI news today", "The search query")
// 	remoteURL := flag.String("url", "ws://0.0.0.0:9222/devtools/browser/adc8c3e9-9e5f-...-...", "URL to the Chrome DevTools instance")

// 	flag.Parse()

// 	// Use *searchQuery and *remoteURL to access the string values
// 	fmt.Println("Search Query:", *searchQuery)
// 	fmt.Println("Remote URL:", *remoteURL)

// 	// Get search results from DuckDuckGo
// 	results := SearchDuckDuckGo(*remoteURL, *searchQuery)

// 	webDocuments := make(map[string]string)

// 	for _, pageURL := range results {
// 		pageContent, _ := RetrieveAndSanitizePage(*remoteURL, pageURL)
// 		webDocuments[pageURL] = pageContent

// 		fmt.Println(pageURL)
// 		fmt.Println(pageContent)
// 	}

// 	// Write the webDocuments map to a JSON file named web_docs.json
// 	file, err := os.Create("web_docs.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	encoder := json.NewEncoder(file)
// 	err = encoder.Encode(webDocuments)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

var (
	resultURLs []string
)

// SearchDuckDuckGo performs a search on DuckDuckGo and retrieves the HTML of the first page of results.
func SearchDuckDuckGo(chromeUrl, query string) []string {
	instanceUrl := chromeUrl

	// Create allocator context for using existing Chrome instance
	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), instanceUrl)
	defer cancel()

	// Create context with logging for actions performed by chromedp
	ctx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Set a timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Run tasks to perform the search and get HTML content
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://duckduckgo.com/`),
		chromedp.WaitVisible(`input[name="q"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="q"]`, query+kb.Enter, chromedp.ByQuery),

		// Wait for JavaScript to load the search results
		chromedp.Sleep(5*time.Second),

		// Wait for search results to be visible
		chromedp.WaitVisible(`button[id="more-results"]`, chromedp.ByQuery),

		chromedp.Nodes(`a`, &nodes, chromedp.ByQueryAll),

		chromedp.ActionFunc(func(c context.Context) error {
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			dom.RequestChildNodes(nodes[0].NodeID).WithDepth(-1).Do(c)

			// Compile the regex outside of the loop for efficiency
			re, err := regexp.Compile(`^http[s]?://`)
			if err != nil {
				return err
			}

			uniqueUrls := make(map[string]bool)

			// loop through the nodes and print the href attributes of a elements
			for _, n := range nodes {
				for _, attr := range n.Attributes {
					if re.MatchString(attr) {
						// Check if URL contains 'duckduckgo'
						if !strings.Contains(attr, "duckduckgo") {
							uniqueUrls[attr] = true
						}
					}
				}
			}

			// Convert map keys to slice
			resultURLs = make([]string, 0, len(uniqueUrls))
			for url := range uniqueUrls {
				resultURLs = append(resultURLs, url)
			}

			for _, l := range resultURLs {
				fmt.Println(l)
			}

			return nil
		}),
	)
	if err != nil {
		return resultURLs
	}

	return resultURLs
}

// RetrieveAndSanitizePage retrieves the HTML content of a page, sanitizes it, and returns the text.
func RetrieveAndSanitizePage(chromeUrl, pageAddress string) (string, error) {
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

	// Run tasks to get HTML content
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageAddress),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return "", err
	}

	fmt.Println("############################################")
	fmt.Println("############################################")
	fmt.Println(htmlContent)
	fmt.Println("############################################")
	fmt.Println("############################################")

	// Sanitize and format the HTML content
	sanitizedText, err := sanitizeHTML(htmlContent)
	if err != nil {
		return "", err
	}

	fmt.Println(sanitizedText)

	return sanitizedText, nil
}

// sanitizeHTML takes HTML content as input and returns sanitized text.
// It uses goquery to extract paragraphs, remove scripts and styles.
// Regex patterns are used to remove unnecessary words.
func sanitizeHTML(htmlContent string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var text string
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text += s.Text() + "\n\n"
	})

	doc.Find("script, style").Remove()

	// Regex patterns for removing unnecessary text
	redundantPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\b(very|really|quite|extremely|just|actually|basically|literally|surprisingly|interestingly)\b`),
		regexp.MustCompile(`\b(slightly|somewhat|kind of|sort of)\b`),
		regexp.MustCompile(`\b(past history|unexpected surprise|free gift|new innovations|true facts)\b`),
		regexp.MustCompile(`\b(whom|whereby|herein)\b`),
		regexp.MustCompile(`\b(utilize|ascertain)\b`),
		regexp.MustCompile(`\b(and|but|so|because|however)\b`),
		regexp.MustCompile(`\b(he|she|it|this|that)\b`),
		// TODO: Research more patterns.
	}

	// Apply regex patterns to remove redundant text
	for _, pattern := range redundantPatterns {
		text = pattern.ReplaceAllString(text, "")
	}

	// Additional text processing
	newlineRegex := regexp.MustCompile(`[\n\r]+`)
	text = newlineRegex.ReplaceAllString(text, "\n")

	spaceRegex := regexp.MustCompile(`[\s\t]+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	// Remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	text = htmlRegex.ReplaceAllString(text, "")

	return strings.TrimSpace(text), nil
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

// travelSubtree illustrates how to ask chromedp to populate a subtree of a node.
//
// https://github.com/chromedp/chromedp/issues/632#issuecomment-654213589
// @mvdan explains why node.Children is almost always empty:
// Nodes are only obtained from the browser on an on-demand basis.
// If we always held the entire DOM node tree in memory,
// our CPU and memory usage in Go would be far higher.
// And chromedp.FromNode can be used to retrieve the child nodes.
//
// Users get confused sometimes (why node.Children is empty while node.ChildNodeCount > 0?).
// And some users want to travel a subtree of the DOM more easy.
// So here comes the example.
func travelSubtree(pageUrl, of string, opts ...chromedp.QueryOption) chromedp.Tasks {
	var nodes []*cdp.Node
	return chromedp.Tasks{
		chromedp.Navigate(pageUrl),
		chromedp.Nodes(of, &nodes, opts...),
		// ask chromedp to populate the subtree of a node
		chromedp.ActionFunc(func(c context.Context) error {
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			return dom.RequestChildNodes(nodes[0].NodeID).WithDepth(-1).Do(c)
		}),
		// wait a little while for dom.EventSetChildNodes to be fired and handled
		chromedp.Sleep(time.Second),
		chromedp.ActionFunc(func(c context.Context) error {
			printNodes(os.Stdout, nodes, "", "  ")
			return nil
		}),
	}
}

func printNodes(w io.Writer, nodes []*cdp.Node, padding, indent string) {
	// This will block until the chromedp listener closes the channel
	for _, node := range nodes {
		switch {
		case node.NodeName == "#text":
			fmt.Fprintf(w, "%s#text: %q\n", padding, node.NodeValue)
		default:
			fmt.Fprintf(w, "%s%s:\n", padding, strings.ToLower(node.NodeName))
			if n := len(node.Attributes); n > 0 {
				fmt.Fprintf(w, "%sattributes:\n", padding+indent)
				for i := 0; i < n; i += 2 {
					fmt.Fprintf(w, "%s%s: %q\n", padding+indent+indent, node.Attributes[i], node.Attributes[i+1])
				}
			}
		}
		if node.ChildNodeCount > 0 {
			fmt.Fprintf(w, "%schildren:\n", padding+indent)
			printNodes(w, node.Children, padding+indent+indent, indent)
		}
	}
}
