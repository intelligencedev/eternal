package main

import (
	"bytes"
	"encoding/json"
	"eternal/pkg/llm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestIndexRoute(t *testing.T) {
	app := fiber.New()

	// Set up routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Send a GET request to the root URL
	req, _ := http.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)    // Read the response body
	assert.Equal(t, "OK", string(body)) // Convert body to string and compare with "OK"
}

func TestConfigRoute(t *testing.T) {
	app := fiber.New()

	// Set up routes
	app.Get("/config", func(c *fiber.Ctx) error {
		config := &AppConfig{
			CurrentUser:   "test_user",
			AssistantName: "test_assistant",
		}
		return c.JSON(config)
	})

	// Send a GET request to the /config URL
	req, _ := http.NewRequest("GET", "/config", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var config AppConfig
	err = json.NewDecoder(resp.Body).Decode(&config)
	assert.NoError(t, err)
	assert.Equal(t, "test_user", config.CurrentUser)
	assert.Equal(t, "test_assistant", config.AssistantName)
}

func TestToolRoute(t *testing.T) {
	app := fiber.New()

	// Set up routes
	app.Post("/tool/:toolName", func(c *fiber.Ctx) error {
		toolName := c.Params("toolName")
		var index int
		found := false
		for i, t := range tools {
			if t.Name == toolName {
				index = i
				found = true
				break
			}
		}
		if !found {
			return c.Status(404).SendString("Tool not found")
		}
		tools[index].Enabled = !tools[index].Enabled
		return c.JSON(tools[index])
	})

	// Set up test data
	tools = []Tool{
		{Name: "websearch", Enabled: false},
		{Name: "imagegen", Enabled: false},
	}

	// Test enabling a tool
	body := []byte{}
	req, _ := http.NewRequest("POST", "/tool/websearch", bytes.NewBuffer(body))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var tool Tool
	err = json.NewDecoder(resp.Body).Decode(&tool)
	assert.NoError(t, err)
	assert.Equal(t, "websearch", tool.Name)
	assert.True(t, tool.Enabled)

	// Test disabling a tool
	req, _ = http.NewRequest("POST", "/tool/websearch", bytes.NewBuffer(body))
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&tool)
	assert.NoError(t, err)
	assert.Equal(t, "websearch", tool.Name)
	assert.False(t, tool.Enabled)

	// Test non-existent tool
	req, _ = http.NewRequest("POST", "/tool/nonexistent", bytes.NewBuffer(body))
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Equal(t, "Tool not found", string(body))
}

func TestModelDataRoute(t *testing.T) {
	app := fiber.New()

	app.Get("/modeldata/:modelName", func(c *fiber.Ctx) error {
		mockData := ModelParams{
			ID:         3,
			Name:       "eternal-120b",
			Homepage:   "https://huggingface.co/intelligence-dev/eternal-120b",
			GGUFInfo:   "https://huggingface.co/intelligence-dev/eternal-120b",
			Downloads:  "",
			Downloaded: true,
			Options: &llm.GGUFOptions{
				Prompt: "Test",
			},
		}

		return c.JSON(mockData)
	})

	req := httptest.NewRequest("GET", "/modeldata/test-model", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var modelData ModelParams
	err = json.NewDecoder(resp.Body).Decode(&modelData)

	assert.NoError(t, err)
	assert.Equal(t, "eternal-120b", modelData.Name)
}
