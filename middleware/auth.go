package middleware

import (
	"vid-msa/model"
	"vid-msa/pkg/casbin"
	"vid-msa/pkg/logger"
	"vid-msa/serializer"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetUser get local user
func GetUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// // session
		// uid := session.Get(c, "user_id")
		// if uid == nil {
		// 	return c.Next()
		// }
		//
		// user, err := model.GetUser(uid)
		// if err == nil {
		// 	c.Locals("user", user)
		// 	return c.Next()
		// }
		//
		// c.Locals("user", nil)
		// jwt
		uid := c.Locals("user_id")
		if uid == nil {
			return c.Next()
		}
		user, err := model.GetUser(uid)
		if err == nil {
			c.Locals("user", user)
			return c.Next()
		}
		return nil
	}
}

// CheckRight check role
func CheckRight() fiber.Handler {
	adapter, err := gormadapter.NewAdapterByDB(model.DB)

	if err != nil {
		logger.Error("[RBAC]",
			zap.Error(err))
	}

	return func(c *fiber.Ctx) error {
		userVal := c.Locals("user")
		if userVal == nil {
			res := serializer.LoginErr()
			return c.JSON(res)
		}
		user, ok := userVal.(*model.User)
		if !ok {
			res := serializer.LoginErr()
			return c.JSON(res)
		}
		ok, err := casbin.Enforce(adapter, user.Role, c)
		if err != nil {
			res := serializer.ServerErr("check role err", err)
			return c.JSON(res)
		} else if !ok {
			res := serializer.NoRightErr()
			return c.JSON(res)
		}
		return c.Next()
	}
}
