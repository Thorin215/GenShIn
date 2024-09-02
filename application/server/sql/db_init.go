package sql

import (
	"application/setting"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql(cfg *setting.MysqlConfig) {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	cfg.Password,
	// 	cfg.Host,
	// 	cfg.Port,
	// 	cfg.Database,
	// )
	dsn := "root:123456@tcp(127.0.0.1:3307)/genshin?charset=utf8mb4&parseTime=True&loc=Local"
	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
