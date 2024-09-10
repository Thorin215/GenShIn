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
		apiV1.POST("/queryUser", v1.QueryUser)
		apiV1.POST("/queryAllUsers", v1.QueryAllUsers)

		// dataset
		apiV1.POST("/createDataset", v1.CreateDataset)
		apiV1.POST("/queryAllDatasets", v1.QueryAllDatasets)
		apiV1.POST("/queryDatasetMetadata", v1.QueryDatasetMetadata)
		apiV1.POST("/addDatasetVersion", v1.AddDatasetVersion)

		// file
		apiV1.POST("/uploadFile", v1.UploadFile)
		apiV1.POST("/downloadFile", v1.DownloadFile)
		apiV1.POST("/downloadFilesCompressed", v1.DownloadFilesCompressed)

		// record
		apiV1.POST("/queryRecordsByUser", v1.QueryRecordsByUser)
		apiV1.POST("/queryRecordsByDataset", v1.QueryRecordsByDataset)

	}
	return r
}
