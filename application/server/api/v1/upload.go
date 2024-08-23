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

type SentenceRequestBody struct {
	Sentence  string `json:"sentence" binding:"required"` // 句子
	Label     bool   `json:"label" binding:"required"`  // 标签
	DatasetID string `json:"dataset_id" binding:"required"` // 所属数据集编号
}

func UploadSentence(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SentenceRequestBody)

	// 解析请求体
	if err := c.ShouldBindJSON(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	// 文件路径
	dataFilePath := filepath.Join("data", "setRecord.json")

	// 读取数据集记录
	records, err := readRecords(dataFilePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取记录时发生错误: %s", err.Error()))
		return
	}

	var datasetRecord *Record
	for i := range records {
		if records[i].DatasetID == body.DatasetID {
			datasetRecord = &records[i]
			break
		}
	}

	if datasetRecord == nil {
		appG.Response(http.StatusNotFound, "失败", "数据集ID不存在")
		return
	}

	// 生成句子编号
	sentenceID := fmt.Sprintf("%d", datasetRecord.Count)

	// 更新数据集记录的计数
	datasetRecord.Count++
	datasetRecord.LastModified = time.Now()
	// 将更新后的记录写入文件
	if err := writeRecords(dataFilePath, records); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入记录时发生错误: %s", err.Error()))
		return
	}

	// 文件名格式为 DataSet{DatasetID}.json
	fileName := fmt.Sprintf("%s_%s.json", "DataSet", body.DatasetID)
	filePath := filepath.Join("data", fileName) // 存储在 data 目录下

	// 创建 data 目录（如果不存在）
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建目录失败: %s", err.Error()))
		return
	}

	// 读取现有数据
	var sentenceDataList []struct {
		Sentence  string `json:"sentence"`
		Label     bool   `json:"label"`
		DatasetID string `json:"dataset_id"`
		SentenceID string `json:"sentence_id"`
	}

	if _, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("打开文件失败: %s", err.Error()))
			return
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(&sentenceDataList); err != nil && err != io.EOF {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取文件数据失败: %s", err.Error()))
			return
		}
	}

	// 添加新数据
	sentenceDataList = append(sentenceDataList, struct {
		Sentence  string `json:"sentence"`
		Label     bool   `json:"label"`
		DatasetID string `json:"dataset_id"`
		SentenceID string `json:"sentence_id"`
	}{
		Sentence:  body.Sentence,
		Label:     body.Label,
		DatasetID: body.DatasetID,
		SentenceID: sentenceID,
	})

	// 打开文件进行写入（覆盖原文件）
	file, err := os.Create(filePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("无法创建文件: %s", err.Error()))
		return
	}
	defer file.Close()

	// 将数据写入文件
	if err := json.NewEncoder(file).Encode(sentenceDataList); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入文件失败: %s", err.Error()))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", gin.H{
		"sentence_id": sentenceID,
	})
}


// SetRequestBody 用于接收账户ID和数据集ID的请求体
type SetRequestBody struct {
	AccountID    string    `json:"account_id"`
	DatasetID    string    `json:"dataset_id"`
	CreationTime time.Time `json:"creation_time"`
}

type Record struct {
	AccountID    string    `json:"account_id"`
	DatasetID    string    `json:"dataset_id"`
	CreationTime time.Time `json:"creation_time"`
	LastModified time.Time `json:"last_modified"`
	Count        int       `json:"count"` // 默认为 0
}

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

	// 读取现有记录
	records, err := readRecords(filePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取记录时发生错误: %s", err.Error()))
		return
	}

	// 检查是否存在相同 dataset_id 的记录
	for _, record := range records {
		if record.DatasetID == body.DatasetID {
			appG.Response(http.StatusConflict, "失败", "数据集ID已存在")
			return
		}
	}

	// 添加新的记录
	newRecord := Record{
		AccountID:    body.AccountID,
		DatasetID:    body.DatasetID,
		CreationTime: body.CreationTime, // 从请求体中获取创建时间
		LastModified: time.Now(),        // 设置最后修改时间为当前时间
		Count:        0,                 // 默认数量为 0
	}
	records = append(records, newRecord)

	// 写入记录到文件
	if err := writeRecords(filePath, records); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入记录时发生错误: %s", err.Error()))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", gin.H{
		"dataset_id":    newRecord.DatasetID,
		"account_id":    newRecord.AccountID,
		"creation_time": newRecord.CreationTime.Format(time.RFC3339), // 返回创建时间
		"last_modified": newRecord.LastModified.Format(time.RFC3339), // 返回最后修改时间
	})
}

// 从文件中读取现有记录
func readRecords(filePath string) ([]Record, error) {
	var records []Record
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return records, nil // 文件不存在，则返回空记录
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return nil, err
	}

	return records, nil
}

// 将记录写入文件
func writeRecords(filePath string, records []Record) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 设置格式化输出
	if err := encoder.Encode(records); err != nil {
		return err
	}

	return nil
}