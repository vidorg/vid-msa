package user

import (
	"strconv"
	"time"
	"vid-msa/middleware"
	"vid-msa/model"
	"vid-msa/pkg/cache"
	"vid-msa/pkg/gredis"
	"vid-msa/pkg/jwt"
	"vid-msa/serializer"

	"github.com/gofiber/fiber/v2"
)

// LogoutService 管理用户注销的服务
type LogoutService struct{}

// Logout 用户登出
func (u *LogoutService) Logout(c *fiber.Ctx) *serializer.Response {
	userInfo := c.Locals("user").(*model.User)
	// 设置 session
	// session.Destroy(c)
	// session.Delete(c, "user_id")
	// token
	token := middleware.GetToken(c)
	// 删除cookie
	c.Cookie(&fiber.Cookie{
		Name: "token",
		// Set expiry date to the past
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})
	// 缓存
	pattern := jwt.RedisTokenConcat(strconv.FormatUint(uint64(userInfo.ID), 10), token)
	err := gredis.Delete(pattern)
	if err != nil {
		return serializer.DBErr("logout err", err)
	}
	if err := cache.Delete(pattern); err != nil {
		return serializer.DBErr("logout err", err)
	}
	return &serializer.Response{
		Code: 200,
		Msg:  "注销成功",
	}
}
