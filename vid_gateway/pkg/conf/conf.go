package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"gofiber-scaffold/client"
	"gofiber-scaffold/pkg/logger"
	"path/filepath"
)

// InitConfig initialize config
func InitConfig(configPath string) error {

	ext := filepath.Ext(configPath)                   // .yml
	filename := configPath[:len(configPath)-len(ext)] // xxx
	viper.SetConfigName(filename)
	viper.SetConfigType(ext[1:])
	viper.AddConfigPath(".")
	viper.SetDefault("meta.port", 3000)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("no such config file: %s", err))
		} else {
			panic(fmt.Errorf("read config error: %s", err))
		}
	}

	// initialize logger
	logger.NewLogger(logger.SetAppName(viper.GetString("meta.name")),
		logger.SetDebugFileName("debug"),
		logger.SetErrorFileName("error"),
		logger.SetWarnFileName("warn"),
		logger.SetInfoFileName("info"),
		logger.SetLogFileDir(viper.GetString("log.log-path")),
		logger.SetMaxAge(viper.GetInt("log.max-age")),
		logger.SetMaxBackups(viper.GetInt("log.max-back-up")),
		logger.SetMaxSize(viper.GetInt("log.max-size")),
		logger.SetLevel(zapcore.DebugLevel),
		logger.SetDevelopment(viper.GetBool("meta.debug")))

	client.Init()
	return nil
}
