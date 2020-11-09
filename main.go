package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"vid-msa/middleware"
	"vid-msa/pkg/conf"
	"vid-msa/pkg/logger"
	"vid-msa/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/spf13/viper"
)

var (
	fConfig = flag.String("config", "./config.toml", "配置文件路径")
	fHelp   = flag.Bool("h", false, "show help")
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
	// views engine
	engine := html.New("./views", ".html")
	// engine debug
	engine.Debug(viper.GetBool("meta.debug"))
	engine.Reload(true)
	app := fiber.New(fiber.Config{
		Prefork:           viper.GetBool("meta.prefork"),
		Views:             engine,
		ReduceMemoryUsage: viper.GetBool("meta.reduce_memory"),
		// https://github.com/gofiber/fiber/issues/426
		Immutable: false,
		// close LOGO
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Info(fmt.Sprintf("[%s] shutting down", viper.GetString("meta.name")))
		_ = app.Shutdown()
	}()
	// initialize middleware
	middleware.InitMiddleWares(app)
	logger.Info(fmt.Sprintf("[%s] initalize router", viper.GetString("meta.name")))
	// initialize router
	router.InitRouter(app)
	logger.Info(fmt.Sprintf("[%s] start http server...", viper.GetString("meta.name")))
	// http
	if err := app.Listen(viper.GetString("meta.server_address")); err != nil {
		panic("Port is already in use !!!!")
	}
}
