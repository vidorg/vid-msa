package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
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
	// limiter
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.JSON(&fiber.Map{
				"code":    500,
				"message": "盗刷流量？？？！！！",
			})
		},
	}))
	// Recover
	//app.Use(recover.New())
}
