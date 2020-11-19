package engine

import (
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/client/redis"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
	"vid_user/internal/app/gredis"
	grpc "vid_user/internal/app/grpc/user"
	"vid_user/pb/user"
)

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		xgo.ParallelWithError(
			eng.InitRedisClient,
			eng.serveGRPC,
			eng.serveHTTP,
		),
	); err != nil {
		xlog.Panic("startup user service err", xlog.Any("err", err))
	}
	return eng
}

func (eng *Engine) serveHTTP() error {
	server := xecho.StdConfig("http").Build()

	//support proxy for http to grpc controller
	g := grpc.UserServiceServer{}
	group2 := server.Group("/grpc")
	group2.GET("/get", xecho.GRPCProxyWrapper(g.SayHello))
	return eng.Serve(server)
}

func (eng *Engine) serveGRPC() error {
	server := xgrpc.StdConfig("grpc").Build()

	user.RegisterUserServiceServer(
		server.Server,
		&grpc.UserServiceServer{},
	)
	return eng.Serve(server)
}

func (eng *Engine) InitRedisClient() error {
	gredis.RedisClient = redis.StdRedisStubConfig("redis").Build()
	return nil
}
