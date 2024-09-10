package sql

import (
	"application/model"
	// "encoding/json"

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
	Owner      string   `gorm:"primaryKey" json:"owner"`
	Name       string   `gorm:"primaryKey" json:"name"`
	Tasks      []string `json:"tasks" gorm:"type:json"`      // 以 JSON 格式存储
	Modalities []string `json:"modalities" gorm:"type:json"` // 以 JSON 格式存储
	Formats    []string `json:"formats" gorm:"type:json"`    // 以 JSON 格式存储
	SubTasks   []string `json:"sub_tasks" gorm:"type:json"`  // 以 JSON 格式存储
	Languages  []string `json:"languages" gorm:"type:json"`  // 以 JSON 格式存储
	Libraries  []string `json:"libraries" gorm:"type:json"`  // 以 JSON 格式存储
	Tags       []string `json:"tags" gorm:"type:json"`       // 以 JSON 格式存储
	License    string   `json:"license" gorm:"type:json"`    // 以 JSON 格式存储
}

func MigrateMetadata(db *gorm.DB) error {
	// 确保创建 MetadataTable 表
	err := db.AutoMigrate(&MetadataTable{})
	if err != nil {
		return err
	}
	return nil
}

func CreateMetaData(metadataBody *MetadataBody) error {
	// 创建 MetadataTable 实例并映射 MetadataBody 的各字段
	tableData := MetadataTable{
		Owner:      metadataBody.Owner,
		Name:       metadataBody.Name,
		Tasks:      metadataBody.Metadata.Tasks,
		Modalities: metadataBody.Metadata.Modalities,
		Formats:    metadataBody.Metadata.Formats,
		SubTasks:   metadataBody.Metadata.SubTasks,
		Languages:  metadataBody.Metadata.Languages,
		Libraries:  metadataBody.Metadata.Libraries,
		Tags:       metadataBody.Metadata.Tags,
		License:    metadataBody.Metadata.License,
	}

	// 存入数据库
	result := DB.Create(&tableData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetMetaData(name string, owner string) (*MetadataBody, error) {
	var metadataTable MetadataTable
	result := DB.Where("owner = ? AND name = ?", owner, name).First(&metadataTable)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 记录未找到
		}
		return nil, result.Error // 其他错误
	}

	// 将 MetadataTable 的字段映射回 Metadata
	metadata := model.Metadata{
		Tasks:      metadataTable.Tasks,
		Modalities: metadataTable.Modalities,
		Formats:    metadataTable.Formats,
		SubTasks:   metadataTable.SubTasks,
		Languages:  metadataTable.Languages,
		Libraries:  metadataTable.Libraries,
		Tags:       metadataTable.Tags,
		License:    metadataTable.License,
	}

	return &MetadataBody{
		Owner:    metadataTable.Owner,
		Name:     metadataTable.Name,
		Metadata: metadata,
	}, nil
}
