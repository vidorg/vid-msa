package model

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB database instance
var DB *gorm.DB

// Database initialize database
func Database(dialector gorm.Dialector) error {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // 彩色打印
		},
	)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
		// disable foreign key constraint
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: &schema.NamingStrategy{
			TablePrefix: "tb_",
			// table name singular
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	DB = db
	s, err := db.DB()
	if err != nil {
		return err
	}
	// pool
	// idle
	s.SetMaxIdleConns(viper.GetInt("mysql.max-idle"))
	// open
	s.SetMaxOpenConns(viper.GetInt("mysql.max-active"))
	// timeout
	s.SetConnMaxLifetime(time.Duration(viper.GetInt("mysql.max-lifetime")) * time.Second)
	// auto migration
	migration()
	return nil
}
