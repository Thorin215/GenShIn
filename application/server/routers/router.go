package routers

import (
	"net/http"
	v1 "application/api/v1"
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
		apiV1.GET("/hello", v1.Hello)
		apiV1.POST("/queryAccountList", v1.QueryAccountList)
		apiV1.POST("/checkAccount", v1.CheckAccount)
		apiV1.POST("/uploadSet", v1.UploadSet) // 增加上传句子接口
		apiV1.POST("/getalldataset",v1.GetAllDataSet)
		apiV1.GET("/getChangeLog", v1.GetChangeLog)
		//apiV1.POST("/updateVersion", v1.UpdateVersion)
	}
	return r
}