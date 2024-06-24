package main

import (
	"errors"
	"eternal/pkg/hfutils"
	"eternal/pkg/llm"
	"eternal/pkg/llm/openai"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pterm/pterm"
	"gorm.io/gorm"
)

// handleModelData retrieves and returns data for a specific model.
func handleModelData() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var model ModelParams
		modelName := c.Params("modelName")
		err := sqliteDB.First(modelName, &model)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("Model not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		return c.JSON(model)
	}
}

// handleModelDownloadUpdate updates the download status of a model.
func handleModelDownloadUpdate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		modelName := c.Params("modelName")
		var payload struct {
			Downloaded bool `json:"downloaded"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		err := sqliteDB.UpdateDownloadedByName(modelName, payload.Downloaded)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update model: %v", err)})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Model 'Downloaded' status updated successfully",
		})
	}
}

// handleModelUpdate updates the model data in the database.
func handleModelUpdate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var model ModelParams
		if err := c.BodyParser(&model); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Cannot parse JSON")
		}

		err := sqliteDB.UpdateByName(model.Name, model)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		return c.JSON(model)
	}
}

// handleModelCards retrieves and renders model cards.
func handleModelCards(modelParams []ModelParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := sqliteDB.Find(&modelParams)

		if err != nil {
			log.Errorf("Database error: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		return c.Render("templates/model", fiber.Map{"models": modelParams})
	}
}

// handleModelSelect handles the selection of models for use.
func handleModelSelect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		modelName := c.Params("name")
		action := c.Params("action")

		if action == "add" {
			if err := AddSelectedModel(sqliteDB.db, modelName); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		} else if action == "remove" {
			if err := RemoveSelectedModel(sqliteDB.db, modelName); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		} else {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid action")
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// handleSelectedModels retrieves and returns the list of selected models.
func handleSelectedModels() fiber.Handler {
	return func(c *fiber.Ctx) error {
		selectedModels, err := GetSelectedModels(sqliteDB.db)

		if err != nil {
			log.Errorf("Error getting selected models: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		var selectedModelNames []string
		for _, model := range selectedModels {
			selectedModelNames = append(selectedModelNames, model.ModelName)
		}

		return c.JSON(selectedModelNames)
	}
}

// handleModelDownload handles the download of a specified model.
func handleModelDownload(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pterm.Error.Println("Download route hit")
		modelName := c.Query("model")

		if modelName == "" {
			log.Errorf("Missing parameters for download")
			return c.Status(fiber.StatusBadRequest).SendString("Missing parameters")
		}

		var downloadURL string
		for _, model := range config.LanguageModels {
			if model.Name == modelName {
				downloadURL = model.Downloads[0]
				break
			}
		}

		modelFileName := filepath.Base(downloadURL)
		modelPath := filepath.Join(config.DataPath, "models", modelName, modelFileName)

		var partialDownload bool
		if info, err := os.Stat(modelPath); err == nil {
			if info.Size() > 0 {
				expectedSize, err := llm.GetExpectedFileSize(downloadURL)
				if err != nil {
					log.Errorf("Error getting expected file size: %v", err)
				}
				partialDownload = info.Size() < expectedSize
			}
		}

		go func() {
			var err error

			if partialDownload {
				pterm.Info.Printf("Resuming download for model: %s\n", modelName)
				err = llm.Download(downloadURL, modelPath)
			} else {
				pterm.Info.Printf("Starting download for model: %s\n", modelName)
				err = llm.Download(downloadURL, modelPath)
			}

			if err != nil {
				log.Errorf("Error in download: %v", err)
			} else {
				err = sqliteDB.UpdateDownloadedByName(modelName, true)
				if err != nil {
					log.Errorf("Failed to update model downloaded state: %v", err)
				}
			}
		}()

		progressErr := fmt.Sprintf("<div class='w-100' id='progress-download-%s' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message' hx-trigger='load'></div>", modelName)

		return c.SendString(progressErr)
	}
}

// handleImgModelDownload handles the download of image generation models.
func handleImgModelDownload(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		config.Tools.ImgGen.Enabled = true

		modelName := c.Query("model")

		var downloadURL string
		for _, model := range config.ImageModels {
			if model.Name == modelName {
				downloadURL = model.Downloads[0]
			}
		}

		modelFileName := strings.Split(downloadURL, "/")[len(strings.Split(downloadURL, "/"))-1]

		if modelName == "" {
			log.Errorf("Missing parameters for download")
			return c.Status(fiber.StatusBadRequest).SendString("Missing parameters")
		}

		modelRoot := fmt.Sprintf("%s/models/%s", config.DataPath, modelName)
		modelPath := fmt.Sprintf("%s/models/%s/%s", config.DataPath, modelName, modelFileName)
		tmpPath := fmt.Sprintf("%s/tmp", config.DataPath)

		if _, err := os.Stat(modelRoot); os.IsNotExist(err) {
			if err := os.MkdirAll(modelRoot, 0755); err != nil {
				log.Errorf("Error creating model directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
			if err := os.MkdirAll(tmpPath, 0755); err != nil {
				log.Errorf("Error creating tmp directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		if _, err := os.Stat(modelPath); err != nil {
			dm := hfutils.ConcurrentDownloadManager{
				FileName:    modelFileName,
				URL:         downloadURL,
				Destination: modelPath,
				NumParts:    1,
				TempDir:     tmpPath,
			}

			go dm.PrintProgress()

			if err := dm.Download(); err != nil {
				fmt.Println("Download failed:", err)
			} else {
				fmt.Println("Download successful!")
			}
		}

		vaeName := "sdxl_vae.safetensors"
		vaeURL := "https://huggingface.co/madebyollin/sdxl-vae-fp16-fix/blob/main/sdxl_vae.safetensors"
		vaePath := fmt.Sprintf("%s/models/%s/%s", config.DataPath, modelName, vaeName)

		if _, err := os.Stat(modelRoot); os.IsNotExist(err) {
			if err := os.MkdirAll(modelRoot, 0755); err != nil {
				log.Errorf("Error creating model directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		if _, err := os.Stat(vaePath); os.IsNotExist(err) {
			go func() {
				response, err := http.Get(vaeURL)
				if err != nil {
					pterm.Error.Printf("Failed to download file: %v", err)
					return
				}
				defer response.Body.Close()

				file, err := os.Create(vaePath)
				if err != nil {
					pterm.Error.Printf("Failed to create file: %v", err)
					return
				}
				defer file.Close()

				_, err = io.Copy(file, response.Body)
				if err != nil {
					pterm.Error.Printf("Failed to write to file: %v", err)
					return
				}

				pterm.Info.Printf("Downloaded file: %s", vaeName)
			}()
		}

		progressErr := "<div name='sse-messages' class='w-100' id='sse-messages' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message'></div>"

		return c.SendString(progressErr)
	}
}

// handleOpenAIModels retrieves and returns a list of OpenAI models.
func handleOpenAIModels(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		client := openai.NewClient(config.OAIKey)
		modelsResponse, err := openai.GetModels(client)

		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

		var gptModels []string
		for _, model := range modelsResponse.Data {
			if strings.HasPrefix(model.ID, "gpt") {
				gptModels = append(gptModels, model.ID)
			}
		}

		return c.JSON(fiber.Map{
			"object": "list",
			"data":   gptModels,
		})
	}
}
