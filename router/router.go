package router

import (
	"github.com/gofiber/fiber/v2"
)

// InitRouter initialize router
func InitRouter(app *fiber.App) *fiber.App {
	// assets
	app.Static("/file", "./static")
	// views
	app.Get("/view", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})
	app.Get("/layout", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main")
	})
	// ping
	app.All("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"code": 200,
			"msg":  "success",
		})
	})
	// init user router
	InitUserRouter(app)
	// init file router
	InitFileRouter(app)
	// init grpc router
	InitGRPCRouter(app)

	return app
}
