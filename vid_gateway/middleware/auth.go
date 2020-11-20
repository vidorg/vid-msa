package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"gofiber-scaffold/client"
	pb "gofiber-scaffold/pb/user"
)

// CheckRight check role
func CheckRight() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := GetToken(c)
		if token == "" {
			return c.JSON(&fiber.Map{
				"code":    20001,
				"message": "请先登录",
			})
		}
		user, _ := client.UserClient.GetUser(context.Background(), &pb.GetUserRequest{
			Token: token,
		})
		if user.User == nil {
			return c.JSON(&fiber.Map{
				"code":    20001,
				"message": "请先登录",
			})
		}
		return c.Next()
	}
}
