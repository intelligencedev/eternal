package main

import (
	"eternal/pkg/llm"
	"eternal/pkg/sd"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type AppConfig struct {
	ServerID       string                            `yaml:"server_id"`
	CurrentUser    string                            `yaml:"current_user"`
	AssistantName  string                            `yaml:"assistant_name"`
	ControlHost    string                            `yaml:"control_host"`
	ControlPort    string                            `yaml:"control_port"`
	DataPath       string                            `yaml:"data_path"`
	ServiceHosts   map[string]map[string]BackendHost `yaml:"service_hosts"`
	ChromedpKey    string                            `yaml:"chromedp_key"`
	OAIKey         string                            `yaml:"oai_key"`
	AnthropicKey   string                            `yaml:"anthropic_key"`
	GoogleKey      string                            `yaml:"google_key"`
	LanguageModels []llm.Model                       `yaml:"language_models"`
	ImageModels    []sd.ImageModel                   `yaml:"image_models"`
	AssistantRoles []struct {
		Name         string `yaml:"name"`
		Instructions string `yaml:"instructions"`
	} `yaml:"assistant_roles"`
	Tools struct {
		Memory struct {
			Enabled bool `yaml:"enabled"`
			TopN    int  `yaml:"top_n"`
		} `yaml:"memory"`
		WebGet struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"webget"`
		WebSearch struct {
			Enabled bool `yaml:"enabled"`
			TopN    int  `yaml:"top_n"`
		} `yaml:"websearch"`
		ImgGen struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"img_gen"`
	} `yaml:"tools"`
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

// LoadConfig loads configuration from a YAML file.
func LoadConfig(fs afero.Fs, path string) (*AppConfig, error) {
	config := &AppConfig{}

	// Use Afero to read the file
	file, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
