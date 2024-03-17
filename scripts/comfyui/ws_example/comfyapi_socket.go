package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

const serverAddress = "127.0.0.1:8188"

func queuePrompt(prompt map[string]interface{}, clientID string) map[string]interface{} {
	p := map[string]interface{}{
		"prompt":    prompt,
		"client_id": clientID,
	}
	data, _ := json.Marshal(p)

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/prompt", serverAddress), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func getImage(filename, subfolder, folderType string) []byte {
	data := url.Values{
		"filename":  {filename},
		"subfolder": {subfolder},
		"type":      {folderType},
	}

	resp, _ := http.Get(fmt.Sprintf("http://%s/view?%s", serverAddress, data.Encode()))
	defer resp.Body.Close()

	image, _ := ioutil.ReadAll(resp.Body)
	return image
}

func getHistory(promptID string) map[string]interface{} {
	resp, _ := http.Get(fmt.Sprintf("http://%s/history/%s", serverAddress, promptID))
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func getImages(c *websocket.Conn, prompt map[string]interface{}) map[string][][]byte {
	clientID := uuid.New().String()
	promptID := queuePrompt(prompt, clientID)["prompt_id"].(string)
	outputImages := make(map[string][][]byte)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		var data map[string]interface{}
		json.Unmarshal(message, &data)

		if data["type"] == "executing" {
			executionData := data["data"].(map[string]interface{})
			if executionData["node"] == nil && executionData["prompt_id"] == promptID {
				break // Execution is done
			}
		}
	}

	history := getHistory(promptID)[promptID].(map[string]interface{})
	outputs := history["outputs"].(map[string]interface{})

	for nodeID := range outputs {
		nodeOutput := outputs[nodeID].(map[string]interface{})
		if images, ok := nodeOutput["images"].([]interface{}); ok {
			var imagesOutput [][]byte
			for _, image := range images {
				imageData := image.(map[string]interface{})
				filename := imageData["filename"].(string)
				subfolder := imageData["subfolder"].(string)
				imageType := imageData["type"].(string)
				imageBytes := getImage(filename, subfolder, imageType)
				imagesOutput = append(imagesOutput, imageBytes)
			}
			outputImages[nodeID] = imagesOutput
		}
	}

	return outputImages
}

func main() {
	app := fiber.New()

	app.Post("/prompt", func(c *fiber.Ctx) error {
		// Load the JSON file
		file, err := os.Open("comfy_workflow.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to load JSON file")
		}
		defer file.Close()

		var prompt map[string]interface{}
		err = json.NewDecoder(file).Decode(&prompt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse JSON
			file")
		}

		// Check if the required keys exist and have the correct data type
		if val, ok := prompt["6"].(map[string]interface{}); ok {
			if inputs, ok := val["inputs"].(map[string]interface{}); ok {
				inputs["text"] = "a beautiful witch"
			} else {
				return c.Status(fiber.StatusInternalServerError).SendString("Invalid structure for key '6'")
			}
		}
		else {
			return c.Status(fiber.StatusInternalServerError).SendString("Missing or invalid key '6'")
		}

		if val, ok := prompt["3"].(map[string]interface{}); ok {
			if inputs, ok := val["inputs"].(map[string]interface{}); ok {
				inputs["seed"] = 12345
			} else {
				return c.Status(fiber.StatusInternalServerError).SendString("Invalid structure for key '3'")
			}
		}

		// Queue the prompt
		promptID := queuePrompt(prompt, uuid.New().String())["prompt_id"].(string)
		fmt.Println("Prompt ID:", promptID)

		// Create a websocket connection
		ws, err := websocket.Connect(c.Context(), fmt.Sprintf("ws://%s/ws", serverAddress))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to connect to websocket")
		}
		defer ws.Close()

		// Get the images from the prompt
		outputImages := getImages(ws, prompt)
		fmt.Println("Images:", outputImages)

		return c.SendString("Prompt queued")
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Load the JSON data from the file
		file, err := os.Open("comfy_workflow.json")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		var prompt map[string]interface{}
		err = json.NewDecoder(file).Decode(&prompt)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Set the text prompt for our positive CLIPTextEncode
		prompt["6"].(map[string]interface{})["inputs"].(map[string]interface{})["text"] = "masterpiece best quality man"

		// Set the seed for our KSampler node
		prompt["3"].(map[string]interface{})["inputs"].(map[string]interface{})["seed"] = 5

		//images := getImages(c, prompt)

		// Commented out code to display the output images:

		// for nodeID := range images {
		// 	for _, imageData := range images[nodeID] {
		// 		// Process the image data as needed
		// 	}
		// }

		c.Close()
	}))

	app.Listen(":3000")
}
