package model

// User 用户
type User struct {
	ID   string `json:"id"`   // 用户ID
	Name string `json:"name"` // 用户名
}

// File 文件
type File struct {
	Hash           string `json:"hash"`            // 文件哈希 (key)
	Size           int64  `json:"size"`            // 文件大小
	ReferenceCount int32  `json:"reference_count"` // 引用计数
}

// DatasetFile 数据集文件
type DatasetFile struct {
	Hash     string `json:"hash"`     // 文件哈希
	FileName string `json:"filename"` // 文件名
}

// Version 数据集的一个版本
type Version struct {
	Files        []DatasetFile `json:"files"`         // 文件列表
	Rows         int32         `json:"rows"`          // 行数
	CreationTime string        `json:"creation_time"` // 创建时间
	ChangeLog    string        `json:"change_log"`    // 版本说明
}

// Dataset 数据集
type Dataset struct {
	Owner    string    `json:"owner"`    // 所有者ID
	Name     string    `json:"name"`     // 数据集名
	Versions []Version `json:"versions"` // 版本列表
	Deleted  bool      `json:"deleted"`  // 已删除
}

// Record 下载记录
type Record struct {
	DatasetOwner string        `json:"dataset_owner"` // 数据集所有者
	DatasetName  string        `json:"dataset_name"`  // 数据集名
	User         string        `json:"user"`          // 下载者ID
	Files        []DatasetFile `json:"files"`         // 文件列表
	Time         string        `json:"time"`          // 下载时间
}

const (
	UserKey          = "user"
	FileKey          = "file"
	DatasetKey       = "dataset"
	RecordUserKey    = "record-user"
	RecordDatasetKey = "record-dataset"
)
