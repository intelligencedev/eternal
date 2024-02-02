package hf

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"eternal/pkg/web"
)

// ModelData represents a Hugging Face model.
type ModelData struct {
	ID          string    `json:"_id" gorm:"primaryKey"`
	ModelID     string    `json:"modelId"`
	Likes       int       `json:"likes"`
	IsPrivate   bool      `json:"private"`
	Downloads   int       `json:"downloads"`
	Tags        []string  `json:"tags" gorm:"type:text"`
	PipelineTag string    `json:"pipeline_tag"`
	LibraryName string    `json:"library_name"`
	CreatedAt   time.Time `json:"createdAt"`
}

// GetHFModels makes a Get requeset to the Hugging Face API to retrieve GGUF models
// search parameter is the name of the user/organization such as TheBloke, Microsoft, etc.
func GetHFModels(search string) ([]ModelData, error) {
	// Construct URL
	url := "https://huggingface.co/api/models?filter=gguf"
	if search != "" {
		url += "&search=" + search
	}

	// Make request
	httpClient := &web.HttpClient{}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err = fmt.Errorf("failed to close response body: %w", err)
			fmt.Println(err)
			return
		}
	}(resp.Body)

	// Decode response
	var models []ModelData
	err = json.NewDecoder(resp.Body).Decode(&models)
	if err != nil {
		return nil, err
	}

	return models, nil
}
