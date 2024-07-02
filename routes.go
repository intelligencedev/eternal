// eternal/routes.go - API routes

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// setupRoutes sets up the routes for the application
func setupRoutes(app *fiber.App, config *AppConfig, modelParams []ModelParams) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("templates/index", fiber.Map{})
	})

	app.Get("/config", func(c *fiber.Ctx) error {
		return c.JSON(config)
	})

	app.Get("/flow", func(c *fiber.Ctx) error {
		return c.Render("templates/flow", fiber.Map{})
	})

	// Project routes
	app.Post("/projects", handleGetProjects())

	// Chat session routes
	app.Post("/chatsubmit", handleChatSubmit(config))
	app.Post("/chat/role/:name", handleRoleSelection(config))

	// Model management routes
	app.Post("/modelcards", handleModelCards(modelParams))
	app.Post("/model/select/:name/:action", handleModelSelect())
	app.Get("/model/selected", handleSelectedModels())
	app.Post("/model/download", handleModelDownload(config))
	app.Post("/imgmodel/download", handleImgModelDownload(config))
	app.Post("/model/set/params", handleModelUpdate())

	// Model - Database routes
	app.Get("/modeldata/:modelName", handleModelData())
	app.Put("/modeldata/:modelName/downloaded", handleModelDownloadUpdate())

	// Chat - Database routes
	app.Get("/chats", handleGetChats())
	app.Get("/chats/:id", handleGetChatByID())
	app.Put("/chats/:id", handleUpdateChat())
	app.Delete("/chats/:id", handleDeleteChat())

	// Tool routes
	app.Get("/tools/list", handleToolList(config))
	app.Post("/tool/:toolName/:enabled/:topN", handleToolToggle(config))
	app.Get("/dpsearch", handleDPSearch())

	// Utility routes
	// return the app config
	app.Post("/config", func(c *fiber.Ctx) error {
		return c.JSON(config)
	})
	app.Post("/upload", handleUpload(config))
	app.Get("/sseupdates", handleSSEUpdates())
	app.Get("/ws", websocket.New(handleWebSocket(config)))

	// OpenAI routes
	app.Get("/openai/models", handleOpenAIModels(config))
	//app.Get("/wsoai", websocket.New(handleOpenAIWebSocket(config)))

	// Anthropic routes
	//app.Get("/wsanthropic", websocket.New(handleAnthropicWebSocket(config)))

	// Google routes
	// app.Get("/wsgoogle", websocket.New(handleGoogleWebSocket(config)))
}
