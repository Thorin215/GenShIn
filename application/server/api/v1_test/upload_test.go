package v1_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"application/routers"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	// 创建测试路由
	router := routers.InitRouter()

	// 创建一个新的 buffer 作为 body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 添加文件字段到 multipart 表单中
	filePath := "/path/to/your/file.txt"
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	// 结束 multipart 写入
	err = writer.Close()
	assert.NoError(t, err)

	// 创建一个新的 HTTP 请求
	req, _ := http.NewRequest("POST", "/api/v1/uploadFile", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// 记录响应
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "成功")
}
