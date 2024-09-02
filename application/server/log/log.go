package log

import (
	"application/sql"
	"time"
)

type Log struct {
	LogID     int32     `json:"log_id" gorm:"primary_key"`
	DataSetID int32     `json:"dataset_id"`
	ChangeLog string    `json:"change_log"`
	TimeStamp time.Time `json:"time_stamp"`
}

func CreateLog(log *Log) error {
	return sql.DB.Create(log).Error
}

func UpdateLog(log *Log) error {
	return sql.DB.Updates(log).Error
}

func GetLog(dataSetID int32) ([]Log, error) {
	logs := []Log{}
	err := sql.DB.Where("data_set_id = ?", dataSetID).Find(&logs).Error
	return logs, err
}
