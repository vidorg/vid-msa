package user

import (
	"strconv"
	"time"
	"vid-msa/model"
	"vid-msa/pkg/gredis"
	"vid-msa/pkg/jwt"
	"vid-msa/serializer"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// LoginService 管理用户登录的服务
type LoginService struct {
	UserName string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

// Login 用户登录
func (u *LoginService) Login(c *fiber.Ctx) *serializer.Response {
	user := &model.User{}
	// 查找用户
	if rdb := model.DB.Where("username = ?", u.UserName).First(&user); rdb.RowsAffected == 0 {
		return serializer.ParamErr("账号或密码错误", nil)
	} else if err := rdb.Error; err != nil {
		return serializer.DBErr("查找用户出错", err)
	}
	// 检查密码
	ok, err := user.MatchPassword(u.Password)
	if err != nil || !ok {
		return serializer.EncryptErr("密码校验失败", err)
	}
	// 检查用户状态
	switch user.Status {
	case model.Inactive:
		return serializer.UserStatusErr("账号未被激活")
	case model.Suspend:
		return serializer.UserStatusErr("账号等待验证")
	default:
	}
	// 设置 session
	// session.Set(c, "user_id", user.ID)
	// c.Locals("user", user)
	// return serializer.BuildUserResponse(user)
	// JWT
	token, err := jwt.GenerateToken(int(user.ID))
	if err != nil {
		return serializer.EncryptErr("token令牌生成失败,请联系管理员", err)
	}
	// 保存到cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Second * time.Duration(viper.GetInt64("jwt.expire"))),
		Secure:   false,
		HTTPOnly: true,
		SameSite: "lax",
	})
	// 将token保存到redis
	pattern := jwt.RedisTokenConcat(strconv.FormatUint(uint64(user.ID), 10), token)
	if err := gredis.Set(pattern, user.ID, time.Duration(viper.GetInt("jwt.expire"))); err != nil {
		return serializer.DBErr("cache store err", err)
	}
	//if err := cache.Set(pattern, strconv.Itoa(int(user.ID))); err != nil {
	//	return serializer.DBErr("cache store err", err)
	//}
	return serializer.BuildLoginResponse(user, token)
}
