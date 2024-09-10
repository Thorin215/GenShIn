package v1

import (
	bc "application/blockchain"
	"application/model"
	"application/pkg/app"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

func QueryRecordsByUser(c *gin.Context) {
	appG := app.Gin{C: c}

	var body struct {
		User string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err))
		return
	}

	res, err := bc.ChannelQuery("queryRecordsByUser", [][]byte{[]byte(body.User)})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err))
		return
	}

	var records []model.Record
	if err = json.Unmarshal(res.Payload, &records); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", records)
}

func QueryRecordsByDataset(c *gin.Context) {
	appG := app.Gin{C: c}

	var body struct {
		Owner string `json:"owner" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err))
		return
	}

	res, err := bc.ChannelQuery("queryRecordsByDataset", [][]byte{[]byte(body.Owner), []byte(body.Name)})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err))
		return
	}

	var records []model.Record
	if err = json.Unmarshal(res.Payload, &records); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", records)
}
