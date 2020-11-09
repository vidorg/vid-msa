package middleware

import (
	"strconv"
	"vid-msa/pkg/gredis"
	"vid-msa/pkg/jwt"
	"vid-msa/serializer"

	"github.com/gofiber/fiber/v2"
)

// GetToken get token
func GetToken(c *fiber.Ctx) string {
	token := c.Cookies("token", "")
	if token == "" {
		token = c.Get("Authorization")
		if token != "" {
			return token
		}
		return c.Query("token", "")
	}
	return token
}

// JWT JWT middleware
func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := GetToken(c)
		if token == "" {
			return c.Next()
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			return c.Next()
		}
		pattern := jwt.RedisTokenConcat(strconv.FormatUint(uint64(claims.UID), 10), token)
		if !gredis.Exists(pattern) {
			res := serializer.LoginErr()
			return c.JSON(res)
		}
		//if !cache.Exists(pattern) {
		//	res := serializer.LoginExpiredErr()
		//	c.ClearCookie("token")
		//	return c.JSON(res)
		//}
		uid := claims.UID
		c.Locals("user_id", uid)
		return c.Next()
	}
}
