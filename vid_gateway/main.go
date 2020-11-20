package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gofiber-scaffold/middleware"
	"gofiber-scaffold/pkg/conf"
	"gofiber-scaffold/pkg/logger"
	"gofiber-scaffold/router"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	fConfig = flag.String("config", "./config.toml", "配置文件路径")
	fHelp   = flag.Bool("h", false, "show help")

	survivalTimeout = int(3e9)
)

func main() {
	flag.Parse()
	if *fHelp {
		flag.Usage()
	} else {
		run()
	}
}

func run() {
	// initialize config
	if err := conf.InitConfig(*fConfig); err != nil {
		panic(err)
	}
	app := fiber.New(fiber.Config{
		Prefork:               viper.GetBool("meta.prefork"),
		ReduceMemoryUsage:     viper.GetBool("meta.reduce_memory"),
		Immutable:             false,
		DisableStartupMessage: false,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  500,
				"msg":   "server error !!",
				"error": err.Error(),
			})
		},
	})
	// graceful shutdown
	go func() {
		signals := make(chan os.Signal, 1)
		// It is not possible to block SIGKILL or syscall.SIGSTOP
		signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			sig := <-signals
			logger.Info("get signal" + sig.String())
			switch sig {
			case syscall.SIGHUP:
				// reload()
			default:
				time.AfterFunc(time.Duration(survivalTimeout), func() {
					logger.Info(fmt.Sprintf("[%s] shutting down", viper.GetString("meta.name")))
					_ = app.Shutdown()
				})

				return
			}
		}
	}()
	// initialize middleware
	middleware.InitMiddleWares(app)
	// initialize router
	router.InitRouter(app)
	if err := app.Listen(viper.GetString("meta.server_address")); err != nil {
		panic("Port is already in use !!!!")
	}
}
