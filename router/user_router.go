package router

import (
	"vid-msa/middleware"
	"vid-msa/router/controller"

	"github.com/gofiber/fiber/v2"
)

func InitUserRouter(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Post("/UserLogin", controller.Login)
	v1.Post("/UserRegister", controller.Register)
	v1.Post("/ForgetPassword", controller.ResetPasswordByValidEmail)
	v1.Post("/ResetPassword", controller.ResetPassword)
	v1.Get("/GetUserList", controller.GetUserList)
	// Auth Required
	auth := app.Group("/api/v1", middleware.CheckRight())
	auth.Post("/UserLogout", controller.Logout)
	auth.Get("/GetLocalUser", controller.GetUserInfo)
	auth.Put("/ChangeUserInfo", controller.ChangeInfo)
}
