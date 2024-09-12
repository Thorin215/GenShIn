// package sql

// import (
// 	"application/model"
// 	// "encoding/json"

// 	"gorm.io/gorm"
// )

// // MetadataBody 用于返回数据
// type MetadataBody struct {
// 	Owner    string         `json:"owner"`
// 	Name     string         `json:"name"`
// 	Metadata model.Metadata `json:"metadata"` // 使用 Metadata 类型
// }

// // MetadataTable 数据库表结构
// type MetadataTable struct {
// 	Owner      string   `gorm:"primaryKey" json:"owner"`
// 	Name       string   `gorm:"primaryKey" json:"name"`
// 	Tasks      []string `json:"tasks" gorm:"type:json"`      // 以 JSON 格式存储
// 	Modalities []string `json:"modalities" gorm:"type:json"` // 以 JSON 格式存储
// 	Formats    []string `json:"formats" gorm:"type:json"`    // 以 JSON 格式存储
// 	SubTasks   []string `json:"sub_tasks" gorm:"type:json"`  // 以 JSON 格式存储
// 	Languages  []string `json:"languages" gorm:"type:json"`  // 以 JSON 格式存储
// 	Libraries  []string `json:"libraries" gorm:"type:json"`  // 以 JSON 格式存储
// 	Tags       []string `json:"tags" gorm:"type:json"`       // 以 JSON 格式存储
// 	License    string   `json:"license"`                     // 普通字符串
// }

// func MigrateMetadata(db *gorm.DB) error {
// 	// 确保创建 MetadataTable 表
// 	err := db.AutoMigrate(&MetadataTable{})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func CreateMetaData(metadataBody *MetadataBody) error {
// 	// 创建 MetadataTable 实例并映射 MetadataBody 的各字段
// 	tableData := MetadataTable{
// 		Owner:      metadataBody.Owner,
// 		Name:       metadataBody.Name,
// 		Tasks:      metadataBody.Metadata.Tasks,
// 		Modalities: metadataBody.Metadata.Modalities,
// 		Formats:    metadataBody.Metadata.Formats,
// 		SubTasks:   metadataBody.Metadata.SubTasks,
// 		Languages:  metadataBody.Metadata.Languages,
// 		Libraries:  metadataBody.Metadata.Libraries,
// 		Tags:       metadataBody.Metadata.Tags,
// 		License:    metadataBody.Metadata.License,
// 	}

// 	// 存入数据库
// 	result := DB.Create(&tableData)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// func GetMetaData(name string, owner string) (*MetadataBody, error) {
// 	var metadataTable MetadataTable
// 	result := DB.Where("owner = ? AND name = ?", owner, name).First(&metadataTable)
// 	if result.Error != nil {
// 		if result.Error == gorm.ErrRecordNotFound {
// 			return nil, nil // 记录未找到
// 		}
// 		return nil, result.Error // 其他错误
// 	}

// 	// 将 MetadataTable 的字段映射回 Metadata
// 	metadata := model.Metadata{
// 		Tasks:      metadataTable.Tasks,
// 		Modalities: metadataTable.Modalities,
// 		Formats:    metadataTable.Formats,
// 		SubTasks:   metadataTable.SubTasks,
// 		Languages:  metadataTable.Languages,
// 		Libraries:  metadataTable.Libraries,
// 		Tags:       metadataTable.Tags,
// 		License:    metadataTable.License,
// 	}

// 	return &MetadataBody{
// 		Owner:    metadataTable.Owner,
// 		Name:     metadataTable.Name,
// 		Metadata: metadata,
// 	}, nil
// }

package sql

import (
	"application/model"
	"strings"

	"gorm.io/gorm"
)

// MetadataBody 用于返回数据
type MetadataBody struct {
	Owner    string         `json:"owner"`
	Name     string         `json:"name"`
	Metadata model.Metadata `json:"metadata"` // 使用 Metadata 类型
}

// MetadataTable 数据库表结构
type MetadataTable struct {
	Owner      string `gorm:"primaryKey" json:"owner"`
	Name       string `gorm:"primaryKey" json:"name"`
	Tasks      string `json:"tasks"`      // 存储为字符串
	Modalities string `json:"modalities"` // 存储为字符串
	Formats    string `json:"formats"`    // 存储为字符串
	SubTasks   string `json:"sub_tasks"`  // 存储为字符串
	Languages  string `json:"languages"`  // 存储为字符串
	Libraries  string `json:"libraries"`  // 存储为字符串
	Tags       string `json:"tags"`       // 存储为字符串
	License    string `json:"license"`    // 普通字符串
}

type QueryMetadataResp struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// MigrateMetadata 确保创建 MetadataTable 表
func MigrateMetadata(db *gorm.DB) error {
	err := db.AutoMigrate(&MetadataTable{})
	if err != nil {
		return err
	}
	return nil
}

// CreateMetaData 插入数据
func CreateMetaData(metadataBody *MetadataBody) error {
	// 将 []string 转换为字符串
	tasksStr := strings.Join(metadataBody.Metadata.Tasks, ",")
	modalitiesStr := strings.Join(metadataBody.Metadata.Modalities, ",")
	formatsStr := strings.Join(metadataBody.Metadata.Formats, ",")
	subTasksStr := strings.Join(metadataBody.Metadata.SubTasks, ",")
	languagesStr := strings.Join(metadataBody.Metadata.Languages, ",")
	librariesStr := strings.Join(metadataBody.Metadata.Libraries, ",")
	tagsStr := strings.Join(metadataBody.Metadata.Tags, ",")

	// 创建 MetadataTable 实例并映射 MetadataBody 的各字段
	tableData := MetadataTable{
		Owner:      metadataBody.Owner,
		Name:       metadataBody.Name,
		Tasks:      tasksStr,
		Modalities: modalitiesStr,
		Formats:    formatsStr,
		SubTasks:   subTasksStr,
		Languages:  languagesStr,
		Libraries:  librariesStr,
		Tags:       tagsStr,
		License:    metadataBody.Metadata.License,
	}

	// 存入数据库
	result := DB.Create(&tableData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetMetaData 从数据库获取数据
func GetMetaData(name string, owner string) (*MetadataBody, error) {
	var metadataTable MetadataTable
	result := DB.Where("owner = ? AND name = ?", owner, name).First(&metadataTable)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 记录未找到
		}
		return nil, result.Error // 其他错误
	}

	// 将字符串拆分为 []string
	tasks := strings.Split(metadataTable.Tasks, ",")
	modalities := strings.Split(metadataTable.Modalities, ",")
	formats := strings.Split(metadataTable.Formats, ",")
	subTasks := strings.Split(metadataTable.SubTasks, ",")
	languages := strings.Split(metadataTable.Languages, ",")
	libraries := strings.Split(metadataTable.Libraries, ",")
	tags := strings.Split(metadataTable.Tags, ",")

	// 将 MetadataTable 的字段映射回 Metadata
	metadata := model.Metadata{
		Tasks:      tasks,
		Modalities: modalities,
		Formats:    formats,
		SubTasks:   subTasks,
		Languages:  languages,
		Libraries:  libraries,
		Tags:       tags,
		License:    metadataTable.License,
	}

	return &MetadataBody{
		Owner:    metadataTable.Owner,
		Name:     metadataTable.Name,
		Metadata: metadata,
	}, nil
}

// QueryMetadata 根据输入查询数据
func QueryMetadata(metadata model.Metadata) ([]QueryMetadataResp, error) {
	var results []MetadataTable

	query := DB.Model(&MetadataTable{})

	if metadata.Tasks != nil {
		query = query.Where("tasks IN ?", metadata.Tasks)
	}
	if metadata.Modalities != nil {
		query = query.Where("modalities IN ?", metadata.Modalities)
	}
	if metadata.Formats != nil {
		query = query.Where("formats IN ?", metadata.Formats)
	}
	if metadata.SubTasks != nil {
		query = query.Where("sub_tasks IN ?", metadata.SubTasks)
	}
	if metadata.Languages != nil {
		query = query.Where("languages IN ?", metadata.Languages)
	}
	if metadata.Libraries != nil {
		query = query.Where("libraries IN ?", metadata.Libraries)
	}
	if metadata.Tags != nil {
		query = query.Where("tags IN ?", metadata.Tags)
	}
	if metadata.License != "" {
		query = query.Where("license = ?", metadata.License)
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	var resp []QueryMetadataResp
	for _, result := range results {
		resp = append(resp, QueryMetadataResp{
			Owner: result.Owner,
			Name:  result.Name,
		})
	}

	return resp, nil
}
