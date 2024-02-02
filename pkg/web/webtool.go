package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/temoto/robotstxt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
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
		cleanedURL := cleanURL(match) // Assuming cleanURL is a function you've defined elsewhere
		cleanedURLs = append(cleanedURLs, cleanedURL)
	}

	return cleanedURLs
}

func Web2HTML(url string, contentOnly bool) string {
	// Check robots.txt
	if !isAllowed(url, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36") {
		log.Println("Access to the URL is disallowed by robots.txt")
		return ""
	} else {
		// Fetch the URL
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching URL:", err)
			return ""
		}
		defer resp.Body.Close()

		// Convert the response body to UTF-8
		utf8Body, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
		if err != nil {
			fmt.Println("Error converting response body to UTF-8:", err)
			return ""
		}

		// Export the HTML
		html, err := io.ReadAll(utf8Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return ""
		}

		// Extract the article contents if contentOnly is true
		var contents []string
		if contentOnly {
			contents = extractArticleContents(string(html))

			// Remove all HTML tags
			re := regexp.MustCompile(`<[^>]*>`)
			for i, content := range contents {
				contents[i] = re.ReplaceAllString(content, "")
			}

			// // Remove all newlines
			// re = regexp.MustCompile(`\n`)
			// for i, content := range contents {
			// 	contents[i] = re.ReplaceAllString(content, "")
			// }

			// // Remove all tabs
			// re = regexp.MustCompile(`\t`)
			// for i, content := range contents {
			// 	contents[i] = re.ReplaceAllString(content, "")
			// }

			// // Remove all double spaces
			// re = regexp.MustCompile(`\s{2,}`)
			// for i, content := range contents {
			// 	contents[i] = re.ReplaceAllString(content, " ")
			// }

			// // Remove all leading and trailing spaces
			// re = regexp.MustCompile(`^\s+|\s+$`)
			// for i, content := range contents {
			// 	contents[i] = re.ReplaceAllString(content, "")
			// }

			// // Remove all special characters
			// re = regexp.MustCompile(`[^a-zA-Z0-9\s]`)
			// for i, content := range contents {
			// 	contents[i] = re.ReplaceAllString(content, "")
			// }

			// Convert to UTF-8
			for i, content := range contents {
				contents[i] = strings.ToValidUTF8(content, "")
			}

		} else {
			contents = append(contents, string(html))
		}

		// Convert the contents to []byte
		var contentsBytes []byte
		for _, content := range contents {
			contentsBytes = append(contentsBytes, []byte(content)...)
		}

		// Check if ./webcache exists and if not create it
		if _, err := os.Stat("./webcache"); os.IsNotExist(err) {
			os.Mkdir("./webcache", 0755)
		}

		// Save the HTML
		err = os.WriteFile("./webcache/output.html", contentsBytes, 0644)
		if err != nil {
			fmt.Println("Error saving HTML:", err)
			return ""
		}

		return string(contentsBytes)
	}
}

// cleanURL removes illegal trailing characters from the URL.
func cleanURL(url string) string {
	// Define illegal trailing characters. Adjust as needed.
	illegalTrailingChars := []rune{'.', ',', ';', '!', '?', ')'}

	for _, char := range illegalTrailingChars {
		if url[len(url)-1] == byte(char) {
			url = url[:len(url)-1]
		}
	}

	return url
}

// We must remove all html tags for the llm to process it.
func extractArticleContents(input string) []string {
	var contents []string

	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			var innerContent string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				innerContent += renderNode(c)
			}
			contents = append(contents, innerContent)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return contents
}

func renderNode(n *html.Node) string {
	var buf strings.Builder
	html.Render(&buf, n)
	return buf.String()
}

func isAllowed(targetURL, userAgent string) bool {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return false
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)
	resp, err := http.Get(robotsURL)
	if err != nil {
		// If there's an error, assume the URL is allowed
		return true
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		// If there's an error, assume the URL is allowed
		return true
	}

	robots, err := robotstxt.FromBytes(data)
	if err != nil {
		// If there's an error, assume the URL is allowed
		return true
	}

	group := robots.FindGroup(userAgent)
	return group.Test(targetURL)
}
