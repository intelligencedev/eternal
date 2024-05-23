package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

prompt := "A beautiful purple flower in a dark forest."

// Define the JSON payload
var payload = `
{
    "prompt": {
        "4": {
            "inputs": {
                "ckpt_name": "canvasdarkxl_v10.safetensors"
            },
            "class_type": "CheckpointLoaderSimple",
            "_meta": {
                "title": "Load Checkpoint"
            }
        },
        "56": {
            "inputs": {
                "resolution": "1280x768 (1.67)",
                "batch_size": 1
            },
            "class_type": "SDXLEmptyLatentSizePicker+",
            "_meta": {
                "title": "🔧 SDXL Empty Latent Size Picker"
            }
        },
        "77": {
            "inputs": {
                "stop_at_clip_layer": -1,
                "clip": [
                    "4",
                    1
                ]
            },
            "class_type": "CLIPSetLastLayer",
            "_meta": {
                "title": "CLIP Set Last Layer"
            }
        },
        "225": {
            "inputs": {
                "text": "A beautiful purple flower in a dark forest.",
                "clip": [
                    "229",
                    1
                ]
            },
            "class_type": "CLIPTextEncode",
            "_meta": {
                "title": "CLIP Text Encode (Positive Prompt)"
            }
        },
        "226": {
            "inputs": {
                "text": "bad quality, low quality, watermark, text, embedding:FastNegativeV2, ",
                "clip": [
                    "77",
                    0
                ]
            },
            "class_type": "CLIPTextEncode",
            "_meta": {
                "title": "CLIP Text Encode (Negative Prompt)"
            }
        },
        "227": {
            "inputs": {
                "lora_name": "midjourney_20230624181825.safetensors",
                "strength_model": 0.8,
                "strength_clip": 0.8,
                "model": [
                    "4",
                    0
                ],
                "clip": [
                    "77",
                    0
                ]
            },
            "class_type": "LoraLoader",
            "_meta": {
                "title": "Load LoRA"
            }
        },
        "228": {
            "inputs": {
                "lora_name": "Expressive_H-000001.safetensors",
                "strength_model": 0.8,
                "strength_clip": 0.8,
                "model": [
                    "227",
                    0
                ],
                "clip": [
                    "227",
                    1
                ]
            },
            "class_type": "LoraLoader",
            "_meta": {
                "title": "Load LoRA"
            }
        },
        "229": {
            "inputs": {
                "lora_name": "sinfully_stylish_SDXL.safetensors",
                "strength_model": 1,
                "strength_clip": 1,
                "model": [
                    "228",
                    0
                ],
                "clip": [
                    "228",
                    1
                ]
            },
            "class_type": "LoraLoader",
            "_meta": {
                "title": "Load LoRA"
            }
        },
        "356": {
            "inputs": {
                "filename_prefix": "api_img_",
                "images": [
                    "237:2",
                    0
                ]
            },
            "class_type": "SaveImage",
            "_meta": {
                "title": "Save Image"
            }
        },
        "233:0": {
            "inputs": {
                "scale": 1.25,
                "adaptive_scale": 0,
                "unet_block": "middle",
                "unet_block_id": 0,
                "sigma_start": false,
                "sigma_end": -1,
                "rescale": false,
                "rescale_mode": "full",
                "model": [
                    "4",
                    0
                ]
            },
            "class_type": "PerturbedAttention",
            "_meta": {
                "title": "Perturbed-Attention Guidance (Advanced)"
            }
        },
        "233:1": {
            "inputs": {
                "hard_mode": true,
                "boost": false,
                "model": [
                    "233:0",
                    0
                ]
            },
            "class_type": "Automatic CFG",
            "_meta": {
                "title": "Automatic CFG"
            }
        },
        "237:0": {
            "inputs": {
                "vae_name": "sdxl_vae.safetensors"
            },
            "class_type": "VAELoader",
            "_meta": {
                "title": "Load VAE"
            }
        },
        "237:1": {
            "inputs": {
                "seed": 76445043577479,
                "steps": 35,
                "cfg": 7,
                "sampler_name": "dpm_2",
                "scheduler": "karras",
                "denoise": 1,
                "model": [
                    "233:1",
                    0
                ],
                "positive": [
                    "225",
                    0
                ],
                "negative": [
                    "226",
                    0
                ],
                "latent_image": [
                    "56",
                    0
                ]
            },
            "class_type": "KSampler",
            "_meta": {
                "title": "KSampler"
            }
        },
        "237:2": {
            "inputs": {
                "samples": [
                    "237:1",
                    0
                ],
                "vae": [
                    "237:0",
                    0
                ]
            },
            "class_type": "VAEDecode",
            "_meta": {
                "title": "VAE Decode"
            }
        }
    }
}`

// SendRequest sends a POST request to the specified endpoint with the provided JSON data
func SendRequest() error {
	// Define the endpoint URL
	url := "http://192.168.0.148:8188/prompt"

	// Generate a random seed
	seed := make([]byte, 8)

	_, err := rand.Read(seed)
	if err != nil {
		return err
	}

	// Convert the seed to a string
	seedStr := strconv.FormatInt(int64(seed[0])<<56|int64(seed[1])<<48|int64(seed[2])<<40|int64(seed[3])<<32|int64(seed[4])<<24|int64(seed[5])<<16|int64(seed[6])<<8|int64(seed[7]), 10)

	// Replace the seed in the payload
	payload = string(bytes.ReplaceAll([]byte(payload), []byte("76445043577479"), []byte(seedStr)))

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Close the response body
	defer resp.Body.Close()

	// Decode the response body
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Print the response
	fmt.Println(result)

	return nil
}

func main() {
	// Send the HTTP request
	if err := SendRequest(); err != nil {
		fmt.Println("Error:", err)
	}
}
