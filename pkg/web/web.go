// web.go
package web

import (
	"context"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/go-shiori/go-readability"
	"github.com/pterm/pterm"
)

var (
	unwantedURLs = []string{
		"web.archive.org",
		"www.youtube.com",
		"www.youtube.com/watch",
		"www.wired.com",
		"www.techcrunch.com",
		"www.wsj.com",
		"www.cnn.com",
		"www.nytimes.com",
		"www.forbes.com",
		"www.businessinsider.com",
		"www.theverge.com",
		"www.thehill.com",
		"www.theatlantic.com",
		"www.foxnews.com",
		"www.theguardian.com",
		"www.nbcnews.com",
		"www.msn.com",
		"www.sciencedaily.com",
		"reuters.com",
		"bbc.com",
		"thenewstack.io",
		"abcnews.go.com",
		"apnews.com",
		"bloomberg.com",
		"polygon.com",
		"reddit.com",
		"indeed.com",
		"test.com",
		// Add more URLs to block from search results
	}

	resultURLs []string
)

func WebGetHandler(address string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var docs string
	err := chromedp.Run(ctx,
		chromedp.Navigate(address),
		chromedp.ActionFunc(func(ctx context.Context) error {
			headers := map[string]interface{}{
				"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
				"Referer":         "https://www.google.com/",
				"Accept-Language": "en-US,en;q=0.9",
				"X-Forwarded-For": "203.0.113.195",
				"Accept-Encoding": "gzip, deflate, br",
				"Connection":      "keep-alive",
				"DNT":             "1",
			}
			return network.SetExtraHTTPHeaders(network.Headers(headers)).Do(ctx)
		}),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &docs),
	)

	if err != nil {
		log.Println("Error retrieving page:", err)
		return "", err
	}

	// Convert url to url.URL
	getUrl, err := url.Parse(address)
	if err != nil {
		log.Println("Error parsing URL:", err)
		return "", err
	}

	article, err := readability.FromReader(strings.NewReader(docs), getUrl)
	if err != nil {
		log.Println("Error parsing reader view:", err)
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
	if err != nil {
		log.Println("Error parsing document:", err)
		return "", err
	}

	text := doc.Find("body").Text()

	text = removeEmptyRows(text)

	//pterm.Info.Println("Document:", text)

	return text, nil
}

func ExtractURLs(input string) []string {
	urlRegex := `http.*?://[^\s<>{}|\\^` + "`" + `"]+`
	re := regexp.MustCompile(urlRegex)

	matches := re.FindAllString(input, -1)

	var cleanedURLs []string
	for _, match := range matches {
		cleanedURL := cleanURL(match)
		cleanedURLs = append(cleanedURLs, cleanedURL)
	}

	return cleanedURLs
}

func RemoveUrl(input []string) []string {
	urlRegex := `http.*?://[^\s<>{}|\\^` + "`" + `"]+`
	re := regexp.MustCompile(urlRegex)

	for i, str := range input {
		matches := re.FindAllString(str, -1)
		for _, match := range matches {
			input[i] = strings.ReplaceAll(input[i], match, "")
		}
	}

	return input
}

func cleanURL(url string) string {
	illegalTrailingChars := []rune{'.', ',', ';', '!', '?', ')'}

	for _, char := range illegalTrailingChars {
		if url[len(url)-1] == byte(char) {
			url = url[:len(url)-1]
		}
	}

	return url
}

func SearchDDG(query string) []string {
	resultURLs = nil

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var nodes []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://lite.duckduckgo.com/lite/`),
		chromedp.WaitVisible(`input[name="q"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="q"]`, query+kb.Enter, chromedp.ByQuery),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`input[name="q"]`, chromedp.ByQuery),
		chromedp.Nodes(`a`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Printf("Error during search: %v", err)
		return nil
	}

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

			for u := range uniqueUrls {
				resultURLs = append(resultURLs, u)
			}

			return nil
		}),
	)

	if err != nil {
		log.Printf("Error processing results: %v", err)
		return nil
	}

	resultURLs = RemoveUnwantedURLs(resultURLs)

	pterm.Warning.Println("Search results:", resultURLs)

	return resultURLs
}

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
	}

	return resultHTML
}

func RemoveUnwantedURLs(urls []string) []string {
	var filteredURLs []string
	for _, u := range urls {
		pterm.Info.Printf("Checking URL: %s", u)

		unwanted := false
		for _, unwantedURL := range unwantedURLs {
			if strings.Contains(u, unwantedURL) {
				pterm.Warning.Printf("URL %s contains unwanted URL %s", u, unwantedURL)
				unwanted = true
				break
			}
		}
		if !unwanted {
			filteredURLs = append(filteredURLs, u)
		}
	}

	pterm.Info.Printf("Filtered URLs: %v", filteredURLs)

	return filteredURLs
}

func GetPageScreen(chromeUrl string, pageAddress string) string {
	instanceUrl := chromeUrl

	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), instanceUrl)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageAddress),
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(pageAddress)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	filename := u.Hostname() + "-" + t.Format("20060102150405") + ".png"

	err = os.WriteFile(filename, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}

func RemoveUrls(input string) string {
	urlRegex := `http.*?://[^\s<>{}|\\^` + "`" + `"]+`
	re := regexp.MustCompile(urlRegex)

	matches := re.FindAllString(input, -1)

	for _, match := range matches {
		input = strings.ReplaceAll(input, match, "")
	}

	return input
}

func removeEmptyRows(input string) string {
	lines := strings.Split(input, "\n")
	var filteredLines []string

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}

	return strings.Join(filteredLines, "\n")
}
