package v1

import (
	"application/log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetChangeLog(c *gin.Context) {
	dataSetIDStr := c.Query("data_set_id")
	dataSetID64, _ := strconv.ParseInt(dataSetIDStr, 10, 32)
	dataSetID32 := int32(dataSetID64)
	logs, _ := log.GetLog(dataSetID32)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": logs,
	})
}
