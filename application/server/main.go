package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"application/blockchain"
	"application/pkg/cron"
	"application/routers"
	"application/setting"
	"application/sql"

	// "github.com/jinzhu/gorm"
	opLog "application/log"
	"application/model"
)

const DBfile = "./conf/config.ini"

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("时区设置失败 %s", err)
	}
	time.Local = timeLocal

	if err := setting.Init(DBfile); err != nil {
		log.Printf("配置数据库文件初始化失败 %s", err)
		return
	}

	sql.InitMysql(setting.Conf.MysqlConfig)

	sql.DB.AutoMigrate(&model.DataSet{})
	sql.DB.AutoMigrate(&model.MetaData{})
	sql.DB.AutoMigrate(&opLog.Log{})

	model.CreateDataSet(&model.DataSet{
		Name:         "default",
		AccountID:    1,
		Stars:        0,
		DataSetID:    1,
		ModifiedTime: time.Now(),
		CreateTime:   time.Now(),
		Data:         []byte("default"),
	}, &model.MetaData{
		DataSetID:  1,
		Tasks:      "default",
		Modalities: "default",
		Formats:    "default",
		SubTasks:   "default",
		Languages:  "default",
		Libararies: "default",
		Tags:       "default",
		License:    "default",
		Rows:       0,
	})

	opLog.CreateLog(&opLog.Log{
		LogID:     1,
		DataSetID: 1,
		OpMsg:     "default",
		Status:    "default",
		TimeStamp: time.Now(),
	})

	blockchain.Init()
	go cron.Init()

	endPoint := fmt.Sprintf("0.0.0.0:%d", 8888)
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
