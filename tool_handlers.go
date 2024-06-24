package main

import (
	"eternal/pkg/llm/openai"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pterm/pterm"
)

// handleToolToggle toggles the state of various tools based on the provided tool name.
func handleToolToggle(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		toolName := c.Params("toolName")
		enabled := c.Params("enabled")
		topN := c.Params("topN")

		pterm.Info.Println(enabled)

		// Convert the enabled parameter to a boolean.
		enabledBool, err := strconv.ParseBool(enabled)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid enabled parameter")
		}

		// Convert the topN parameter to an integer.
		topNInt, err := strconv.Atoi(topN)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid topN parameter")
		}

		// Print the params to the console.
		pterm.Info.Println("Params:")
		pterm.Info.Println(toolName)

		switch toolName {
		case "memory":
			pterm.Warning.Sprintf("Memory tool toggled: %t\n", config.Tools.Memory.Enabled)
			config.Tools.Memory.Enabled = enabledBool
			config.Tools.Memory.TopN = topNInt
		case "webget":
			pterm.Warning.Sprintf("WebGet tool toggled: %t\n", config.Tools.WebGet.Enabled)
			config.Tools.WebGet.Enabled = !config.Tools.WebGet.Enabled
		case "websearch":
			pterm.Warning.Sprintf("WebSearch tool toggled: %t\n", config.Tools.WebSearch.Enabled)
			config.Tools.WebSearch.Enabled = enabledBool
			config.Tools.WebSearch.TopN = topNInt
		case "imggen":
			config.Tools.ImgGen.Enabled = true
		default:
			return c.Status(fiber.StatusNotFound).SendString("Tool not found")
		}

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Tool %s toggled", toolName)})
	}
}

// handleToolList retrieves and returns a list of tools from the configuration with all parameters.
func handleToolList(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(config.Tools)
	}
}

// handleOpenAIModels retrieves and returns a list of OpenAI models.
func handleOpenAIModels(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		client := openai.NewClient(config.OAIKey)
		modelsResponse, err := openai.GetModels(client)

		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

		var gptModels []string
		for _, model := range modelsResponse.Data {
			if strings.HasPrefix(model.ID, "gpt") {
				gptModels = append(gptModels, model.ID)
			}
		}

		return c.JSON(fiber.Map{
			"object": "list",
			"data":   gptModels,
		})
	}
}

// handleRoleSelection handles the selection of assistant roles.
func handleRoleSelection(config *AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName := c.Params("name")
		var foundRole *struct {
			Name         string `yaml:"name"`
			Instructions string `yaml:"instructions"`
		}

		for i := range config.AssistantRoles {
			if config.AssistantRoles[i].Name == roleName {
				foundRole = &config.AssistantRoles[i]
				break
			}
		}

		if foundRole == nil {
			pterm.Warning.Printf("Role %s not found. Defaulting to 'chat'.\n", roleName)
			for i := range config.AssistantRoles {
				if config.AssistantRoles[i].Name == "chat" {
					foundRole = &config.AssistantRoles[i]
					break
				}
			}
		}

		if foundRole == nil && len(config.AssistantRoles) > 0 {
			foundRole = &config.AssistantRoles[0]
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("Role set to %s", foundRole.Name),
			})
		}

		if foundRole != nil {
			config.CurrentRoleInstructions = foundRole.Instructions
			pterm.Info.Printf("Role set to: %s\n", foundRole.Name)
			pterm.Info.Println(foundRole.Instructions)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("Role set to %s", foundRole.Name),
			})
		}

		return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
	}
}
