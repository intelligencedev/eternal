package web

import (
	"encoding/json"
	"os"

	g "github.com/serpapi/google-search-results-golang"
)

type SearchData struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func GetSerpResults(query string, apikey string) (*[]SearchData, error) {
	parameter := map[string]string{
		"engine":  "google",
		"q":       query,
		"api_key": apikey,
	}

	search := g.NewGoogleSearch(parameter, apikey)
	results, err := search.GetJSON()
	if err != nil {
		return nil, err
	}

	resdata := []SearchData{}

	organic_results := results["organic_results"].([]interface{})
	for _, organic_result := range organic_results {
		organic_result := organic_result.(map[string]interface{})
		link := organic_result["link"].(string)
		title := organic_result["title"].(string)

		res := new(SearchData)
		res.Title = title
		res.Link = link

		// Only append if the link is not a PDF
		if len(link) > 4 && link[len(link)-4:] != ".pdf" {
			resdata = append(resdata, *res)
		}
	}

	// Write to JSON file
	file, _ := json.MarshalIndent(resdata, "", " ")
	_ = os.WriteFile("data.json", file, 0644)

	return &resdata, nil
}
