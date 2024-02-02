package main

import (
	"eternal/pkg/llm"
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
	LanguageModels []llm.Model                       `yaml:"language_models"`
}

type BackendHost struct {
	ID        uint           `gorm:"primaryKey" yaml:"-"`
	Host      string         `yaml:"host" gorm:"column:host"`
	Port      string         `yaml:"port" gorm:"column:port"`
	ModelType string         `yaml:"model_type" gorm:"column:model_type"`
	CreatedAt time.Time      `yaml:"-"`
	UpdatedAt time.Time      `yaml:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" yaml:"-"`
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
