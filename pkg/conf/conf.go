package conf

import (
	"fmt"
	"path/filepath"
	"vid-msa/model"
	"vid-msa/pkg/gredis"
	"vid-msa/pkg/logger"
	_ "vid-msa/pkg/session"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
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

	// initialize database
	dbParams := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		viper.GetString("mysql.user"), viper.GetString("mysql.password"),
		viper.GetString("mysql.host"), viper.GetInt("mysql.port"),
		viper.GetString("mysql.db"), viper.GetString("mysql.charset"),
	)

	err := model.Database(mysql.Open(dbParams))
	// Failed to connect to MySQL, switch to sqlite
	if err != nil {
		logger.Error(fmt.Sprintf("[%s] Failed to connect to MySQL, switch to sqlite", viper.GetString("meta.name")))
		err := model.Database(sqlite.Open(viper.GetString("meta.name") + ".db"))
		if err != nil {
			logger.Error(fmt.Sprintf("[%s] Failed to enable the sqlite", viper.GetString("meta.name")))
			return err
		}
	}
	logger.Info(fmt.Sprintf("[%s] connect to Redis", viper.GetString("meta.name")))
	// initialize Redis
	if err := gredis.InitRedisClient(); err != nil {
		logger.Panic(fmt.Sprintf("[%s] \"Failed to start redis!\", err", viper.GetString("meta.name")))
	}
	// session store to redis
	//session.Session = session.InitSession()
	// initialize Schtasks
	//if err := cron.StartTasks();err!=nil{
	//	return err
	//}
	return nil
}
