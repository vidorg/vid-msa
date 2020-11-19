package user

import (
	"context"

	"vid_user/pb/user"
)

type UserServiceServer struct{}

func (server *UserServiceServer) SayHello(context context.Context, request *user.UserRequest) (response *user.UserResponse, err error) {
	response = &user.UserResponse{
		Reply: "xxxx",
	}
	return response, nil
}
