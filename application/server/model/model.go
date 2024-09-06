package model

import (
	"application/sql"
	"gorm.io/gorm"
	"encoding/json"
)

// DatasetMetadata 数据集元数据
type DatasetMetadata struct {
	Tasks      []string `json:"tasks" gorm:"type:json"`      // 以 JSON 格式存储
	Modalities []string `json:"modalities" gorm:"type:json"` // 以 JSON 格式存储
	Formats    []string `json:"formats" gorm:"type:json"`    // 以 JSON 格式存储
	SubTasks   []string `json:"sub_tasks" gorm:"type:json"`  // 以 JSON 格式存储
	Languages  []string `json:"languages" gorm:"type:json"`  // 以 JSON 格式存储
	Libraries  []string `json:"libraries" gorm:"type:json"`  // 以 JSON 格式存储
	Tags       []string `json:"tags" gorm:"type:json"`       // 以 JSON 格式存储
	License    string   `json:"license"`                    // 普通字段
}

// MetadataBody 用于返回数据
type MetadataBody struct {
	Name     string          `json:"name"`
	Owner    string          `json:"owner"`
	Metadata DatasetMetadata `json:"metadata"` // 使用 DatasetMetadata 类型
}

// MetadataTable 数据库表结构
type MetadataTable struct {
	Name     string          `gorm:"primaryKey" json:"name"`
	Owner    string          `gorm:"primaryKey" json:"owner"`
	Metadata json.RawMessage `json:"metadata" gorm:"type:json"` // 存储 JSON 数据
}

func InitializeDatabase(db *gorm.DB) error {
	// 确保创建 MetadataTable 表
	err := db.AutoMigrate(&MetadataTable{})
	if err != nil {
		return err
	}
	return nil
}

func CreateMetaData(metadataBody *MetadataBody) error {
	if err := InitializeDatabase(sql.DB); err != nil {
		return err
	}

	// 将 DatasetMetadata 转换为 JSON 数据
	metadataJSON, err := json.Marshal(metadataBody.Metadata)
	if err != nil {
		return err
	}

	// 创建 MetadataTable 实例
	tableData := MetadataTable{
		Name:     metadataBody.Name,
		Owner:    metadataBody.Owner,
		Metadata: metadataJSON,
	}

	result := sql.DB.Create(&tableData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetMetaData(name string, owner string) (*MetadataBody, error) {
	var metadataTable MetadataTable
	result := sql.DB.Where("name = ? AND owner = ?", name, owner).First(&metadataTable)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 记录未找到
		}
		return nil, result.Error // 其他错误
	}

	// 将 JSON 数据解析为 DatasetMetadata
	var metadata DatasetMetadata
	if err := json.Unmarshal(metadataTable.Metadata, &metadata); err != nil {
		return nil, err
	}

	return &MetadataBody{
		Name:     metadataTable.Name,
		Owner:    metadataTable.Owner,
		Metadata: metadata,
	}, nil
}
