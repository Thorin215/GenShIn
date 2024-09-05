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
    "archive/zip"
    // "bytes"
    "io/ioutil"
	"github.com/gin-gonic/gin"
    "path/filepath"
    // "mime"
    "encoding/base64"
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
    fmt.Println(body.Name)
    fmt.Println(body.Files)
    // 构造 DatasetVersion 列表
    versions := []DatasetVersion{
        {
            CreationTime: body.CreationTime.Format(time.RFC3339),
            ChangeLog:    "Initial Version",
            Rows:         body.Rows,
            Files:        body.Files,
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

    file, err := os.OpenFile("data/setRecord.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("无法打开文件: %s", err.Error()))
        return
    }
    defer file.Close()

    // 将现有内容读入
    var records []map[string]interface{}
    fileContent, _ := os.ReadFile("data/setRecord.json")
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
    filePath := "data/setRecord.json" // 更新为 setRecord.json 文件的实际路径
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

    // 重新打开文件
    file.Seek(0, io.SeekStart)
    
    // 创建目录（如果不存在）
    dir := "data/Files"
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建目录出错: %s", err.Error()))
        return
    }

    // 保存文件到本地路径
    filePath := filepath.Join(dir, header.Filename)
    outFile, err := os.Create(filePath)
    if err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("保存文件出错: %s", err.Error()))
        return
    }
    defer outFile.Close()

    // 将文件内容写入本地
    if _, err := io.Copy(outFile, file); err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入文件出错: %s", err.Error()))
        return
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

func DownloadDataSet(c *gin.Context) {
    appG := app.Gin{C: c}
    body := new(struct {
        Files  []string `json:"files"`
        Name   string   `json:"name"`
        Owner  string   `json:"owner"`
    })

    // 解析 Body 参数
    if err := c.ShouldBindJSON(body); err != nil {
        appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
        return
    }

    // 使用 ToJson 函数将文件哈希值数组序列化为 JSON 字符串
    bodyBytes := ToJson(body.Files)

    // 调用智能合约查询文件信息
    resp, err := bc.ChannelQuery("queryMultipleFiles", [][]byte{bodyBytes})
    if err != nil {
        fmt.Printf("调用智能合约出错: %s", err.Error())
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
        return
    }

    // 反序列化 JSON
    var files []DatasetFile
    if err = json.Unmarshal(resp.Payload, &files); err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
        return
    }

    // 创建一个临时文件用于存储压缩包
    tmpFile, err := os.CreateTemp("", "dataset-*.zip")
    if err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建临时文件失败: %s", err.Error()))
        return
    }
    defer os.Remove(tmpFile.Name())

    // 创建一个 zip.Writer 对象
    zipWriter := zip.NewWriter(tmpFile)
    defer zipWriter.Close()

    // 压缩文件
    for _, file := range files {
        fileName := file.FileName
        if fileName == "" {
            appG.Response(http.StatusBadRequest, "失败", "文件名不能为空")
            return
        }

        // 构建文件路径
        filePath := fmt.Sprintf("data/Files/%s", fileName)

        // 打开文件
        f, err := os.Open(filePath)
        if err != nil {
            appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("无法读取文件: %s", err.Error()))
            return
        }
        defer f.Close()

        // 创建一个 zip 文件
        zipFile, err := zipWriter.Create(fileName)
        if err != nil {
            appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建 zip 文件失败: %s", err.Error()))
            return
        }

        // 将文件内容写入 zip 文件
        _, err = io.Copy(zipFile, f)
        if err != nil {
            appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入 zip 文件失败: %s", err.Error()))
            return
        }
    }

    // 关闭 zip.Writer，确保所有内容都写入了压缩包
    if err := zipWriter.Close(); err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("关闭 zip.Writer 失败: %s", err.Error()))
        return
    }

    // 读取压缩包的内容
    tmpFile.Seek(0, io.SeekStart)
    zipContent, err := io.ReadAll(tmpFile)
    if err != nil {
        appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取压缩包内容失败: %s", err.Error()))
        return
    }

    // Base64 编码压缩包内容
    encodedZipContent := base64.StdEncoding.EncodeToString(zipContent)

    // 返回 JSON 响应
    appG.Response(http.StatusOK, "成功", map[string]interface{}{
        "files": []map[string]string{
            {
                "filename": fmt.Sprintf("%s.zip", body.Name),
                "content":  encodedZipContent,
            },
        },
        "name":  body.Name,
        "owner": body.Owner,
    })
}