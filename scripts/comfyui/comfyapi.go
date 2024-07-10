package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const promptText = `
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
            "height": 512,
            "width": 512
        }
    },
    "6": {
        "class_type": "CLIPTextEncode",
        "inputs": {
            "clip": [
                "4",
                1
            ],
            "text": "beautiful scenery nature glass bottle landscape, , purple galaxy bottle,"
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
    "9": {
        "class_type": "SaveImage",
        "inputs": {
            "filename_prefix": "ComfyUI",
            "images": [
                "8",
                0
            ]
        }
    }
}
`

func queuePrompt(prompt map[string]interface{}) error {
	p := map[string]interface{}{"prompt": prompt}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://192.168.0.148:8188/prompt", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func main() {
	var prompt map[string]interface{}
	err := json.Unmarshal([]byte(promptText), &prompt)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// set the text prompt for our positive CLIPTextEncode
	prompt["6"].(map[string]interface{})["inputs"].(map[string]interface{})["text"] = "beautiful scenery nature glass bottle landscape, , purple galaxy bottle,"

	// set the seed for our KSampler node
	prompt["3"].(map[string]interface{})["inputs"].(map[string]interface{})["seed"] = 5

	err = queuePrompt(prompt)
	if err != nil {
		fmt.Println("Error queuing prompt:", err)
		return
	}

	fmt.Println("Prompt queued successfully")
}
