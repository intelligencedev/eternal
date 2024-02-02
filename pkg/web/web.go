package web

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
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

	docs = strings.ReplaceAll(docs, "\n\n", "\n")
	fmt.Println("Document:", docs)

	return docs, nil
}

// serpapi handler, not free
func SearchHandler(c *fiber.Ctx) error {
	query := c.Query("q")
	apikey := c.Query("apikey")
	res, err := GetSerpResults(query, apikey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Only return the first 3 results
	// if len(*res) > 3 {
	// 	*res = (*res)[:3]
	// }

	// docs := ""

	// for _, r := range *res {
	// 	fmt.Printf("%s\n", r.Title)
	// 	fmt.Printf("%s\n", r.Link)

	// 	doc := web.Web2HTML(r.Link, true)

	// 	docs = fmt.Sprintf("%s %s", docs, doc)
	// }

	doc := (*res)[0].Link
	docs := Web2HTML(doc, true)

	return c.SendString(docs)
}

// ChromeDP Search Handler - Free
// func SearchChromeDPHandler(c *fiber.Ctx) error {
// 	docs := ""
// 	query := c.Query("q")
// 	apikey := c.Query("apikey")
// 	chromeDPAddress := fmt.Sprintf("ws://0.0.0.0:9222/devtools/browser/%s", apikey)
// 	res := SearchDuckDuckGo(chromeDPAddress, query)
// 	if len(res) == 0 {
// 		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving search results")
// 	}

// 	// TODO: Build a list of domains we want to exclude from our web workflows and expose that list
// 	// via a config file and/or frontend. Then filter out any results from those domains.
// 	// Remove youtube results
// 	for i := len(res) - 1; i >= 0; i-- {
// 		if strings.Contains(res[i], "youtube.com") {
// 			res = append(res[:i], res[i+1:]...)
// 		}
// 	}

// 	// Only return the first 3 results
// 	if len(res) > 3 {
// 		res = res[len(res)-5:]
// 	}

// 	for _, r := range res {
// 		fmt.Printf("%s\n", r)

// 		doc, _ := RetrieveAndSanitizePage(chromeDPAddress, r)

// 		log.Printf("Document: %s\n", doc)

// 		docs = fmt.Sprintf("%s %s", docs, doc)
// 	}

// 	// Remove all empty lines
// 	docs = strings.ReplaceAll(docs, "\n\n", "\n")

// 	return c.SendString(docs)
// }
