package user

import (
	"vid-msa/model"
	"vid-msa/serializer"

	"golang.org/x/crypto/bcrypt"
)

// ResetPasswordService 管理用户重新设置密码的服务
type ResetPasswordService struct {
	UserName    string `form:"username" json:"username" validate:"required"`
	Password    string `form:"password" json:"password" validate:"required"`
	NewPassword string `form:"new_password" json:"new_password" validate:"required"`
}

// ResetPassword 更新密码
func (u *ResetPasswordService) ResetPassword() *serializer.Response {
	user := &model.User{}
	// 查找用户
	if rdb := model.DB.Where("username = ?", u.UserName).First(&user); rdb.RowsAffected == 0 {
		return serializer.ParamErr("账号或密码错误", nil)
	} else if err := rdb.Error; err != nil {
		return serializer.DBErr("search user err", err)
	}
	// 检查密码
	if ok, err := user.MatchPassword(u.Password); err != nil {
		return serializer.EncryptErr("密码校验失败", err)
	} else if !ok {
		return serializer.ParamErr("账号或密码错误", nil)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.NewPassword), model.PasswordCost)
	if err != nil {
		return serializer.ParamErr("param err", err)
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
