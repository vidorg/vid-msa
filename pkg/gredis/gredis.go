package gredis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	Rdb *redis.Client
	ctx = context.Background()
)

// InitRedisClient initialize redis client
func InitRedisClient() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	pong, err := Rdb.Ping(ctx).Result()
	if err != nil || pong == "PONG" {
		return err
	}
	return nil
}

// Get get value from redis
func Get(key string) (string, error) {
	return Rdb.Get(ctx, key).Result()
}

// Set set value to redis
func Set(key string, value interface{}, expiration time.Duration) error {
	_, err := Rdb.Set(ctx, key, value, expiration).Result()
	return err
}

// Exists exist in redis
func Exists(key string) bool {
	result, err := Rdb.Exists(ctx, key).Result()
	if err != nil || result != 1 {
		return true
	}
	return false
}

// Delete delete from redis
func Delete(key string) error {
	_, err := Rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// Do
func Do(args ...interface{}) (interface{}, error) {
	result, err := Rdb.Do(ctx, args).Result()
	return result, err
}

// Publish publish to channel
func Publish(channel string, message interface{}) error {
	return Rdb.Publish(ctx, channel, message).Err()
}

// Subscribe subscribe a channel
func Subscribe(channel string) *redis.PubSub {
	return Rdb.Subscribe(ctx, channel)
}

// SubscribeChan subscribe a channel,return <- chan
func SubscribeChan(channel string) <-chan *redis.Message {
	return Rdb.Subscribe(ctx, channel).Channel()
}
