package sql

import (
	"application/model"
	"strings"

	"gorm.io/gorm"
)

// MetadataBody 用于返回数据
type MetadataBody struct {
	Owner     string         `json:"owner"`
	Name      string         `json:"name"`
	Metadata  model.Metadata `json:"metadata"`  // 使用 Metadata 类型
	Downloads int            `json:"downloads"` // 下载次数
	Deleted   bool           `json:"deleted"`   // 是否删除
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
	Downloads  int    `json:"downloads"`  // 下载次数
	Deleted    bool   `json:"deleted"`    // 是否删除
}

func MigrateMetadata(db *gorm.DB) error {
	err := db.AutoMigrate(&MetadataTable{})
	if err != nil {
		return err
	}
	return nil
}

func CreateMetadata(metadataBody *MetadataBody) error {
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
		Downloads:  metadataBody.Downloads, // 设置下载次数
		Deleted:    metadataBody.Deleted,   // 设置删除标记
	}

	// 存入数据库
	result := DB.Create(&tableData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetMetadata(owner string, name string) (*MetadataBody, error) {
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
		Owner:     metadataTable.Owner,
		Name:      metadataTable.Name,
		Metadata:  metadata,
		Downloads: metadataTable.Downloads, // 返回下载次数
		Deleted:   metadataTable.Deleted,   // 返回删除标记
	}, nil
}

func IncrementDownloads(owner string, name string) error {
	result := DB.Model(&MetadataTable{}).Where("owner = ? AND name = ?", owner, name).Update("downloads", gorm.Expr("downloads + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func QueryDownloads(owner string, name string) (int, error) {
	var metadata MetadataTable
	result := DB.Where("owner = ? AND name = ?", owner, name).First(&metadata)
	if result.Error != nil {
		return 0, result.Error
	}
	return metadata.Downloads, nil
}

func MarkDeleted(owner string, name string) error {
	result := DB.Model(&MetadataTable{}).Where("owner = ? AND name = ?", owner, name).Update("deleted", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
