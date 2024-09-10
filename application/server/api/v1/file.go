package v1

import (
	bc "application/blockchain"
	"application/model"
	"application/pkg/app"
	"application/pkg/utils"
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	appG := app.Gin{C: c}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("获取文件出错: %s", err.Error()))
		return
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("计算哈希值出错: %s", err.Error()))
		return
	}
	hashSum := hash.Sum(nil)
	hashString := hex.EncodeToString(hashSum)

	// 检查链上文件是否已存在。如果存在，不再重复上传，直接返回
	_, err = bc.ChannelQuery("queryFile", [][]byte{[]byte(hashString)})
	if err == nil {
		appG.Response(http.StatusOK, "成功", hashString)
		return
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
	filePath := filepath.Join(dir, hashString)
	outFile, err := os.Create(filePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("打开文件出错: %s", err.Error()))
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
		[]byte(hashString),
		[]byte(fmt.Sprintf("%d", header.Size)),
	}

	_, err = bc.ChannelExecute("createFile", args)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 返回结果
	appG.Response(http.StatusOK, "成功", hashString)
}

func DownloadFile(c *gin.Context) {
	appG := app.Gin{C: c}

	var body struct {
		File         model.DatasetFile `json:"file"`
		DatasetOwner string            `json:"dataset_owner"`
		DatasetName  string            `json:"dataset_name"`
		User         string            `json:"user"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
		return
	}

	hash := body.File.Hash
	fileName := body.File.FileName

	// 检查 hash 是否为 SHA-256
	if ok, _ := regexp.MatchString("^[a-f0-9]{64}$", hash); !ok {
		appG.Response(http.StatusBadRequest, "失败", "文件哈希格式错误")
		return
	}

	// 检查本地文件是否存在
	filePath := filepath.Join("data/Files", hash)
	if _, err := os.Stat(filePath); err != nil {
		appG.Response(http.StatusNotFound, "失败", "文件不存在")
		return
	}

	// 将 [File] 转换为 JSON 字符串
	file, err := json.Marshal([]model.DatasetFile{body.File})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("序列化出错: %s", err.Error()))
		return
	}

	// 上传下载记录
	args := [][]byte{
		[]byte(body.DatasetOwner),     // args[0]: 所有者ID | string
		[]byte(body.DatasetName),      // args[1]: 数据集名 | string
		[]byte(body.User),             // args[2]: 下载者ID | string
		[]byte(file),                  // args[3]: 文件信息 []DatasetFile | string (JSON)
		[]byte(utils.GetTimeString()), // args[4]: 下载时间 | string
	}

	_, err = bc.ChannelExecute("createRecord", args)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	appG.C.FileAttachment(filePath, fileName)
}

func DownloadFilesCompressed(c *gin.Context) {
	appG := app.Gin{C: c}

	var body struct {
		Files        []model.DatasetFile `json:"files"`
		ZipName      string              `json:"zipname"`
		DatasetOwner string              `json:"dataset_owner"`
		DatasetName  string              `json:"dataset_name"`
		User         string              `json:"user"`
	}

	// 解析 Body 参数
	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	re := regexp.MustCompile("^[a-f0-9]{64}$")
	for _, file := range body.Files {
		// 检查 hash 是否为 SHA-256
		if ok := re.MatchString(file.Hash); !ok {
			appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("文件哈希格式错误: %s", file.Hash))
			return
		}

		// 检查本地文件是否存在
		filePath := filepath.Join("data", "Files", file.Hash)
		if _, err := os.Stat(filePath); err != nil {
			appG.Response(http.StatusNotFound, "失败", fmt.Sprintf("文件不存在: %s", file.Hash))
			return
		}
	}

	// 将 Files 转换为 JSON 字符串
	files, err := json.Marshal(body.Files)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("序列化出错: %s", err.Error()))
		return
	}

	// 上传下载记录
	args := [][]byte{
		[]byte(body.DatasetOwner),     // args[0]: 所有者ID | string
		[]byte(body.DatasetName),      // args[1]: 数据集名 | string
		[]byte(body.User),             // args[2]: 下载者ID | string
		[]byte(files),                 // args[3]: 文件列表 []DatasetFile | string (JSON)
		[]byte(utils.GetTimeString()), // args[4]: 下载时间 | string
	}

	_, err = bc.ChannelExecute("createRecord", args)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 1. 创建一个临时压缩文件
	zipFile, err := os.CreateTemp("data", "download-*.zip")
	filePath := zipFile.Name()
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建压缩文件出错: %s", err.Error()))
		return
	}

	// 2. 将所有文件添加到压缩文件
	zipWriter := zip.NewWriter(zipFile)
	for _, file := range body.Files {
		fileReader, err := os.Open(filepath.Join("data", "Files", file.Hash))
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("打开文件出错: %s", err.Error()))
			return
		}

		zipFileWriter, err := zipWriter.Create(file.FileName)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("添加压缩文件出错: %s", err.Error()))
			return
		}

		if _, err := io.Copy(zipFileWriter, fileReader); err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入压缩文件出错: %s", err.Error()))
			return
		}
		fileReader.Close()
	}
	zipWriter.Close()

	// 3. 返回压缩文件
	var filePathAbs string
	if filePathAbs, err = filepath.Abs(filePath); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("获取文件路径出错: %s", err.Error()))
		return
	}
	appG.C.FileAttachment(filePathAbs, fmt.Sprintf("%s.zip", body.ZipName))
	defer os.Remove(filePath)
}
