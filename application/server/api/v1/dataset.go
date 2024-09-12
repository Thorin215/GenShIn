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
		Owner    string         `json:"owner" binding:"required"`
		Name     string         `json:"name" binding:"required"`
		Metadata model.Metadata `json:"metadata" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	args := [][]byte{
		[]byte(body.Owner),
		[]byte(body.Name),
	}

	_, err := bc.ChannelExecute("createDataset", args)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	metadataBody := &sql.MetadataBody{
		Owner:     body.Owner,
		Name:      body.Name,
		Metadata:  body.Metadata,
		Downloads: 0,
		Deleted:   false,
	}

	if err := sql.CreateMetadata(metadataBody); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入数据库失败: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", "")
}

func QueryAllDatasets(c *gin.Context) {
	appG := app.Gin{C: c}

	res, err := bc.ChannelQuery("queryAllDatasets", nil)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
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

func QueryDatasetMetadata(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		Owner string `json:"owner" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("解析请求体失败: %s", err.Error()))
		return
	}

	owner := body.Owner
	name := body.Name

	if owner == "" || name == "" {
		appG.Response(http.StatusBadRequest, "失败", "所有者和数据集名称不能为空")
		return
	}

	metadataBody, err := sql.GetMetadata(owner, name)
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
		Owner   string        `json:"owner" binding:"required"`
		Name    string        `json:"name" binding:"required"`
		Version model.Version `json:"version" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	// 确认文件不为空
	if len(body.Version.Files) == 0 {
		appG.Response(http.StatusBadRequest, "失败", "文件不能为空")
		return
	}

	// 调用链码
	_, err := bc.ChannelExecute("addDatasetVersion", [][]byte{
		[]byte(body.Owner),
		[]byte(body.Name),
		[]byte(utils.ToJson(body.Version)),
	})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", "")
}

func QueryAllVersions(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		Owner string `json:"owner" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

	// 绑定并验证请求的 JSON 参数
	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	owner := body.Owner
	name := body.Name

	if owner == "" || name == "" {
		appG.Response(http.StatusBadRequest, "失败", "所有者和数据集名称不能为空")
		return
	}

	// 调用智能合约从区块链中查询数据集
	res, err := bc.ChannelQuery("queryDataset", [][]byte{
		[]byte(owner),
		[]byte(name),
	})

	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 解析返回的 JSON 数据
	var dataset model.Dataset
	if err = json.Unmarshal(res.Payload, &dataset); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	// 检查数据集是否已删除
	if dataset.Deleted {
		appG.Response(http.StatusBadRequest, "失败", "该数据集已被删除")
		return
	}

	// 返回版本列表
	appG.Response(http.StatusOK, "成功", dataset.Versions)
}
