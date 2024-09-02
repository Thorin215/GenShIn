package log

import (
	"application/sql"
	"time"
)

type Log struct {
	LogID     int       `json:"log_id" gorm:"primary_key"`
	DataSetID int       `json:"dataset_id"`
	OpMsg     string    `json:"op_msg"`
	Status    string    `json:"status"`
	TimeStamp time.Time `json:"time_stamp"`
}

func CreateLog(log *Log) error {
	return sql.DB.Create(log).Error
}

func UpdateLog(log *Log) error {
	return sql.DB.Save(log).Error
}

func GetLog(dataSetID int) ([]Log, error) {
	logs := []Log{}
	err := sql.DB.Where("dataset_id = ?", dataSetID).Find(&logs).Error
	return logs, err
}
