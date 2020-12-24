package main

import (
	"context"
	"fmt"
	"log"
	"user/user"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc"
)

func main() {
	cc := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "user.rpc",
		},
	})

	client := user.NewUserClient(cc.Conn())

	reply, err := client.Login(context.Background(), &user.LoginRequest{
		Username: "test",
		Password: "123",
	})
	if err != nil {
		log.Fatal(err)
	}
	reply2, err := client.Register(context.Background(), &user.RegisterRequest{
		Username: "test1231",
		Password: "123",
	})
	if err != nil {
		log.Fatal(err)
	}
	reply3, err := client.GetUser(context.Background(), &user.GetUserRequset{
		Token: reply.Token,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.Token)
	fmt.Println(reply2.Username)
	fmt.Println(reply2.Token)
	fmt.Println(reply3.Username)
}
