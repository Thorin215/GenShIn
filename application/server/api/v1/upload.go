package v1

import (
	//bc "application/blockchain"
	"application/pkg/app"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"io"
	"github.com/gin-gonic/gin"
)

// type SentenceRequestBody struct {
// 	Sentence  string `json:"sentence" binding:"required"` // 句子
// 	Label     bool   `json:"label" binding:"required"`  // 标签
// 	DatasetID string `json:"dataset_id" binding:"required"` // 所属数据集编号
// }

// func UploadSentence(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	body := new(SentenceRequestBody)

// 	// 解析请求体
// 	if err := c.ShouldBindJSON(body); err != nil {
// 		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
// 		return
// 	}

// 	// 文件路径
// 	dataFilePath := filepath.Join("data", "setRecord.json")

// 	// 读取数据集记录
// 	records, err := readRecords(dataFilePath)
// 	if err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取记录时发生错误: %s", err.Error()))
// 		return
// 	}

// 	var datasetRecord *Record
// 	for i := range records {
// 		if records[i].DatasetID == body.DatasetID {
// 			datasetRecord = &records[i]
// 			break
// 		}
// 	}

// 	if datasetRecord == nil {
// 		appG.Response(http.StatusNotFound, "失败", "数据集ID不存在")
// 		return
// 	}

// 	// 生成句子编号
// 	sentenceID := fmt.Sprintf("%d", datasetRecord.Count)

// 	// 更新数据集记录的计数
// 	datasetRecord.Count++
// 	datasetRecord.LastModified = time.Now()
// 	// 将更新后的记录写入文件
// 	if err := writeRecords(dataFilePath, records); err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入记录时发生错误: %s", err.Error()))
// 		return
// 	}

// 	// 文件名格式为 DataSet{DatasetID}.json
// 	fileName := fmt.Sprintf("%s_%s.json", "DataSet", body.DatasetID)
// 	filePath := filepath.Join("data", fileName) // 存储在 data 目录下

// 	// 创建 data 目录（如果不存在）
// 	if err := os.MkdirAll("data", os.ModePerm); err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建目录失败: %s", err.Error()))
// 		return
// 	}

// 	// 读取现有数据
// 	var sentenceDataList []struct {
// 		Sentence  string `json:"sentence"`
// 		Label     bool   `json:"label"`
// 		DatasetID string `json:"dataset_id"`
// 		SentenceID string `json:"sentence_id"`
// 	}

// 	if _, err := os.Stat(filePath); err == nil {
// 		file, err := os.Open(filePath)
// 		if err != nil {
// 			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("打开文件失败: %s", err.Error()))
// 			return
// 		}
// 		defer file.Close()

// 		if err := json.NewDecoder(file).Decode(&sentenceDataList); err != nil && err != io.EOF {
// 			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取文件数据失败: %s", err.Error()))
// 			return
// 		}
// 	}

// 	// 添加新数据
// 	sentenceDataList = append(sentenceDataList, struct {
// 		Sentence  string `json:"sentence"`
// 		Label     bool   `json:"label"`
// 		DatasetID string `json:"dataset_id"`
// 		SentenceID string `json:"sentence_id"`
// 	}{
// 		Sentence:  body.Sentence,
// 		Label:     body.Label,
// 		DatasetID: body.DatasetID,
// 		SentenceID: sentenceID,
// 	})

// 	// 打开文件进行写入（覆盖原文件）
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("无法创建文件: %s", err.Error()))
// 		return
// 	}
// 	defer file.Close()

// 	// 将数据写入文件
// 	if err := json.NewEncoder(file).Encode(sentenceDataList); err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入文件失败: %s", err.Error()))
// 		return
// 	}

// 	// 成功响应
// 	appG.Response(http.StatusOK, "成功", gin.H{
// 		"sentence_id": sentenceID,
// 	})
// }

// SetRequestBody 用于接收账户ID和数据集的请求体
type SetRequestBody struct {
	Owner        string          `json:"owner"`
	Name         string          `json:"name"`         // 数据集名称
	CreationTime time.Time       `json:"creation_time"`
	Metadata     DatasetMetadata `json:"metadata"`     // 元数据
}

// DatasetMetadata 数据集元数据
type DatasetMetadata struct {
	Tasks      []string `json:"tasks"`      // 目标任务
	Modalities []string `json:"modalities"` // 数据模态
	Formats    []string `json:"formats"`    // 文件格式
	SubTasks   []string `json:"sub_tasks"`  // 子任务
	Languages  []string `json:"languages"`  // 语言
	Libraries  []string `json:"libraries"`  // 适用库
	Tags       []string `json:"tags"`       // 标签
	License    string   `json:"license"`    // 许可证
}

// Dataset 数据集
type Dataset struct {
	Owner     string           `json:"owner"`     // 所有者ID
	Name      string           `json:"name"`      // 数据集名
	Metadata  DatasetMetadata `json:"metadata"`  // 元数据
	Versions  []DatasetVersion `json:"versions"` // 版本列表
}

// DatasetVersion 数据集一个版本
type DatasetVersion struct {
	CreationTime string   `json:"creation_time"` // 创建时间
	ChangeLog    string   `json:"change_log"`    // 版本说明
	Rows         int32    `json:"rows"`          // 行数
	Files        []string `json:"files"`         // 文件哈希列表
}

// UploadSet 处理数据集上传请求
func UploadSet(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SetRequestBody)

	// 解析请求体
	if err := c.ShouldBindJSON(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	// 文件路径
	filePath := filepath.Join("data", "setRecord.json")

	// 读取现有数据集
	datasets, err := readDatasets(filePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取数据集时发生错误: %s", err.Error()))
		return
	}

	// 检查是否存在相同数据集名称的记录
	for _, dataset := range datasets {
		if dataset.Name == body.Name {
			appG.Response(http.StatusConflict, "失败", "数据集名称已存在")
			return
		}
	}

	// 创建新的数据集记录
	newDataset := Dataset{
		Owner:     body.Owner,
		Name:      body.Name,
		Metadata:  body.Metadata,
		Versions:  []DatasetVersion{}, // 初始时没有版本
	}

	// 添加新的记录
	datasets = append(datasets, newDataset)

	// 写入数据集到文件
	if err := writeDatasets(filePath, datasets); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入数据集时发生错误: %s", err.Error()))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", gin.H{
		"dataset_name":  newDataset.Name,
		"owner":         newDataset.Owner,
		// "creation_time": body.CreationTime, // 返回创建时间
	})
}


func writeDatasets(filePath string, datasets []Dataset) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ") // 设置格式化输出
    if err := encoder.Encode(datasets); err != nil {
        return err
    }

    return nil
}

func initializeFile(filePath string) error {
    return writeDatasets(filePath, []Dataset{})
}

func readDatasets(filePath string) ([]Dataset, error) {
    var datasets []Dataset
    file, err := os.Open(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            // 文件不存在，返回空数据集
            return datasets, nil
        }
        return nil, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&datasets)
    if err != nil && err != io.EOF {
        return nil, err
    }

    return datasets, nil
}

