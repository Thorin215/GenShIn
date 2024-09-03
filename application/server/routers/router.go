// package routers
//
// import (
// 	v1 "application/api/v1"
// 	"github.com/gin-gonic/gin"
// )
//
// // InitRouter 初始化路由信息
// func InitRouter() *gin.Engine {
// 	r := gin.Default()
//
// 	apiV1 := r.Group("/api/v1")
// 	{
// 		apiV1.GET("/hello", v1.Hello)
// 		apiV1.POST("/queryAccountList", v1.QueryAccountList)
// 		apiV1.POST("/createRealEstate", v1.CreateRealEstate)
// 		apiV1.POST("/queryRealEstateList", v1.QueryRealEstateList)
// 		apiV1.POST("/createSelling", v1.CreateSelling)
// 		apiV1.POST("/createSellingByBuy", v1.CreateSellingByBuy)
// 		apiV1.POST("/querySellingList", v1.QuerySellingList)
// 		apiV1.POST("/querySellingListByBuyer", v1.QuerySellingListByBuyer)
// 		apiV1.POST("/updateSelling", v1.UpdateSelling)
// 		apiV1.POST("/createDonating", v1.CreateDonating)
// 		apiV1.POST("/queryDonatingList", v1.QueryDonatingList)
// 		apiV1.POST("/queryDonatingListByGrantee", v1.QueryDonatingListByGrantee)
// 		apiV1.POST("/updateDonating", v1.UpdateDonating)
// 	}
// 	return r
// }
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
		apiV1.POST("/uploadSet", v1.UploadSet)
		apiV1.POST("/updateVersion", v1.UpdateVersion)
	}
	return r
}
