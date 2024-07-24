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

	// Project routes
	app.Post("/v1/projects", handleGetProjects())

	// Chat session routes
	app.Post("/v1/chat/submit", handleChatSubmit(config))
	app.Post("/v1/chat/role/:name", handleRoleSelection(config))

	// app.Post("/stop-streaming/:turnID", handleStopStreaming()) //Not currently implemented

	// Model management routes
	app.Post("/v1/views/models", handleModelCards(modelParams))
	app.Post("/v1/models/select/:name/:action", handleModelSelect())
	app.Get("/v1/models/selected", handleSelectedModels())
	app.Post("/v1/models/download", handleModelDownload(config))
	app.Post("/v1/models/img/download", handleImgModelDownload(config))
	app.Post("/v1/models/set", handleModelUpdate())

	// Roles routes
	app.Get("/v1/roles", handleGetRoles(config))

	// Model - Database routes
	app.Get("/v1/modeldata/:modelName", handleModelData())
	app.Put("/v1/modeldata/:modelName/downloaded", handleModelDownloadUpdate())

	// Chat - Database routes
	app.Get("/v1/chats", handleGetChats())
	app.Get("/v1/chats/:id", handleGetChatByID())
	app.Put("/v1/chats/:id", handleUpdateChat())
	app.Delete("/v1/chats/:id", handleDeleteChat())

	// Tool routes
	app.Get("/v1/tools", handleRenderTools(config))
	app.Get("/v1/tools/list", handleToolList(config))
	app.Post("/v1/tools/:toolName/:enabled/:topN", handleToolToggle(config))
	app.Get("/v1/dpsearch", handleDPSearch())
	app.Post("/v1/tools/img/workflow/set", handleImgSetWorkflow(config))
	app.Post("/v1/tools/img/resolution/set", handleImgSetResolution(config))

	// Utility routes
	app.Get("/v1/config", func(c *fiber.Ctx) error {
		return c.JSON(config)
	})
	app.Post("/upload", handleUpload(config))
	app.Get("/sseupdates", handleSSEUpdates())
	app.Get("/ws", websocket.New(handleWebSocket(config)))

	// Experimental routes
	// app.Get("/flow", func(c *fiber.Ctx) error {
	// 	return c.Render("templates/flow", fiber.Map{})
	// })
}
