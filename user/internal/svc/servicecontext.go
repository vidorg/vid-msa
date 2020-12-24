package svc

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"user/internal/config"
	"user/internal/model"
)

type ServiceContext struct {
	c      config.Config
	Secret []byte
	DB     *gorm.DB
	Rdb    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSourceName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_",
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(&model.User{}); err != nil {
		fmt.Println("db AutoMigrate err,", err.Error())
	}
	rdb := redis.NewRedis(c.RedisAddr, redis.NodeType)
	if !rdb.Ping() {
		panic(errors.New("redis ping err"))
	}
	return &ServiceContext{
		c:      c,
		Secret: []byte(c.Secret),
		DB:     db,
		Rdb:    rdb,
	}
}
