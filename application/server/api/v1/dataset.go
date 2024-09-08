package v1

import (
	bc "application/blockchain"
	"application/model"
	"application/pkg/app"
	"application/pkg/utils"
	"application/sql"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDataset(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		Owner    string         `json:"owner"`
		Name     string         `json:"name"`
		Metadata model.Metadata `json:"metadata"`
	}

	if err := c.ShouldBindJSON(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	args := [][]byte{
		[]byte(body.Owner),
		[]byte(body.Name),
	}

	res, err := bc.ChannelExecute("createDataset", args)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	payload := string(res.Payload)
	if res.ChaincodeStatus != 200 {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("智能合约出错: %s", payload))
		return
	}

	metadataBody := &sql.MetadataBody{
		Owner:    body.Owner,
		Name:     body.Name,
		Metadata: body.Metadata,
	}

	if err := sql.CreateMetaData(metadataBody); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入数据库失败: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", "")
}

func GetAllDatasets(c *gin.Context) {
	appG := app.Gin{C: c}

	res, err := bc.ChannelQuery("queryAllDatasets", nil)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}
	if res.ChaincodeStatus != 200 {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("智能合约出错: %s", string(res.Payload)))
		return
	}

	// 反序列化 JSON
	var datasets []model.Dataset
	if err = json.Unmarshal(res.Payload, &datasets); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", datasets)
}

func GetDatasetMetadata(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		Owner string `json:"owner"`
		Name  string `json:"name"`
	}

	if err := c.BindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("解析请求体失败: %s", err.Error()))
		return
	}

	owner := body.Owner
	name := body.Name

	if owner == "" || name == "" {
		appG.Response(http.StatusBadRequest, "失败", "所有者和数据集名称不能为空")
		return
	}

	metadataBody, err := sql.GetMetaData(name, owner)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("查询数据库失败: %s", err.Error()))
		return
	}

	if metadataBody == nil {
		appG.Response(http.StatusNotFound, "失败", "数据集不存在")
		return
	}

	metadata := metadataBody.Metadata
	appG.Response(http.StatusOK, "成功", metadata)
}

func AddDatasetVersion(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		Owner   string        `json:"owner"`
		Name    string        `json:"name"`
		Version model.Version `json:"version"`
	}

	if err := c.ShouldBindJSON(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	// 调用链码
	res, err := bc.ChannelExecute("appendDatasetVersion", [][]byte{
		[]byte(body.Owner),
		[]byte(body.Name),
		[]byte(utils.ToJson(body.Version)),
	})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}
	if res.ChaincodeStatus != 200 {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("智能合约出错: %s", string(res.Payload)))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", "")
}
