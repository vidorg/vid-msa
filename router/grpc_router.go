package router

import (
	"vid-msa/middleware"
	"vid-msa/router/controller"

	"github.com/gofiber/fiber/v2"
)

func InitGRPCRouter(app *fiber.App) {
	// Auth Required
	auth := app.Group("/api/v1", middleware.CheckRight())
	auth.Get("/grpc/add/:a/:b", controller.Add)
	auth.Get("/grpc/multi/:a/:b", controller.Multi)
}
