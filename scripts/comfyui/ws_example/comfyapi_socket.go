package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

const serverAddress = "192.168.0.148:8188"

type Prompt struct {
	Prompt   map[string]interface{} `json:"prompt"`
	ClientID string                 `json:"client_id"`
}

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func queuePrompt(prompt map[string]interface{}, clientID string) (map[string]interface{}, error) {
	p := Prompt{
		Prompt:   prompt,
		ClientID: clientID,
	}
	jsonData, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/prompt", serverAddress), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func getImages(c *websocket.Conn, prompt map[string]interface{}) (map[string][][]byte, error) {
	clientID := uuid.New().String()
	result, err := queuePrompt(prompt, clientID)
	if err != nil {
		return nil, err
	}

	promptID := result["prompt_id"].(string)
	outputImages := make(map[string][][]byte)
	currentNode := ""

	for {
		messageType, msg, err := c.ReadMessage()
		if err != nil {
			return nil, err
		}

		if messageType == websocket.TextMessage {
			var message Message
			err = json.Unmarshal(msg, &message)
			if err != nil {
				return nil, err
			}

			if message.Type == "executing" {
				data := message.Data
				if data["prompt_id"] == promptID {
					if data["node"] == nil {
						break // Execution is done
					} else {
						currentNode = data["node"].(string)
					}
				}
			}
		} else if messageType == websocket.BinaryMessage {
			if currentNode == "save_image_websocket_node" {
				outputImages[currentNode] = append(outputImages[currentNode], msg[8:])
			}
		}
	}

	return outputImages, nil
}

func main() {
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		promptText := `
		{
			"3": {
				"class_type": "KSampler",
				"inputs": {
					"cfg": 8,
					"denoise": 1,
					"latent_image": [
						"5",
						0
					],
					"model": [
						"4",
						0
					],
					"negative": [
						"7",
						0
					],
					"positive": [
						"6",
						0
					],
					"sampler_name": "euler",
					"scheduler": "normal",
					"seed": 8566257,
					"steps": 20
				}
			},
			"4": {
				"class_type": "CheckpointLoaderSimple",
				"inputs": {
					"ckpt_name": "canvasdarkxl_v10.safetensors"
				}
			},
			"5": {
				"class_type": "EmptyLatentImage",
				"inputs": {
					"batch_size": 1,
					"height": 1024,
					"width": 1024
				}
			},
			"6": {
				"class_type": "CLIPTextEncode",
				"inputs": {
					"clip": [
						"4",
						1
					],
					"text": "masterpiece best quality girl"
				}
			},
			"7": {
				"class_type": "CLIPTextEncode",
				"inputs": {
					"clip": [
						"4",
						1
					],
					"text": "bad hands"
				}
			},
			"8": {
				"class_type": "VAEDecode",
				"inputs": {
					"samples": [
						"3",
						0
					],
					"vae": [
						"4",
						2
					]
				}
			},
			"save_image_websocket_node": {
				"class_type": "SaveImageWebsocket",
				"inputs": {
					"images": [
						"8",
						0
					]
				}
			}
		}
		`

		var prompt map[string]interface{}
		err := json.Unmarshal([]byte(promptText), &prompt)
		if err != nil {
			fmt.Println("Error unmarshaling prompt:", err)
			return
		}

		// Set the text prompt for our positive CLIPTextEncode
		prompt["6"].(map[string]interface{})["inputs"].(map[string]interface{})["text"] = "beautiful scenery nature glass bottle landscape, , purple galaxy bottle,"

		// Set the seed for our KSampler node
		prompt["3"].(map[string]interface{})["inputs"].(map[string]interface{})["seed"] = 5

		images, err := getImages(c, prompt)
		if err != nil {
			fmt.Println("Error getting images:", err)
			return
		}

		fmt.Printf("Received %d images\n", len(images["save_image_websocket_node"]))
	}))

	err := app.Listen(":8081")
	if err != nil {
		panic(err)
	}
}
