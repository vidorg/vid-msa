package gredis

import (
	"github.com/douyu/jupiter/pkg/client/redis"
	"time"
	jwt2 "vid_user/internal/app/jwt"
)

var RedisClient *redis.Redis

const ValidCodePrefix = "vid:phone_code:"

func CheckPhoneCode(phone string, code string) bool {
	validCode := RedisClient.Get(ValidCodePrefix + phone)
	if validCode == "" {
		return false
	}
	return validCode == code
}

func SetToken(token string, uid int64) bool {
	return RedisClient.Set(jwt2.RedisTokenConcat(uid, token), "1", 1000*time.Hour)
}

func CheckToken(token string, uid int64) bool {
	return !(RedisClient.Get(jwt2.RedisTokenConcat(uid, token)) == "")
}
