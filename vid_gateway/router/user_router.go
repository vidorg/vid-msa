package router

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-scaffold/middleware"
	"gofiber-scaffold/router/controller"
)

func InitUserRouter(app *fiber.App) {
	app.Get("/test1", controller.Test1)
	app.Post("/auth/login", controller.PhoneLogin)
	auth := app.Use(middleware.CheckRight())
	auth.Get("/test2", controller.Test2)
}
