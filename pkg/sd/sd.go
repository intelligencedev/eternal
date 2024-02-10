package sd

import (
	"fmt"
	"image"
	"math"
	"os/exec"
	"path/filepath"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/pterm/pterm"
)

type CommandOutput struct {
	Output       string `json:"output"`
	Finished     string `json:"finished"`
	SocketNumber string `json:"socketNumber"`
	ModelName    string `json:"modelName"`
}

type SDParams struct {
	Mode            string  `json:"mode" cli:"-M, --mode"`
	Threads         int     `json:"threads" cli:"-t, --threads"`
	Model           string  `json:"model" cli:"-m, --model"`
	VAE             string  `json:"vae" cli:"--vae"`
	TAESDPath       string  `json:"taesd" cli:"--taesd"`
	ControlNet      string  `json:"control_net" cli:"--control-net"`
	EmbdDir         string  `json:"embd_dir" cli:"--embd-dir"`
	UpscaleModel    string  `json:"upscale_model" cli:"--upscale-model"`
	WeightType      string  `json:"type" cli:"--type"`
	LoraModelDir    string  `json:"lora_model_dir" cli:"--lora-model-dir"`
	InitImage       string  `json:"init_img" cli:"-i, --init-img"`
	ControlImage    string  `json:"control_image" cli:"--control-image"`
	Output          string  `json:"output" cli:"-o, --output"`
	Prompt          string  `json:"prompt" cli:"-p, --prompt"`
	NegativePrompt  string  `json:"negative_prompt" cli:"-n, --negative-prompt"`
	CFGScale        float64 `json:"cfg_scale" cli:"--cfg-scale"`
	Strength        float64 `json:"strength" cli:"--strength"`
	ControlStrength float64 `json:"control_strength" cli:"--control-strength"`
	Height          int     `json:"height" cli:"-H, --height"`
	Width           int     `json:"width" cli:"-W, --width"`
	SamplingMethod  string  `json:"sampling_method" cli:"--sampling-method"`
	Steps           int     `json:"steps" cli:"--steps"`
	RNG             string  `json:"rng" cli:"--rng"`
	Seed            int     `json:"seed" cli:"-s, --seed"`
	BatchCount      int     `json:"batch_count" cli:"-b, --batch-count"`
	Schedule        string  `json:"schedule" cli:"--schedule"`
	ClipSkip        int     `json:"clip_skip" cli:"--clip-skip"`
	VAETiling       bool    `json:"vae_tiling" cli:"--vae-tiling"`
	ControlNetCPU   bool    `json:"control_net_cpu" cli:"--control-net-cpu"`
	Canny           bool    `json:"canny" cli:"--canny"`
	Verbose         bool    `json:"verbose" cli:"-v, --verbose"`
}

func BuildCommand(dataPath string, params SDParams) *exec.Cmd {
	//vaePath := filepath.Join(dataPath, "models/StableDiffusion/sd15/sdxl_vae.safetensors")
	modelPath := filepath.Join(dataPath, "models/StableDiffusion/sd15/dreamshaper_8_q5_1.gguf")
	outPath := filepath.Join(dataPath, "web/img/sd_out.png")
	cmdPath := filepath.Join(dataPath, "sd/sd")

	pterm.Println("Command:", cmdPath)

	cmdArgs := []string{
		"-p", params.Prompt,
		"-n", "ugly, low quality, deformed",
		"-m", modelPath,
		//"--vae", vaePath, // NOT WORKING DO NOT USE
		"-o", outPath,
		"--rng", "std_default",
		//"--cfg-scale", "7",
		"--sampling-method", "dpm2",
		//"--steps", "20",
		"--seed", "-1",
		//"--upscale-model", "/mnt/d/StableDiffusionModels/sdxl/upscalers/RealESRGAN_x4plus_anime_6B.pth",
		"--schedule", "karras",
		"--clip-skip", "2",
	}

	// Print cmdArgs
	pterm.Println("Command Args:", cmdArgs)

	// Run the command in a separate process
	return exec.Command(cmdPath, cmdArgs...)

	//return exec.Command(cmdPath, cmdArgs...)
}

func Text2Image(cmdPath string, params *SDParams) error {
	cmd := BuildCommand(cmdPath, *params)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}

	pterm.Warning.Println("Output:", string(out))

	// return the output
	return nil
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
