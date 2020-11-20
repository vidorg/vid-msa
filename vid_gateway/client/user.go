package client

import (
	pb "gofiber-scaffold/pb/user"
	"google.golang.org/grpc"
)

var UserClient pb.UserClient

func Init() {
	conn, err := grpc.Dial("localhost:20102", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	UserClient = pb.NewUserClient(conn)
}
