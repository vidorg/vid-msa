package middleware

import (
	"time"
	"vid-msa/serializer"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis"
	"github.com/spf13/viper"
)

// InitMiddleWares initialize middleware
func InitMiddleWares(app *fiber.App) {
	// Favicon
	app.Use(favicon.New())
	// logger
	app.Use(Logger())
	// CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	// compress
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	// fiber storage
	storage := redis.New(redis.Config{
		Host:     viper.GetString("cache.redis_host"),
		Port:     viper.GetInt("cache.redis_port"),
		Password: viper.GetString("cache.redis_password"),
		Database: viper.GetInt("cache.redis_db"),
		Reset:    false,
	})
	// limiter
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
		Key: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.JSON(serializer.Response{
				Code: 500,
				Msg:  "warning！！！",
			})
		},
		Store: storage,
	}))
	// cache
	if viper.GetBool("cache.enable") {
		app.Use(cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.Query("refresh") == "true"
			},
			Expiration:   time.Duration(viper.GetInt("cache.expiration")) * time.Minute,
			CacheControl: true,
			Store:        storage,
		}))
		// cache etag
		app.Use(etag.New())
	}
	// Recover 用panic正常响应
	app.Use(recover.New())
	// monitor
	app.Get("/", monitor.New())
	// custom middleware
	app.Use(JWT())
	app.Use(GetUser())
}
