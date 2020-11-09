package user

import (
	"vid-msa/model"
	"vid-msa/serializer"
)

// RegisterService 管理用户注册的服务
type RegisterService struct {
	UserName string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	NickName string `form:"nickname" json:"nickname"`
	//Email    string `form:"email" json:"email" validate:"required|email"`
	Email *string `form:"email" json:"email" validate:"email"`
}

// Register 用户注册
func (u *RegisterService) Register() *serializer.Response {
	// 如果没输入nickname，那和username同名
	if u.NickName == "" {
		u.NickName = u.UserName
	}
	user := &model.User{
		UserName: u.UserName,
		Nickname: u.NickName,
		Email:    u.Email,
		Status:   model.Active,
		Role:     "normal",
	}
	// 表单验证
	var count int64 = 0
	model.DB.Model(&model.User{}).Where("username = ?", u.UserName).Count(&count)
	if count > 0 {
		count = 0
		return serializer.ParamErr("用户名已经注册", nil)
	}
	if u.Email != nil {
		model.DB.Model(&model.User{}).Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			return serializer.DBErr("邮箱已经被注册", nil)
		}
	}
	// 加密密码
	if err := user.SetPassword(u.Password); err != nil {
		return serializer.EncryptErr("encrypt err", err)
	}
	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.DBErr("db err", err)
	}
	return serializer.BuildUserResponse(user)
}
