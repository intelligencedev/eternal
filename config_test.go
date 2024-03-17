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

	// Refactored assertServiceHost function
	assertServiceHost := func(service string, hostKey string, expectedHost string, expectedPort string) {
		hostConfig, exists := config.ServiceHosts[service][hostKey]
		assert.True(t, exists)
		assert.Equal(t, expectedHost, hostConfig.Host) // Use Host instead of HostURL
		assert.Equal(t, expectedPort, hostConfig.Port) // Use Port instead of HostPort
	}

	// Updated calls to assertServiceHost with correct field names and types
	assertServiceHost("retrieval", "retrieval_host_1", "localhost", "8081")
	assertServiceHost("image", "image_host_1", "localhost", "8082")
	assertServiceHost("speech", "speech_host_1", "localhost", "8083")
	assertServiceHost("llm", "llm_host_1", "localhost", "8081")

	assert.Equal(t, "sk-...", config.OAIKey)
}
