// package v1

// import (
// 	bc "application/blockchain"
// 	"application/pkg/app"
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // SentenceRequestBody 用于接收句子、标签、所属数据集编号及句子自身编号的请求体
// type SentenceRequestBody struct {
// 	Sentence    string `json:"sentence" binding:"required"`    // 句子
// 	Label       bool   `json:"label" binding:"required"`       // 标签
// 	DatasetID   string `json:"dataset_id" binding:"required"`  // 所属数据集编号
// 	SentenceID  string `json:"sentence_id" binding:"required"` // 句子自身编号
// }

// func UploadSentence(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	body := new(SentenceRequestBody)

// 	// 解析请求体
// 	if err := c.ShouldBindJSON(body); err != nil {
// 		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
// 		return
// 	}

// 	// 在这里你可以处理接收到的句子、标签、数据集编号及句子编号
// 	// 例如：
// 	fmt.Printf("Received sentence: %s, label: %t, dataset ID: %s, sentence ID: %s\n", body.Sentence, body.Label, body.DatasetID, body.SentenceID)

// 	// 成功响应
// 	appG.Response(http.StatusOK, "成功", gin.H{
// 		"sentence":   body.Sentence,
// 		"label":      body.Label,
// 		"dataset_id": body.DatasetID,
// 		"sentence_id": body.SentenceID,
// 	})
// }

package v1

import (
	//bc "application/blockchain"
	"application/pkg/app"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SentenceRequestBody 用于接收句子、标签、所属数据集编号及句子自身编号的请求体
type SentenceRequestBody struct {
	Sentence   string `json:"sentence" binding:"required"`    // 句子
	Label      bool   `json:"label" binding:"required"`       // 标签
	DatasetID  string `json:"dataset_id" binding:"required"`  // 所属数据集编号
	SentenceID string `json:"sentence_id" binding:"required"` // 句子自身编号
}

func UploadSentence(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SentenceRequestBody)

	// 解析请求体
	if err := c.ShouldBindJSON(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	// 文件名格式为 {DatasetID}_{SentenceID}.json
	fileName := fmt.Sprintf("%s_%s.json", body.DatasetID, body.SentenceID)
	filePath := filepath.Join("data", fileName) // 存储在 data 目录下

	// 创建 data 目录（如果不存在）
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建目录失败: %s", err.Error()))
		return
	}

	// 打开文件进行写入
	file, err := os.Create(filePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("无法创建文件: %s", err.Error()))
		return
	}
	defer file.Close()

	// 将数据写入文件
	if err := json.NewEncoder(file).Encode(body); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入文件失败: %s", err.Error()))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", gin.H{
		"file": fileName,
	})
}

// // SetRequestBody 用于接收账户ID和数据集ID的请求体
// type SetRequestBody struct {
// 	AccountID string `json:"account_id" binding:"required"`  // 账户ID
// 	DatasetID string `json:"dataset_id" binding:"required"`  // 数据集ID
// }

// func UploadSet(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	body := new(SetRequestBody)

// 	// 解析请求体
// 	if err := c.ShouldBindJSON(body); err != nil {
// 		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
// 		return
// 	}

// 	// 文件名格式为 {AccountID}_{DatasetID}.json
// 	fileName := fmt.Sprintf("%s_%s.json", body.AccountID, body.DatasetID)
// 	filePath := filepath.Join("data", fileName) // 存储在 data 目录下

// 	// 创建 data 目录（如果不存在）
// 	if err := os.MkdirAll("data", os.ModePerm); err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建目录失败: %s", err.Error()))
// 		return
// 	}

// 	// 打开文件进行写入
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("无法创建文件: %s", err.Error()))
// 		return
// 	}
// 	defer file.Close()

// 	// 将数据写入文件
// 	if err := json.NewEncoder(file).Encode(body); err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入文件失败: %s", err.Error()))
// 		return
// 	}

// 	// 成功响应
// 	appG.Response(http.StatusOK, "成功", gin.H{
// 		"file": fileName,
// 	})
// }

// SetRequestBody 用于接收账户ID和数据集ID的请求体
type SetRequestBody struct {
	AccountID string `json:"account_id" binding:"required"`  // 账户ID
	DatasetID string `json:"dataset_id" binding:"required"`  // 数据集ID
}

type Record struct {
	AccountID string `json:"account_id"`
	DatasetID string `json:"dataset_id"`
}

// 读取现有记录
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

// 写入记录到文件
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
		AccountID: body.AccountID,
		DatasetID: body.DatasetID,
	}
	records = append(records, newRecord)

	// 写入记录到文件
	if err := writeRecords(filePath, records); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入记录时发生错误: %s", err.Error()))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", gin.H{
		"file": filePath,
	})
}