package sql

import (
	"application/model"
	"encoding/json"

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
	Owner    string          `gorm:"primaryKey" json:"owner"`
	Name     string          `gorm:"primaryKey" json:"name"`
	Metadata json.RawMessage `json:"metadata" gorm:"type:json"` // 存储 JSON 数据
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
	// 将 Metadata 转换为 JSON 数据
	metadataJSON, err := json.Marshal(metadataBody.Metadata)
	if err != nil {
		return err
	}

	// 创建 MetadataTable 实例
	tableData := MetadataTable{
		Owner:    metadataBody.Owner,
		Name:     metadataBody.Name,
		Metadata: metadataJSON,
	}

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

	// 将 JSON 数据解析为 Metadata
	var metadata model.Metadata
	if err := json.Unmarshal(metadataTable.Metadata, &metadata); err != nil {
		return nil, err
	}

	return &MetadataBody{
		Owner:    metadataTable.Owner,
		Name:     metadataTable.Name,
		Metadata: metadata,
	}, nil
}
