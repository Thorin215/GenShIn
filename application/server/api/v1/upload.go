package v1

import (
	//bc "application/blockchain"
	"application/pkg/app"
	// "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	//"io/ioutil"
	"application/sql"
	"io"

	"application/log"

	"github.com/gin-gonic/gin"
)

// DatasetMetadata 数据集元数据
type DatasetMetadata struct {
	Tasks      []string `json:"tasks"`      // 目标任务
	Modalities []string `json:"modalities"` // 数据模态
	Formats    []string `json:"formats"`    // 文件格式
	SubTasks   []string `json:"subTasks"`   // 子任务
	Languages  []string `json:"languages"`  // 语言
	Libraries  []string `json:"libraries"`  // 适用库
	Tags       []string `json:"tags"`       // 标签
	License    string   `json:"license"`    // 许可证

}

// 存入数据库的model DatasetMetadata
type MetadataToDB struct {
	DataSetID    int32     `json:"data_set_id" gorm:"primary_key"`
	Owner        string    `json:"owner"`
	Rows         int32     `json:"rows"`
	Name         string    `json:"name" gorm:"not null"` // 数据集名称
	Tasks        string    `json:"tasks"`                // 目标任务
	Modalities   string    `json:"modalities"`           // 数据模态
	Formats      string    `json:"formats"`              // 文件格式
	SubTasks     string    `json:"subTasks"`             // 子任务
	Languages    string    `json:"languages"`            // 语言
	Libraries    string    `json:"libraries"`            // 适用库
	Tags         string    `json:"tags"`                 // 标签
	License      string    `json:"license"`              // 许可证
	CreationTime time.Time `json:"creation_time"`
}

// 存入数据库的model DataSetToDB
type DataSetToDB struct {
	DataSetID     int32     `json:"data_set_id" gorm:"primary_key"`
	Data          []byte    `json:"data"`
	ModeifiedTime time.Time `json:"modified_time"`
}

// SetRequestBody 用于接收账户ID和数据集的请求体
type SetRequestBody struct {
	Name         string          `json:"Name"`
	Owner        string          `json:"Owner"`
	CreationTime time.Time       `json:"CreationTime"`
	Rows         int32           `json:"Rows"`
	Metadata     DatasetMetadata `json:"metadata"` // 元数据
}

// Dataset 数据集
type Dataset struct {
	Owner    string           `json:"owner"`    // 所有者ID
	Name     string           `json:"name"`     // 数据集名
	Metadata DatasetMetadata  `json:"metadata"` // 元数据
	Versions []DatasetVersion `json:"versions"` // 版本列表
}

// DatasetVersion 数据集一个版本
type DatasetVersion struct {
	CreationTime string   `json:"creation_time"` // 创建时间
	ChangeLog    string   `json:"change_log"`    // 版本说明
	Rows         int32    `json:"rows"`          // 行数
	Files        []string `json:"files"`         // 文件哈希列表
}

func GetAllDataSet(c *gin.Context) {
	appG := app.Gin{C: c}

	//--------------------------------------------数据库--------------------------------

	metaDatas := []MetadataToDB{}
	err := sql.DB.Find(&metaDatas).Error
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取数据集时发生错误: %s", err.Error()))
		return
	}

	//-----------------------------文件实现--------------------------------
	// dataFilePath := filepath.Join("data", "setRecord.json")

	// datasets, err := readDatasets(dataFilePath)
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取数据集时发生错误: %s", err.Error()))
	// 	return
	// }
	//-----------------------------文件实现--------------------------------
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": metaDatas,
	})
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

	// ---------------------------------数据库--------------------------------

	metaData := &MetadataToDB{
		Name:         body.Name,
		Owner:        body.Owner,
		Rows:         body.Rows,
		Tasks:        fmt.Sprintf("%v", body.Metadata.Tasks), //[]string 转 string方便存入数据库
		Modalities:   fmt.Sprintf("%v", body.Metadata.Modalities),
		Formats:      fmt.Sprintf("%v", body.Metadata.Formats),
		SubTasks:     fmt.Sprintf("%v", body.Metadata.SubTasks),
		Languages:    fmt.Sprintf("%v", body.Metadata.Languages),
		Libraries:    fmt.Sprintf("%v", body.Metadata.Libraries),
		Tags:         fmt.Sprintf("%v", body.Metadata.Tags),
		License:      body.Metadata.License,
		CreationTime: body.CreationTime,
	}

	if FindDataSet(body.Name, body.Owner) {
		appG.Response(http.StatusConflict, "失败", "数据集名称已存在")
		return
	}

	err := CreateMetaData(metaData)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入数据集时发生错误: %s", err.Error()))
		return
	}

	opLog := &log.Log{
		DataSetID: metaData.DataSetID,
		ChangeLog: "Create",
		TimeStemp: time.Now(),
	}

	log.CreateLog(opLog)

	// metaData_, err := GetMetaData(1)
	// fmt.Println("------------------")
	// fmt.Println(metaData_)
	// fmt.Println("----------------")

	// ---------------------------------数据库--------------------------------

	//----------------------------文件实现--------------------------------
	// // 文件路径
	// filePath := filepath.Join("data", "setRecord.json")

	// // 读取现有数据集
	// datasets, err := readDatasets(filePath)
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取数据集时发生错误: %s", err.Error()))
	// 	return
	// }

	// // 检查是否存在相同数据集名称的记录
	// for _, dataset := range datasets {
	// 	if dataset.Name == body.Name {
	// 		appG.Response(http.StatusConflict, "失败", "数据集名称已存在")
	// 		return
	// 	}
	// }

	// // 创建新的数据集记录
	// newDataset := Dataset{
	// 	Owner:    body.Owner,
	// 	Name:     body.Name,
	// 	Metadata: body.Metadata,
	// 	Versions: []DatasetVersion{
	// 		{
	// 			CreationTime: body.CreationTime.Format(time.RFC3339), // 格式化时间
	// 			ChangeLog:    "initial version",                      // 默认版本说明
	// 			Rows:         body.Rows,                              // 设置行数
	// 			Files:        []string{"initial"},                    // 默认文件哈希列表
	// 		},
	// 	},
	// }

	// // 添加新的记录
	// datasets = append(datasets, newDataset)

	// // 写入数据集到文件
	// if err := writeDatasets(filePath, datasets); err != nil {
	// 	appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入数据集时发生错误: %s", err.Error()))
	// 	return
	// }
	//----------------------------文件实现--------------------------------

	// 成功响应
	appG.Response(http.StatusOK, "成功", gin.H{
		"dataset_name": metaData.Name,
		"owner":        metaData.Owner,
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

// UpdateVersion 更新数据集版本
func UpdateVersion(c *gin.Context) {
	appG := app.Gin{C: c}
	type UpdateVersionRequest struct {
		Name         string   `json:"name"`          // 数据集名称
		Owner        string   `json:"owner"`         // 所有者
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

	// 文件路径
	filePath := filepath.Join("data", "setRecord.json")

	// 读取现有数据集
	datasets, err := readDatasets(filePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取数据集时发生错误: %s", err.Error()))
		return
	}

	// 查找目标数据集
	var dataset *Dataset
	for i := range datasets {
		if datasets[i].Name == body.Name && datasets[i].Owner == body.Owner {
			dataset = &datasets[i]
			break
		}
	}

	if dataset == nil {
		appG.Response(http.StatusNotFound, "失败", "数据集未找到")
		return
	}

	// 创建新的版本
	newVersion := DatasetVersion{
		CreationTime: body.CreationTime,
		ChangeLog:    body.ChangeLog,
		Rows:         body.Rows,
		Files:        body.Files,
	}

	// 更新数据集的版本列表
	dataset.Versions = append(dataset.Versions, newVersion)

	// 写入更新后的数据集到文件
	if err := writeDatasets(filePath, datasets); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入数据集时发生错误: %s", err.Error()))
		return
	}

	// 成功响应
	appG.Response(http.StatusOK, "成功", gin.H{
		"dataset_name": body.Name,
		"owner":        body.Owner,
	})
}

// func UploadMetaData(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	metaData := &DatasetMetadata{}

// 	if err := c.ShouldBindJSON(metaData); err != nil {
// 		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数错误: %s", err.Error()))
// 		return
// 	}

// 	CreateMetaData(metaData)

// 	appG.Response(http.StatusOK, "成功", gin.H{
// 		"meta_data": metaData,
// 	})
// }

func CreateMetaData(metaData *MetadataToDB) error {
	error := sql.DB.Create(metaData).Error
	return error
}

func GetMetaData(dataSetID int32) (*MetadataToDB, error) {
	metaData := &MetadataToDB{}
	err := sql.DB.Where("data_set_id = ?", dataSetID).First(metaData).Error
	if err != nil {
		return nil, err
	}
	return metaData, nil
}

func GetMetaDataByName(name string) (*MetadataToDB, error) {
	metaData := &MetadataToDB{}
	err := sql.DB.Where("name = ?", name).First(metaData).Error
	if err != nil {
		return nil, err
	}
	return metaData, nil
}

func FindDataSet(name string, owner string) bool {
	metaData := &MetadataToDB{}
	err := sql.DB.Where("name = ? AND owner = ?", name, owner).First(metaData).Error
	if err != nil {
		return false
	}
	return true
}
