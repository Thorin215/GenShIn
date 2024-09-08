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
	sql.Migrate()

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

// import (
//     "github.com/gin-gonic/gin"
//     "github.com/gin-contrib/cors"
// )

// func main() {
//     router := gin.Default()

//     // CORS 配置
//     router.Use(cors.New(cors.Config{
//         AllowOrigins:     []string{"http://localhost:8000"},
//         AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
//         AllowHeaders:     []string{"Origin", "Content-Type"},
//         AllowCredentials: true,
//     }))

//     // 其他路由定义...
//     router.Run(":8888")
// }
