package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"gofiber-scaffold/client"
	pb "gofiber-scaffold/pb/user"
)

func Test1(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{
		"msg": "xxx",
	})
}
func Test2(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{
		"msg": "xxx2",
	})
}

func PhoneLogin(c *fiber.Ctx) error {
	phone := c.FormValue("phone", "1231424")
	password := c.FormValue("password", "1231424")
	response, err := client.UserClient.PhoneLogin(context.Background(),
		&pb.PhoneLoginRequest{
			Phone:    phone,
			Password: password,
		})
	if err != nil {
		return err
	}
	return c.JSON(response)
}
