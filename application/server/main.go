package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"application/blockchain"
	"application/conf"
	"application/pkg/cron"
	"application/routers"
	"application/sql"

	"gorm.io/gorm"
)

type TestData struct {
	gorm.Model
}

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("时区设置失败 %s", err)
	}
	time.Local = timeLocal

	if err := conf.Init(); err != nil {
		log.Printf("配置数据库文件初始化失败 %s", err)
		return
	}

	sql.InitMysql(conf.Conf.MysqlConfig)
	err = sql.Migrate()
	if err != nil {
		log.Printf("数据库迁移失败 %s", err)
		return
	}

	blockchain.Init()
	go cron.Init()

	endPoint := fmt.Sprintf("%s:%s", conf.Conf.ServerConfig.Host, conf.Conf.ServerConfig.Port)
	server := &http.Server{
		Addr:    endPoint,
		Handler: routers.InitRouter(),
	}

	log.Printf("[info] start http server listening %s", endPoint)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("start http server failed %s", err)
	}
}
