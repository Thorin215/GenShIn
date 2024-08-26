package model

// User 用户
type User struct {
	ID   string `json:"id"`   // 用户ID
	Name string `json:"name"` // 用户名
}

// DatasetFile 数据集文件
type DatasetFile struct {
	Hash     string `json:"hash"`     // 文件哈希 (key)
	FileName string `json:"filename"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
}

/*
元数据不上链，方便数据库查询
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
*/

// DatasetVersion 数据集一个版本
type DatasetVersion struct {
	CreationTime string   `json:"creation_time"` // 创建时间
	ChangeLog    string   `json:"change_log"`    // 版本说明
	Size         int64    `json:"size"`          // 版本大小
	Rows         int32    `json:"rows"`          // 行数
	Files        []string `json:"files"`         // 文件哈希列表
}

// Dataset 数据集
type Dataset struct {
	Owner     string `json:"owner"`     // 所有者ID
	Name      string `json:"name"`      // 数据集名
	Downloads int32  `json:"downloads"` // 下载次数
	Stars     int32  `json:"stars"`     // 收藏次数
	// Metadata  DatasetMetadata  `json:"metadata"`  // 元数据
	Versions []DatasetVersion `json:"versions"` // 版本列表
}

// DownloadRecord 下载记录
type DownloadRecord struct {
	DatasetOwner string   `json:"dataset_owner"` // 数据集所有者
	DatasetName  string   `json:"dataset_name"`  // 数据集名
	User         string   `json:"user"`          // 下载者ID
	Files        []string `json:"files"`         // 文件哈希列表
	Time         string   `json:"time"`          // 下载时间
}

const (
	UserKey                  = "user"
	DatasetFileKey           = "dataset-file"
	DatasetKey               = "dataset"
	DownloadRecordUserKey    = "download-record-user"
	DownloadRecordDatasetKey = "download-record-dataset"
)
