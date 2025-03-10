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

	// Query the blockchain
	res, err := bc.ChannelQuery("queryAllDatasets", nil)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// Deserialize the JSON response
	var datasets []model.Dataset
	if err = json.Unmarshal(res.Payload, &datasets); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	// Filter out datasets where Deleted is true
	var activeDatasets []model.DatasetEx
	for _, dataset := range datasets {
		if !dataset.Deleted {
			// Query the database for the download count
			downloads, err := sql.QueryDownloads(dataset.Owner, dataset.Name)
			if err != nil {
				appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("数据库错误: %s", err.Error()))
				return
			}
			// Append the dataset to the active datasets
			activeDatasets = append(activeDatasets, model.DatasetEx{
				Owner:     dataset.Owner,
				Name:      dataset.Name,
				Versions:  dataset.Versions,
				Downloads: downloads,
				Deleted:   dataset.Deleted,
			})
		}
	}

	// Return the active datasets
	appG.Response(http.StatusOK, "成功", activeDatasets)
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

func DeleteDataset(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		Owner string `json:"owner" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

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

	// 调用链码删除数据集
	_, err := bc.ChannelExecute("deleteDataset", [][]byte{
		[]byte(owner),
		[]byte(name),
	})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 标注数据集已删除
	if err := sql.MarkDeleted(owner, name); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("数据库出错: %s", err.Error()))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", "success")
}
