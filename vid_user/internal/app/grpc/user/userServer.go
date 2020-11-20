package user

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"vid_user/internal/app/gredis"
	"vid_user/internal/app/jwt"
	"vid_user/internal/app/model"
	"vid_user/internal/app/uuid"

	"vid_user/pb/user"
)

const (
	PHONE = "phone"
	EMAIL = "email"
)

type UserServer struct{}

func (server *UserServer) GetUserID(ctx context.Context, request *user.GetUserIDRequest) (*user.Response, error) {
	panic("xxx")
}

func (server *UserServer) SetUserInfo(ctx context.Context, request *user.SetUserInfoRequest) (response *user.Response, err error) {
	if model.DB.Where("id = ?", request.Uid).
		Updates(&model.User{
			NickName: request.Nickname,
			Gender:   request.Gender,
		}).RowsAffected == 0 {
		response = &user.Response{
			Code:    20006,
			Message: "更新用户信息失败",
		}
		return
	}
	response = &user.Response{
		Code:    200,
		Message: "更新成功",
	}
	return
}

func (server *UserServer) GetUser(context context.Context, request *user.GetUserRequest) (response *user.Response, err error) {
	claims, err := jwt.ParseToken(request.Token)
	if err != nil {
		response = &user.Response{
			Code:    20001,
			Message: "提取token出错",
		}
		return
	}
	if !gredis.CheckToken(request.Token, claims.UID) {
		response = &user.Response{
			Code:    20006,
			Message: "你还没有登录",
		}
		return
	}
	var userInfo model.User
	result := model.DB.Where("id = ?", claims.UID).
		First(&userInfo)
	if result.Error != nil {
		response = &user.Response{
			Code:    20004,
			Message: "当前用户不存在",
		}
		return
	}
	response = &user.Response{
		Code:    200,
		Message: "success",
		User: &user.UserInfo{
			Id:       userInfo.ID,
			Nickname: userInfo.NickName,
			Avatar:   userInfo.Avatar,
			Gender:   userInfo.Gender,
		},
	}
	return
}

func (server *UserServer) PhoneLogin(context context.Context, request *user.PhoneLoginRequest) (response *user.Response, err error) {
	var userAuthorization model.UserAuthorization
	result := model.DB.Model(&model.UserAuthorization{}).
		Where("identity_type = ?", PHONE).
		Where("identifier = ?", request.Phone).
		Where("credential = ?", request.Password).
		First(&userAuthorization)
	if result.Error != nil {
		response = &user.Response{
			Code:    20001,
			Message: "手机号或者密码错误",
		}
		return
	}
	token, err := jwt.GenerateToken(userAuthorization.UserID)
	if err != nil {
		fmt.Println(err)
		response = &user.Response{
			Code:    500,
			Message: "生成token出错",
		}
		return
	}
	if !gredis.SetToken(token, userAuthorization.UserID) {
		response = &user.Response{
			Code:    500,
			Message: "保存到redis出错",
		}
		return
	}
	response = &user.Response{
		Code:    200,
		Message: "success",
		Data: &user.TokenData{
			Token: token,
		},
	}
	return
}

func (server *UserServer) EmailLogin(context context.Context, request *user.EmailLoginRequest) (response *user.Response, err error) {
	var userAuthorization model.UserAuthorization
	result := model.DB.Model(&model.UserAuthorization{}).
		Where("identity_type = ?", EMAIL).
		Where("identifier = ?", request.Email).
		Where("credential = ?", request.Password).
		First(&userAuthorization)
	if result.Error != nil {
		response = &user.Response{
			Code:    20002,
			Message: "手机号或者密码错误",
		}
		return
	}
	token, err := jwt.GenerateToken(userAuthorization.UserID)
	if err != nil {
		response = &user.Response{
			Code:    500,
			Message: "生成token出错",
		}
		return
	}
	if !gredis.SetToken(token, userAuthorization.UserID) {
		response = &user.Response{
			Code:    500,
			Message: "保存到redis出错",
		}
		return
	}
	response = &user.Response{
		Code:    200,
		Message: "success",
		Data: &user.TokenData{
			Token: token,
		},
	}
	return
}

func (server *UserServer) PhoneVerificationCodeLogin(context context.Context, request *user.PhoneCodeLoginRequest) (response *user.Response, err error) {
	if !gredis.CheckPhoneCode(request.Phone, request.Code) {
		response = &user.Response{
			Code:    20003,
			Message: "验证码有误",
		}
		return
	}
	var userAuthorization model.UserAuthorization
	result := model.DB.Model(&model.UserAuthorization{}).
		Where("identity_type = ?", PHONE).
		Where("identifier = ?", request.Phone).
		First(&userAuthorization)
	if result.Error == gorm.ErrRecordNotFound {
		uid := uuid.GenerateID()
		if model.DB.Create(&model.UserAuthorization{
			ID:           uuid.GenerateID(),
			UserID:       uid,
			IdentityType: "phone",
			Identifier:   request.Phone,
			Credential:   request.Phone,
		}).RowsAffected == 0 {
			response = &user.Response{
				Code:    500,
				Message: "新建用户出错",
			}
			return
		}
		if model.DB.Create(&model.User{
			ID:       uid,
			NickName: "",
			Gender:   0,
		}).RowsAffected == 0 {
			response = &user.Response{
				Code:    500,
				Message: "新建用户出错",
			}
			return
		}
		token, err1 := jwt.GenerateToken(uid)
		if err1 != nil {
			response = &user.Response{
				Code:    500,
				Message: "生成token出错",
			}
			return
		}
		if !gredis.SetToken(token, uid) {
			response = &user.Response{
				Code:    500,
				Message: "保存到redis出错",
			}
			return
		}
		response = &user.Response{
			Code:    200,
			Message: "新用户,请完善信息",
		}
		return
	}
	token, err := jwt.GenerateToken(userAuthorization.UserID)
	if err != nil {
		response = &user.Response{
			Code:    500,
			Message: "生成token出错",
		}
		return
	}
	if !gredis.SetToken(token, userAuthorization.UserID) {
		response = &user.Response{
			Code:    500,
			Message: "保存到redis出错",
		}
		return
	}
	response = &user.Response{
		Code:    200,
		Message: "success",
		Data: &user.TokenData{
			Token: token,
		},
	}
	return
}
