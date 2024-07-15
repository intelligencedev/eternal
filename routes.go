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
	app.Post("/stop-streaming/:turnID", handleStopStreaming())

	// Model management routes
	app.Post("/modelcards", handleModelCards(modelParams))
	app.Post("/model/select/:name/:action", handleModelSelect())
	app.Get("/model/selected", handleSelectedModels())
	app.Post("/model/download", handleModelDownload(config))
	app.Post("/imgmodel/download", handleImgModelDownload(config))
	app.Post("/model/set/params", handleModelUpdate())

	// Roles routes
	app.Get("/roles", handleGetRoles(config))

	// Model - Database routes
	app.Get("/modeldata/:modelName", handleModelData())
	app.Put("/modeldata/:modelName/downloaded", handleModelDownloadUpdate())

	// Chat - Database routes
	app.Get("/chats", handleGetChats())
	app.Get("/chats/:id", handleGetChatByID())
	app.Put("/chats/:id", handleUpdateChat())
	app.Delete("/chats/:id", handleDeleteChat())

	// Tool routes
	app.Get("/tools", handleRenderTools(config))
	app.Get("/tools/list", handleToolList(config))
	app.Post("/tool/:toolName/:enabled/:topN", handleToolToggle(config))
	app.Get("/dpsearch", handleDPSearch())
	app.Post("/tools/img/workflow/set", handleImgSetWorkflow(config))
	app.Post("/tools/img/resolution/set", handleImgSetResolution(config))

	// Utility routes
	app.Post("/config", func(c *fiber.Ctx) error {
		return c.JSON(config)
	})
	app.Post("/upload", handleUpload(config))
	app.Get("/sseupdates", handleSSEUpdates())
	app.Get("/ws", websocket.New(handleWebSocket(config)))
}
