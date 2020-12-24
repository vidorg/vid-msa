package config

import "github.com/tal-tech/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Secret         string
	DataSourceName string
	RedisAddr      string
	RedisPass      string
}
