package router

import (
	"github.com/gofiber/fiber/v2"
)

// InitRouter initialize router
func InitRouter(app *fiber.App) *fiber.App {
	// ping
	app.All("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"code": 200,
			"msg":  "success",
		})
	})
	// init user router
	InitUserRouter(app)

	return app
}
