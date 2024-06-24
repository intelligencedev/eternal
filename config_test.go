package main

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fs := afero.NewOsFs()

	// Load example .config.yml from current directory
	config, err := LoadConfig(fs, ".config.yml")
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, "User", config.CurrentUser)
	assert.Equal(t, "Assistant", config.AssistantName)
	assert.Equal(t, "localhost", config.ControlHost)
	assert.Equal(t, "8080", config.ControlPort)
	assert.Equal(t, "/Users/$USER/.eternal", config.DataPath)
}
