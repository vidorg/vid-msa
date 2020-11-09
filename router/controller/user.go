package controller

import (
	"vid-msa/serializer"
	"vid-msa/service/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

// GetUserList get user list
func GetUserList(c *fiber.Ctx) error {
	service := &user.QueryUsersService{}
	if err := c.QueryParser(service); err != nil {
		return c.JSON(serializer.ParamErr("please input page and limit", err))
	}
	return c.JSON(service.QueryUsers())
}

// GetUserInfo get local user
func GetUserInfo(c *fiber.Ctx) error {
	service := &user.NoParamService{}

	return c.JSON(service.GetCurrentUser(c))
}

// Logout logout user
func Logout(c *fiber.Ctx) error {
	service := &user.LogoutService{}

	return c.JSON(service.Logout(c))
}

// ChangeInfo update user profile
func ChangeInfo(c *fiber.Ctx) error {
	service := &user.ChangeUserInfoService{}
	if err := c.BodyParser(service); err != nil {
		return c.JSON(serializer.ParamErr("请检查是否输入正确", err))
	}
	return c.JSON(service.ChangeUserInfo(c))
}

// Login user login
func Login(c *fiber.Ctx) error {
	service := &user.LoginService{}
	if err := c.BodyParser(service); err != nil {
		return c.JSON(serializer.ParamErr("请检查是否输入正确", err))
	}
	v := validate.Struct(service)
	if v.Validate() {
		return c.JSON(service.Login(c))
	}
	return c.JSON(serializer.ParamErr(v.Errors.One(), v.Errors))
}

// ResetPassword reset password
func ResetPassword(c *fiber.Ctx) error {
	service := &user.ResetPasswordService{}
	if err := c.BodyParser(service); err != nil {
		return c.JSON(serializer.ParamErr("请检查是否输入正确", err))
	}
	v := validate.Struct(service)
	if v.Validate() {
		return c.JSON(service.ResetPassword())
	}
	return c.JSON(serializer.ParamErr(v.Errors.One(), v.Errors))
}

// ResetPasswordByValidEmail forget password
func ResetPasswordByValidEmail(c *fiber.Ctx) error {
	service := &user.ForgetPasswordService{}
	if err := c.BodyParser(service); err != nil {
		return c.JSON(serializer.ParamErr("请检查是否输入正确", err))
	}
	v := validate.Struct(service)
	if v.Validate() {
		return c.JSON(service.ResetPassword("forget_password"))
	}
	return c.JSON(serializer.ParamErr(v.Errors.One(), v.Errors))
}

// Register user register
func Register(c *fiber.Ctx) error {
	service := &user.RegisterService{}
	if err := c.BodyParser(service); err != nil {
		return c.JSON(serializer.ParamErr("请检查是否输入正确", err))
	}
	v := validate.Struct(service)
	if v.Validate() {
		return c.JSON(service.Register())
	}
	return c.JSON(serializer.ParamErr(v.Errors.One(), v.Errors))
}
