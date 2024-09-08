package routers

import (
	v1 "application/api/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由信息
func InitRouter() *gin.Engine {
	r := gin.Default()

	// CORS 处理
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Type")

		// Preflight request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	apiV1 := r.Group("/api/v1")
	{
		// hello
		apiV1.GET("/hello", v1.Hello)

		// user
		apiV1.POST("/queryUser", v1.QueryUser)
		apiV1.POST("/queryAllUsers", v1.QueryAllUsers)

		// dataset
		apiV1.POST("/createDataset", v1.CreateDataset)
		apiV1.POST("/getAllDatasets", v1.GetAllDatasets)
		apiV1.POST("/getDatasetMetadata", v1.GetDatasetMetadata)

		// file
		apiV1.POST("/uploadFile", v1.UploadFile)
		apiV1.POST("/downloadFile", v1.DownloadFile)
		apiV1.POST("/downloadFilesCompressed", v1.DownloadFilesCompressed)

	}
	return r
}
