package user

import (
	"vid-msa/model"
	"vid-msa/serializer"

	"github.com/gofiber/fiber/v2"
)

// ChangeUserInfoService 修改用户资料
type ChangeUserInfoService struct {
	UserName string  `form:"username" json:"username"`
	Avatar   string  `form:"avatar" json:"avatar"`
	NickName string  `form:"nickname" json:"nickname"`
	Email    *string `form:"email" json:"email"`
}

// ChangeUserInfo 修改用户资料
func (s *ChangeUserInfoService) ChangeUserInfo(c *fiber.Ctx) *serializer.Response {
	//userID := session.Get(c, "user_id").(int64)
	userID := c.Locals("user_id")
	model.DB.Model(&model.User{}).Where("id = ? ", userID).Updates(&model.User{
		UserName: s.UserName,
		Nickname: s.NickName,
		Avatar:   s.Avatar,
		Email:    s.Email,
	})
	return &serializer.Response{
		Code: 200,
		Msg:  "更新用户信息成功",
	}
}
