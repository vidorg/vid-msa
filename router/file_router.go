package router

import (
	"vid-msa/middleware"
	"vid-msa/router/controller"

	"github.com/gofiber/fiber/v2"
)

func InitFileRouter(app *fiber.App) {
	// Auth Required
	auth := app.Group("/api/v1", middleware.CheckRight())
	auth.Get("/GetQiniuToken", controller.GetUploadToken)
	auth.Post("/UploadFile", controller.UploadFile)
}
