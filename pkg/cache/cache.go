package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/spf13/viper"
)

var cache *bigcache.BigCache

func InitCache() {
	// 缓存设置为2周
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(time.Duration(viper.GetInt("jwt.expire")) * time.Second))
}

func Set(key string, value string) error {
	return cache.Set(key, []byte(value))
}

func Delete(key string) error {
	return cache.Delete(key)
}

func Get(key string) (string, error) {
	get, err := cache.Get(key)
	return string(get), err
}

func Exists(key string) bool {
	_, err := cache.Get(key)
	if err == bigcache.ErrEntryNotFound {
		return false
	}
	return true
}
