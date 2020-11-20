package model

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

//Init ...
func Init() {
	dbParams := conf.GetString("mysql.dsn")
	err := Database(mysql.Open(dbParams))
	if err != nil {
		xlog.Panic("init mysql err")
	}
}

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
	// auto migration
	migration()
	return nil
}

// migration auto migration
func migration() {
	if err := DB.AutoMigrate(&User{}, &UserAuthorization{}); err != nil {
		panic("auto migration err")
	}
}
