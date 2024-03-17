package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func queuePrompt(prompt map[string]interface{}) {
	p := map[string]interface{}{
		"prompt": prompt,
	}
	data, _ := json.Marshal(p)

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8188/prompt", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	app := fiber.New()

	app.Post("/prompt", func(c *fiber.Ctx) error {
		// Load the JSON file
		file, err := ioutil.ReadFile("prompt.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to load JSON file")
		}

		var prompt map[string]interface{}
		err = json.Unmarshal(file, &prompt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse JSON file")
		}

		// Check if the required keys exist and have the correct data type
		if val, ok := prompt["6"].(map[string]interface{}); ok {
			if inputs, ok := val["inputs"].(map[string]interface{}); ok {
				inputs["text"] = "a beautiful witch"
			} else {
				return c.Status(fiber.StatusInternalServerError).SendString("Invalid structure for key '6'")
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString("Missing or invalid key '6'")
		}

		if val, ok := prompt["3"].(map[string]interface{}); ok {
			if inputs, ok := val["inputs"].(map[string]interface{}); ok {
				inputs["seed"] = 5
			} else {
				return c.Status(fiber.StatusInternalServerError).SendString("Invalid structure for key '3'")
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString("Missing or invalid key '3'")
		}

		queuePrompt(prompt)

		return c.SendString("Prompt queued")
	})

	app.Listen(":3000")
}
