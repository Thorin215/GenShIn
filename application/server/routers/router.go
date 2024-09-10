package routers

import (
	v1 "application/api/v1"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由信息
func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiV1 := r.Group("/api/v1")
	{
		// hello
		apiV1.GET("/hello", v1.Hello)

		// user
		apiV1.POST("/user", v1.QueryUser)
		apiV1.POST("/user/all", v1.QueryAllUsers)
		apiV1.POST("/user/create", v1.CreateUser)
		apiV1.POST("/user/login", v1.CheckUserLogin)
		// dataset
		apiV1.POST("/dataset/create", v1.CreateDataset)
		apiV1.POST("/dataset/all", v1.QueryAllDatasets)
		apiV1.POST("/dataset/metadata", v1.QueryDatasetMetadata)
		apiV1.POST("/dataset/version/create", v1.AddDatasetVersion)

		// file
		apiV1.POST("/file/upload", v1.UploadFile)
		apiV1.POST("/file/download", v1.DownloadFile)
		apiV1.POST("/file/download/zip", v1.DownloadFilesCompressed)

		// record
		apiV1.POST("/record/by/user", v1.QueryRecordsByUser)
		apiV1.POST("/record/by/dataset", v1.QueryRecordsByDataset)

	}
	return r
}
