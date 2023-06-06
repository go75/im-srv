package initialize

import (
	"im-srv/ws-srv/global"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DB() {
	// sql日志
	logger := logger.New (
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config {
			SlowThreshold: time.Duration(global.Config.MySQL.SlowThreshold) * time.Millisecond, // 慢SQL阈值
			//LogLevel: logger.Info, // 级别
			LogLevel: logger.LogLevel(global.Config.MySQL.LogLevel),
			Colorful: global.Config.MySQL.Colorful, // 彩色
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(global.Config.MySQL.DSN), &gorm.Config {
		Logger: logger,
	})
	if err!=nil {
		panic(err)
	}
}