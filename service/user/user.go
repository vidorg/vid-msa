package user

import (
	"strconv"
	"vid-msa/model"
	"vid-msa/pkg/cache"
	"vid-msa/pkg/gredis"
	"vid-msa/pkg/jwt"
	"vid-msa/pkg/mail"
	"vid-msa/serializer"

	"github.com/gofiber/fiber/v2"
	"github.com/seefs001/seefslib-go/xstring"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	ID uint `form:"id" json:"id" validate:"required"`
}

// LoginByEmailService 邮箱验证码登录的服务
type LoginByEmailService struct {
	Email     string `form:"email" json:"email" validate:"required|email"`
	ValidCode string `form:"valid_code" json:"valid_code" validate:"required"`
}

// SendEmailService 邮箱验证码登录的服务
type SendEmailService struct {
	Email string `form:"email" json:"email" validate:"required|email"`
}

// NoParamService 无参数service
type NoParamService struct{}

// ForgetPasswordService 管理用户忘记密码的服务
type ForgetPasswordService struct {
	UserName    string `form:"username" json:"username" validate:"required"`
	ValidCode   string `form:"valid_code" json:"valid_code" validate:"required"`
	NewPassword string `form:"new_password" json:"new_password" validate:"required"`
}

// SendValidCode 发送六位验证码
func (u *SendEmailService) SendValidCode(validType string) *serializer.Response {
	code := xstring.GenValidateCode(6)
	validCode, err := xstring.StringToInt64(code)
	/**/ _ = gredis.Delete(validType + ":" + u.Email)
	err = gredis.Set(validType+":"+u.Email, validCode, 30*60)
	err = mail.SendMail(code, u.Email)
	if err != nil {
		return serializer.ParamErr("邮箱输入有误", nil)
	}
	return &serializer.Response{
		Code: 200,
		Msg:  "验证码发送成功",
	}
}

// Login 邮箱验证码登录
func (u *LoginByEmailService) Login(c *fiber.Ctx, validType string) *serializer.Response {
	user := &model.User{}

	// 查找用户
	if rdb := model.DB.Where("email = ?", u.Email).First(&user); rdb.RowsAffected == 0 {
		return serializer.ParamErr("没有该用户", nil)
	} else if err := rdb.Error; err != nil {
		return serializer.DBErr("查找用户出错", err)
	}

	// 检查redis的验证码
	validCode, err := gredis.Get(validType + ":" + u.Email)
	if err != nil {
		return serializer.ParamErr("您还没有发送验证码", err)
	}
	if string(validCode) != u.ValidCode {
		return serializer.ParamErr("验证码错误", err)
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
	// res := serializer.BuildUserResponse(user)
	// return res

	// JWT
	token, err := jwt.GenerateToken(int(user.ID))
	if err != nil {
		return serializer.EncryptErr("token令牌生成失败,请联系管理员", err)
	}

	// 保存到cookie
	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: token,
	})

	// 将token保存到redis
	pattern := jwt.RedisTokenConcat(strconv.FormatUint(uint64(user.ID), 10), token)
	//err = gredis.Set(pattern, user.ID, viper.GetInt("jwt.expire"))
	//if err != nil {
	//	return serializer.DBErr("登录缓存操作失败，请联系管理员", err)
	//}
	if err := cache.Set(pattern, strconv.Itoa(int(user.ID))); err != nil {
		return serializer.DBErr("登录缓存操作失败，请联系管理员", err)
	}
	return serializer.BuildLoginResponse(user, token)

}

// ResetPassword 忘记密码
func (u *ForgetPasswordService) ResetPassword(validType string) *serializer.Response {
	user := &model.User{}

	// 查找用户
	if rdb := model.DB.Where("username = ?", u.UserName).First(&user); rdb.RowsAffected == 0 {
		return serializer.ParamErr("账号或密码错误", nil)
	} else if err := rdb.Error; err != nil {
		return serializer.DBErr("查找用户出错", err)
	}
	if user.Email == nil {
		return serializer.DBErr("您没有设置邮箱，请联系管理员", nil)
	}
	validCode, err := gredis.Get(validType + ":" + *user.Email)
	if err != nil {
		return serializer.ParamErr("您还没有发送验证码", err)
	}
	if string(validCode[:]) != u.ValidCode {
		return serializer.ParamErr("验证码错误", err)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(u.NewPassword), model.PasswordCost)
	if err != nil {
		return serializer.ParamErr("参数错误", err)
	}

	update := model.DB.Model(&model.User{}).Where("username = ?", u.UserName).Update("password", string(bytes))
	if update == nil {
		return serializer.ParamErr("密码加密失败", err)
	}
	return &serializer.Response{
		Code: 200,
		Msg:  "密码修改成功",
	}
}

// GetCurrentUser 获取当前登录用户
func (u *NoParamService) GetCurrentUser(c *fiber.Ctx) *serializer.Response {
	user := c.Locals("user")
	if user == nil {
		return &serializer.Response{
			Code: 403,
			Msg:  "not found locals(user)",
		}
	}
	userInfo := user.(*model.User)
	if userInfo == nil {
		return &serializer.Response{
			Code: 403,
			Msg:  "获取用户信息失败，可能您登录已经过期",
		}
	}
	return serializer.BuildUserResponse(userInfo)
}

func (s *Service) GetUserInfo() *serializer.Response {
	user := model.User{}
	model.DB.Model(&model.User{}).Omit("Password").Where("id = ?", s.ID).First(&user)
	return &serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: user,
	}
}
