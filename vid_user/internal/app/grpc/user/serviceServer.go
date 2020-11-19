package user

import (
	"context"

	"vid_user/pb/user"
)

type ServiceServer struct{}

func (server *ServiceServer) PhoneLogin(context context.Context, request *user.PhoneLoginRequest) (response *user.Response, err error) {
	response = &user.Response{
		Token:        "token" + request.Phone,
		RefreshToken: "resfrsh token",
	}
	err = nil
	return
}

func (server *ServiceServer) EmailLogin(context context.Context, request *user.EmailLoginRequest) (response *user.Response, err error) {
	response = &user.Response{
		Token:        "token1",
		RefreshToken: "resfrsh1 token",
	}
	return
}

func (server *ServiceServer) PhoneVerificationCodeLogin(context context.Context, request *user.PhoneCodeLoginRequest) (response *user.Response, err error) {
	response = &user.Response{
		Token:        "token2",
		RefreshToken: "resfrsh2 token",
	}
	return
}
