package sql

import (
	"application/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql(cfg *conf.MysqlConfig) {
	var err error

	DB, err = gorm.Open(mysql.Open(cfg.Conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
