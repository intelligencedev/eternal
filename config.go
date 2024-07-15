// config.go
package main

import (
	"embed"
	"eternal/pkg/llm"
	"eternal/pkg/sd"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var (
	sqliteDB       *SQLiteDB
	searchIndex    bleve.Index
	localFs        = afero.NewOsFs()
	chatTurn       = 1
	messageCounter int64
)

// AppConfig holds the application configuration.
type AppConfig struct {
	ServerID        string                            `yaml:"server_id"`
	CurrentUser     string                            `yaml:"current_user"`
	AssistantName   string                            `yaml:"assistant_name"`
	ControlHost     string                            `yaml:"control_host"`
	ControlPort     string                            `yaml:"control_port"`
	DataPath        string                            `yaml:"data_path"`
	ServiceHosts    map[string]map[string]BackendHost `yaml:"service_hosts"`
	ChromedpKey     string                            `yaml:"chromedp_key"`
	OAIKey          string                            `yaml:"oai_key"`
	AnthropicKey    string                            `yaml:"anthropic_key"`
	GoogleKey       string                            `yaml:"google_key"`
	LanguageModels  []llm.Model                       `yaml:"language_models"`
	ImageModels     []sd.ImageModel                   `yaml:"image_models"`
	EmbeddingModel  string                            `yaml:"embedding_model"`
	CurrentRoleName string                            `yaml:"current_role_name"`
	//CurrentImgRoleName      string                            `yaml:"current_img_role_name"`
	CurrentRoleInstructions string `yaml:"current_role"`
	AssistantRoles          []struct {
		Name         string `yaml:"name"`
		Instructions string `yaml:"instructions"`
	} `yaml:"assistant_roles"`
	Tools                Tools        `yaml:"tools"`
	DefaultProjectConfig Project      `yaml:"default_project"`
	Projects             []Project    `yaml:"projects"`
	ImgGenConfig         ImgGenConfig `yaml:"img_gen_config"`
}

// BackendHost represents a local or remote backend host.
type BackendHost struct {
	ID            uint           `gorm:"primaryKey" yaml:"-"`
	Host          string         `yaml:"host" gorm:"column:host"`
	Port          string         `yaml:"port" gorm:"column:port"`
	GgufGPULayers int            `yaml:"gpu_layers" gorm:"column:gguf_gpu_layers"`
	ModelType     string         `yaml:"model_type" gorm:"column:model_type"`
	CreatedAt     time.Time      `yaml:"-"`
	UpdatedAt     time.Time      `yaml:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" yaml:"-"`
}

type ImgGenConfig struct {
	CheckpointName string `json:"model_name"`
	Workflow       string `json:"workflow"`
	Seed           int    `json:"seed"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
}

// WebSocketMessage represents a message sent over WebSocket.
type WebSocketMessage struct {
	ChatMessage string                 `json:"chat_message"`
	Model       string                 `json:"model"`
	Headers     map[string]interface{} `json:"HEADERS"`
}

// loadConfig loads configuration from a YAML file.
//
// Parameters:
// - fs: The filesystem to read the configuration file from.
// - path: The path to the configuration file.
//
// Returns:
// - A pointer to the loaded AppConfig.
// - An error if the configuration file could not be read or unmarshaled.
func loadConfig(fs afero.Fs, path string) (*AppConfig, error) {
	config := &AppConfig{}
	file, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// createModelParam creates a ModelParams object from a llm.Model and data path.
//
// Parameters:
// - model: The llm.Model to create parameters for.
// - dataPath: The path to the data directory.
//
// Returns:
// - A ModelParams object populated with the model's parameters.
func createModelParam(model *llm.Model, dataPath string) ModelParams {
	var localPath string
	if model.Downloads != nil {
		fileName := filepath.Base(model.Downloads[0])
		localPath = filepath.Join(dataPath, "models", model.Name, fileName)
	}

	downloaded := fileExists(localPath)

	return ModelParams{
		Name:       model.Name,
		Homepage:   model.Homepage,
		GGUFInfo:   model.GGUF,
		Downloaded: downloaded,
		Options: &llm.GGUFOptions{
			Model:         localPath,
			Prompt:        model.Prompt,
			CtxSize:       model.Ctx,
			Temp:          0.7,
			RepeatPenalty: 1.1,
		},
	}
}

// EnsureDataPath ensures that the data path exists.
//
// Parameters:
// - config: The application configuration containing the data path.
//
// Returns:
// - An error if the data path could not be created.
func EnsureDataPath(config *AppConfig) error {
	if _, err := os.Stat(config.DataPath); os.IsNotExist(err) {
		return localFs.MkdirAll(config.DataPath, os.ModePerm)
	}
	return nil
}

// CopyFiles copies files from an embedded filesystem to a destination directory.
//
// Parameters:
// - fsys: The embedded filesystem to copy files from.
// - srcDir: The source directory in the embedded filesystem.
// - destDir: The destination directory on the local filesystem.
//
// Returns:
// - An error if the files could not be copied.
func CopyFiles(fsys embed.FS, srcDir, destDir string) error {
	fileEntries, err := fsys.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", srcDir, err)
	}

	for _, entry := range fileEntries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			if err := CopyFiles(fsys, srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(fsys, srcPath, destPath); err != nil {
				log.Errorf("Failed to copy file %s: %v", srcPath, err)
			}
		}
	}
	return nil
}

// copyFile copies a single file from the embedded filesystem to the destination.
//
// Parameters:
// - fsys: The embedded filesystem to copy the file from.
// - srcPath: The source file path in the embedded filesystem.
// - destPath: The destination file path on the local filesystem.
//
// Returns:
// - An error if the file could not be copied.
func copyFile(fsys embed.FS, srcPath, destPath string) error {
	fileData, err := fsys.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", srcPath, err)
	}
	return os.WriteFile(destPath, fileData, 0755)
}

// IncrementTurn increments and returns a counter for frontend chat elements.
//
// Returns:
// - The incremented counter value.
func IncrementTurn() int64 {
	return atomic.AddInt64(&messageCounter, 1)
}
