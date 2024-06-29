package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// handleGetProjects returns a handler function that retrieves all projects
func handleGetProjects() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		projects, err := sqliteDB.GetProjects()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("failed to get projects: %v", err),
			})
		}

		// Print the projects
		for _, project := range projects {
			fmt.Printf("Project: %s\n", project.Name)
		}

		// render content in projects template
		return c.Render("templates/projects", fiber.Map{
			"projects": projects,
		})
	}
}
