package sd

import (
	"fmt"
	"image"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/pterm/pterm"
)

type ImageModel struct {
	Name       string   `yaml:"name"`
	Homepage   string   `yaml:"homepage"`
	Prompt     string   `yaml:"prompt"`
	Downloads  []string `yaml:"downloads,omitempty"`
	Downloaded bool     `yaml:"downloaded"`
	LocalPath  string   `yaml:"localPath,omitempty"`
}

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
	Height          int     `json:"height" cli:"--height"`
	Width           int     `json:"width" cli:"--width"`
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
	//modelPath := filepath.Join(dataPath, "models/StableDiffusion/sd15/dreamshaper_8_q5_1.gguf")
	modelPath := filepath.Join(dataPath, "models/dreamshaper-8-sd15/128713")
	outPath := filepath.Join(dataPath, "web/img/sd_out.png")
	//cmdPath := filepath.Join(dataPath, "sd/sd")
	cmdPath := filepath.Join(dataPath, "sd")

	pterm.Println("Command:", cmdPath)

	cmdArgs := []string{
		"-p", params.Prompt,
		"-n", "ugly, low quality, deformed, malformed, floating limbs, bad hands, poorly drawn, bad anatomy, extra limb, blurry, disfigured, realistic, child",
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
		"--width", "512",
		"--height", "768",
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

var (
	downloadProgressMap = make(map[string]DownloadProgress)
	progressMutex       sync.Mutex
	TurnCounter         = 0
)

// DownloadProgress structure to hold download progress information
type DownloadProgress struct {
	Total   int64 `json:"total"`
	Current int64 `json:"current"`
}

// WebModel defines the interface for a model that interacts over the web.
type WebModel interface {
	Connect(endpoint string) (*http.Response, error)
}

// ModelManager defines the interface for managing local models.
type ModelManager interface {
	GetConfig() error
	Download() error
	Delete() error
}

type Model struct {
	Name      string   `yaml:"name"`
	Homepage  string   `yaml:"homepage"`
	Prompt    string   `yaml:"prompt"`
	Ctx       int      `yaml:"ctx"`
	Roles     []string `yaml:"roles"`
	Tags      []string `yaml:"tags,omitempty"`
	GGUF      string   `yaml:"gguf,omitempty"`
	Downloads []string `yaml:"downloads,omitempty"`
	LocalPath string   `yaml:"localPath,omitempty"`
}

type ProgressReader struct {
	Reader        io.Reader
	ProgressBar   *pterm.ProgressbarPrinter
	TotalRead     int64
	ContentLength int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.TotalRead += int64(n)

	// Calculate progress percentage
	//progressPercentage := int64(100.0 * float64(pr.TotalRead) / float64(pr.ContentLength))

	// Update the progress map safely
	progressMutex.Lock()
	downloadProgressMap["sse-progress"] = DownloadProgress{
		Total:   pr.ContentLength,
		Current: pr.TotalRead,
	}
	progressMutex.Unlock()

	pr.ProgressBar.Add(n)

	return n, err
}

func Download(url string, localPath string) error {
	dir := filepath.Dir(localPath)
	// Ensure the directory exists
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open the local file for writing, create it if not exists
	out, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer out.Close()

	// Find out how much has already been downloaded
	fi, err := out.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	size := fi.Size()

	// If already downloaded, no need to download again
	if size > 0 {
		fmt.Printf("Resuming download from byte %d...\n", size)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Range header to request the portion of the file we don't have yet
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", size))

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start file download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status getting file: %s", resp.Status)
	}

	// Seek to the end of the file to start appending data
	_, err = out.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	pterm.Info.Printf("Downloading model:\nURL: %s\nFile: %s\n", url, localPath)

	// Initialize the progress bar
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(int(resp.ContentLength)).WithTitle("Downloading").Start()

	// Wrap the response body in a custom reader that updates the progress bar
	progressReader := &ProgressReader{
		Reader:        resp.Body,
		ProgressBar:   progressBar,
		TotalRead:     size,
		ContentLength: resp.ContentLength + size,
	}

	// Copy the remaining data to the file, updating progress along the way
	_, err = io.Copy(out, progressReader)
	if err != nil {
		return err
	}

	// Ensure the progress bar reflects the complete download
	progressBar.Total = (int(progressReader.TotalRead))

	// Finish the progress bar
	progressBar.Stop()

	pterm.Success.Println("Download completed successfully.")

	return nil
}

func GetDownloadProgress(key string) string {
	progressMutex.Lock()
	defer progressMutex.Unlock()

	// Return the progress for a specific key
	if progress, ok := downloadProgressMap[key]; ok {
		// Calculate progress percentage
		return fmt.Sprintf("%d%%", int64(100.0*float64(progress.Current)/float64(progress.Total)))
	}

	// Fallback if no progress is found
	return "0%"
}

func (m *Model) Delete() error {
	if err := os.Remove(m.LocalPath); err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	return nil
}

func IncrementTurnCounter() {
	TurnCounter++
}

func GetTurnCounter() int {
	return TurnCounter
}
