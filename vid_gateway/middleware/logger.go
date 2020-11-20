package middleware

import (
	"gofiber-scaffold/pkg/logger"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

// Logger logger middleware
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		LoggerWithZap(logger.Log.Logger, start, c)
		return c.Next()
	}
}

// LoggerWithZap zap logger
func LoggerWithZap(logger *zap.Logger, start time.Time, c *fiber.Ctx) {
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	ip := c.IP()
	code := c.Response().StatusCode()
	//resp := c.Response().String()
	length := len(c.Response().Body())
	// ignore monitor,swagger
	if path == "/" || strings.Contains(path, "/doc") {
		return
	}
	if code >= 500 {
		logger.Error("[Web]",
			zap.Int("code", code),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("length", length),
			zap.Any("latency", latency),
			zap.String("IP", ip))
	} else if code >= 400 {
		logger.Warn("[Fiber]",
			zap.Int("code", code),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("length", length),
			zap.Any("latency", latency),
			zap.String("IP", ip))
	} else {
		logger.Info("[Fiber]",
			zap.Int("code", code),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("length", length),
			zap.Any("latency", latency),
			zap.String("IP", ip))
	}
}
