package sd

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
	"os/exec"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gofiber/fiber/v2"
)

type SDParams struct {
	Mode           string  `json:"mode" cli:"-M, --mode"`
	Threads        int     `json:"threads" cli:"-t, --threads"`
	Model          string  `json:"model" cli:"-m, --model"`
	VAE            string  `json:"vae" cli:"--vae"`
	TAESDPath      string  `json:"taesd" cli:"--taesd"`
	WeightType     string  `json:"type" cli:"--type"`
	LoraModelDir   string  `json:"lora_model_dir" cli:"--lora-model-dir"`
	InitImage      string  `json:"init_img" cli:"-i, --init-img"`
	Output         string  `json:"output" cli:"-o, --output"`
	Prompt         string  `json:"prompt" cli:"-p, --prompt"`
	NegativePrompt string  `json:"negative_prompt" cli:"-n, --negative-prompt"`
	CFGScale       float64 `json:"cfg_scale" cli:"--cfg-scale"`
	Strength       float64 `json:"strength" cli:"--strength"`
	Height         int     `json:"height" cli:"-H, --height"`
	Width          int     `json:"width" cli:"-W, --width"`
	SamplingMethod string  `json:"sampling_method" cli:"--sampling-method"`
	Steps          int     `json:"steps" cli:"--steps"`
	RNG            string  `json:"rng" cli:"--rng"`
	Seed           int     `json:"seed" cli:"-s, --seed"`
	BatchCount     int     `json:"batch_count" cli:"-b, --batch-count"`
	Schedule       string  `json:"schedule" cli:"--schedule"`
	Verbose        bool    `json:"verbose" cli:"-v, --verbose"`
}

func ConstructCLICommand(params SDParams) *exec.Cmd {
	cmdPath := "/home/art/.eternal/sd"

	cmdArgs := []string{
		"-m", params.Model,
		"-p", params.Prompt,
		"-o", "./public/img/sd_output.png",
		"--cfg-scale", "7",
		"--sampling-method", "dpm++2m", //"dpm++2m",
		"--steps", "20",
		"--seed", "-1",
		//"--upscale-model", "/mnt/d/StableDiffusionModels/sdxl/upscalers/RealESRGAN_x4plus_anime_6B.pth",
		"--schedule", "karras",
		"--clip-skip", "1",
	}

	// Conditional arguments
	if params.InitImage != "" {
		cmdArgs = append(cmdArgs, "-i", params.InitImage)
	}
	if params.VAE != "" {
		cmdArgs = append(cmdArgs, "--vae", params.VAE)
	}

	if params.Mode != "" {
		cmdArgs = append(cmdArgs, "--mode", params.Mode)
	}

	return exec.Command(cmdPath, cmdArgs...)
}

func Text2Image(c *fiber.Ctx) error {
	var params SDParams

	err := json.Unmarshal(c.Body(), &params)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	cmd := ConstructCLICommand(params)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"output": string(out),
	})
}

// cropAndAdjustContrast crops the image to a 768x768 square while keeping the subject centered,
// increases the contrast slightly, and overwrites the original image.
func CropAndAdjustContrast(filePath string) error {
	// Load the image from the file path
	img, err := imgio.Open(filePath)
	if err != nil {
		return err
	}

	// Determine the size and position for cropping
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	minDim := int(math.Min(float64(width), float64(height)))
	offsetX, offsetY := (width-minDim)/2, (height-minDim)/2

	// Crop the image to a square while centering the subject
	croppedImg := transform.Crop(img, image.Rect(offsetX, offsetY, offsetX+minDim, offsetY+minDim))

	// Resize the image to 768x768 if necessary
	if minDim != 768 {
		croppedImg = transform.Resize(croppedImg, 768, 768, transform.Linear)
	}

	// Increase the contrast
	adjustedImg := adjust.Contrast(croppedImg, 0.2) // Adjust the contrast level as needed

	// Save the image back to the same file
	err = imgio.Save(filePath, adjustedImg, imgio.JPEGEncoder(100))
	if err != nil {
		return err
	}

	return nil
}
