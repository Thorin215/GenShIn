package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"encoding/json"
	"fmt"
	"net/http"
    "crypto/sha256"
	"encoding/hex"
	"io"
	"time"
	"os"
    "io/ioutil"
	"github.com/gin-gonic/gin"
)

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

// SetRequestBody 用于接收账户ID和数据集的请求体
type SetRequestBody struct {
    Name          string                 `json:"Name"`
    Owner         string                 `json:"Owner"`
    CreationTime  time.Time              `json:"CreationTime"`
    Rows          int32                  `json:"Rows"`
	Metadata      DatasetMetadata        `json:"metadata"`     // 元数据
    Files        []string                `json:"files"`         // 文件哈希列表
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

func UploadSet(c *gin.Context) {
    appG := app.Gin{C: c}
    body := new(SetRequestBody)

    // 解析请求体
    if err := c.ShouldBindJSON(body); err != nil {
        appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
        return
    }

    // 构造 DatasetVersion 列表
    versions := []DatasetVersion{
        {
            CreationTime: body.CreationTime.Format(time.RFC3339),
            ChangeLog:    "initial version",
            Rows:         body.Rows,
            Files:        []string{},
        },
    }

    // 调用链码创建数据集
    args := [][]byte{
        []byte(body.Owner),
        []byte(body.Name),
        []byte(ToJson(versions)), // 将 JSON 字符串转换为 []byte
    }

    createResponse, err := bc.ChannelExecute("createDataset", args)
    if err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
        return
    }

    // 检查链码响应
    payload := string(createResponse.Payload)
    if payload != "" {
        appG.Response(http.StatusInternalServerError, "失败", payload)
        return
    }

    // 将数据集信息写入本地 JSON 文件
    setRecord := map[string]interface{}{
        "name":     body.Name,
        "owner":    body.Owner,
        "metadata": body.Metadata,
    }

    file, err := os.OpenFile("setRecord.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("无法打开文件: %s", err.Error()))
        return
    }
    defer file.Close()

    // 将现有内容读入
    var records []map[string]interface{}
    fileContent, _ := os.ReadFile("setRecord.json")
    if len(fileContent) > 0 {
        json.Unmarshal(fileContent, &records)
    }

    // 添加新记录
    records = append(records, setRecord)

    // 写入文件
    file.Truncate(0) // 清空文件内容
    file.Seek(0, 0)  // 将文件指针移动到开头
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(records); err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入文件失败: %s", err.Error()))
        return
    }

    // 成功响应
    appG.Response(http.StatusOK, "成功", gin.H{
        "dataset_name": body.Name,
        "owner":        body.Owner,
    })
}

func ToJson(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes
}

// GetAllDataSet 处理获取所有数据集请求
func GetAllDataSet(c *gin.Context) {
	appG := app.Gin{C: c}

	// 调用链码的 QueryDatasetFullList 函数
	resp, err := bc.ChannelQuery("queryDatasetFullList", nil)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 反序列化 JSON
	var datasets []Dataset
	if err = json.Unmarshal(resp.Payload, &datasets); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", datasets)
}

// UpdateVersion 更新数据集版本
func UpdateVersion(c *gin.Context) {
	appG := app.Gin{C: c}

	// 定义请求体结构
	type UpdateVersionRequest struct {
		Owner        string   `json:"owner"`        // 所有者
		Name         string   `json:"name"`         // 数据集名称
		CreationTime string   `json:"creation_time"` // 创建时间
		ChangeLog    string   `json:"change_log"`    // 版本说明
		Rows         int32    `json:"rows"`          // 行数
		Files        []string `json:"files"`         // 文件哈希列表
	}

	body := new(UpdateVersionRequest)
	if err := c.ShouldBindJSON(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	// 准备链码调用参数
	version := DatasetVersion{
		CreationTime: body.CreationTime,
		ChangeLog:    body.ChangeLog,
		Rows:         body.Rows,
		Files:        body.Files,
	}
	versionBytes, err := json.Marshal(version)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("版本序列化出错: %s", err.Error()))
		return
	}

	// 调用链码的 AppendDatasetVersion 函数
	resp, err := bc.ChannelExecute("appendDatasetVersion", [][]byte{
		[]byte(body.Owner),
		[]byte(body.Name),
		versionBytes,
	})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用链码出错: %s", err.Error()))
        fmt.Println(resp)
		return
    }

	// 成功响应
	appG.Response(http.StatusOK, "成功", gin.H{
		"dataset_name": body.Name,
		"owner":        body.Owner,
	})
}

// GetDatasetMetadata 查询本地 setRecord.json 获取元数据
func GetDatasetMetadata(c *gin.Context) {
    appG := app.Gin{C: c}
    var requestBody struct {
        Owner string `json:"owner"`
        Name  string `json:"name"`
    }

    // 解析请求体
    if err := c.BindJSON(&requestBody); err != nil {
        appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("解析请求体失败: %s", err.Error()))
        return
    }

    owner := requestBody.Owner
    name := requestBody.Name

    if owner == "" || name == "" {
        appG.Response(http.StatusBadRequest, "失败", "所有者和数据集名称不能为空")
        return
    }

    // 读取本地 setRecord.json 文件
    filePath := "setRecord.json" // 更新为 setRecord.json 文件的实际路径
    fileContent, err := ioutil.ReadFile(filePath)
    if err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取文件失败: %s", err.Error()))
        return
    }

    // 定义用于解析的临时结构体
    var datasets []map[string]interface{}

    // 解析 JSON 文件内容
    err = json.Unmarshal(fileContent, &datasets)
    if err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析文件内容失败: %s", err.Error()))
        return
    }

    // 查找匹配的 Dataset
    for _, dataset := range datasets {
        if dataset["owner"] == owner && dataset["name"] == name {
            metadata, ok := dataset["metadata"].(map[string]interface{})
            if !ok {
                appG.Response(http.StatusInternalServerError, "失败", "元数据格式不正确")
                return
            }

            // 转换为 DatasetMetadata 结构体
            var result DatasetMetadata
            metadataBytes, err := json.Marshal(metadata)
            if err != nil {
                appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("转换元数据失败: %s", err.Error()))
                return
            }

            err = json.Unmarshal(metadataBytes, &result)
            if err != nil {
                appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析元数据失败: %s", err.Error()))
                return
            }

            appG.Response(http.StatusOK, "成功", result)
            return
        }
    }

    appG.Response(http.StatusNotFound, "失败", "数据集未找到")
}

type DatasetFile struct {
	Hash     string `json:"hash"`     // 文件哈希 (key)
	FileName string `json:"filename"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
}

func UploadFile(c *gin.Context) {
	appG := app.Gin{C: c}
    fmt.Println("UploadFile")
	// 获取文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("无法获取文件: %s", err.Error()))
		return
	}
	defer file.Close()

	// 计算文件的 SHA-256 哈希值
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("计算哈希值出错: %s", err.Error()))
		return
	}
	hashSum := hash.Sum(nil)
	hashString := hex.EncodeToString(hashSum)

	// 获取文件大小
	fileSize := header.Size

	// 创建 DatasetFile 对象
	datasetFile := DatasetFile{
		Hash:     hashString,
		FileName: header.Filename,
		Size:     fileSize,
	}

	// 将文件信息存储到链上
	args := [][]byte{
		[]byte(datasetFile.FileName),
		[]byte(fmt.Sprintf("%d", datasetFile.Size)),
		[]byte(datasetFile.Hash),
	}

	createResponse, err := bc.ChannelExecute("createFile", args)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 检查链码响应
	payload := string(createResponse.Payload)
	if payload != "" {
		appG.Response(http.StatusInternalServerError, "失败", payload)
		return
	}

	// 返回结果
	appG.Response(http.StatusOK, "成功", datasetFile)
}