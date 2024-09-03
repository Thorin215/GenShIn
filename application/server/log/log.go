package log

import (
	"application/sql"
	"time"
)

type Log struct {
	LogID     int32     `json:"log_id" gorm:"primaryKey"`
	DataSetID int32     `json:"data_set_id"`
	ChangeLog string    `json:"change_log"`
	TimeStemp time.Time `json:"time_stemp"`
}

func GetLog(dataSetID int32) ([]Log, error) {
	logs := []Log{}
	err := sql.DB.Where("data_set_id = ?", dataSetID).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func CreateLog(log *Log) error {
	error := sql.DB.Create(log).Error
	return error
}

func DeleteLog(logID int32) error {
	error := sql.DB.Where("log_id = ?", logID).Delete(&Log{}).Error
	return error
}

func DeleteLogByDataSetID(dataSetID int32) error {
	error := sql.DB.Where("data_set_id = ?", dataSetID).Delete(&Log{}).Error
	return error
}

func UpdateLog(log *Log) error {
	error := sql.DB.Updates(log).Error
	return error
}
