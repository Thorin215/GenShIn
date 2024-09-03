package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CreateDownloadRecord 创建下载记录
// args[0]: 数据集所有者ID string
// args[1]: 数据集名 string
// args[2]: 下载者ID string
// args[3]: 文件哈希列表 string []string as JSON
// args[4]: 下载时间 string
func CreateDownloadRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("CreateDownloadRecord-参数数量错误")
	}

	var files []string
	if err := json.Unmarshal([]byte(args[3]), &files); err != nil {
		return shim.Error(fmt.Sprintf("CreateDownloadRecord-反序列化出错: %s", err))
	}

	record := model.DownloadRecord{
		DatasetOwner: args[0],
		DatasetName:  args[1],
		User:         args[2],
		Files:        files,
		Time:         args[4],
	}

	if err := model.ValidateDownloadRecord(record); err != nil {
		return shim.Error(fmt.Sprintf("CreateDownloadRecord-参数错误: %s", err))
	}

	if exist, err := checkUserExist(stub, record.User); err != nil {
		return shim.Error(fmt.Sprintf("CreateDownloadRecord-查询用户出错: %s", err))
	} else if !exist {
		return shim.Error(fmt.Sprintf("CreateDownloadRecord-参数错误: 用户不存在: %s", record.User))
	}

	if exist, err := checkDatasetExist(stub, record.DatasetOwner, record.DatasetName); err != nil {
		return shim.Error(fmt.Sprintf("CreateDownloadRecord-查询数据集出错: %s", err))
	} else if !exist {
		return shim.Error(fmt.Sprintf("CreateDownloadRecord-参数错误: 数据集不存在: %s/%s", record.DatasetOwner, record.DatasetName))
	}

	for _, fileHash := range record.Files {
		if exist, err := checkFileExist(stub, fileHash); err != nil {
			return shim.Error(fmt.Sprintf("CreateDownloadRecord-查询文件出错: %s", err))
		} else if !exist {
			return shim.Error(fmt.Sprintf("CreateDownloadRecord-参数错误: 文件不存在: %s", fileHash))
		}
	}

	keyDataset := []string{strings.ToLower(record.DatasetOwner), strings.ToLower(record.DatasetName), strings.ToLower(record.User)}
	keyUser := []string{strings.ToLower(record.User), strings.ToLower(record.DatasetOwner), strings.ToLower(record.DatasetName)}

	if err := utils.WriteLedger(record, stub, model.DownloadRecordDatasetKey, keyDataset); err != nil {
		return shim.Error(fmt.Sprintf("CreateDownloadRecord-写入账本出错: %s", err))
	}
	if err := utils.WriteLedger(record, stub, model.DownloadRecordUserKey, keyUser); err != nil {
		return shim.Error(fmt.Sprintf("CreateDownloadRecord-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}

// QueryDownloadRecordListByUser 查询下载记录列表
// args[0]: 下载者ID string
func QueryDownloadRecordListByUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("QueryDownloadRecordListByUser-参数数量错误")
	}

	userID := strings.ToLower(args[0])
	res, err := utils.GetStateByPartialCompositeKeys2(stub, model.DownloadRecordUserKey, []string{userID})
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDownloadRecordListByUser-查询下载记录出错: %s", err))
	}

	var records []model.DownloadRecord
	for _, recordByte := range res {
		var record model.DownloadRecord
		err = json.Unmarshal(recordByte, &record)
		if err != nil {
			return shim.Error(fmt.Sprintf("QueryDownloadRecordListByUser-反序列化出错: %s", err))
		}
		records = append(records, record)
	}

	recordsByte, err := json.Marshal(records)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDownloadRecordListByUser-序列化出错: %s", err))
	}

	return shim.Success(recordsByte)
}

// QueryDownloadRecordListByDataset 查询下载记录列表
// args[0]: 数据集所有者ID string
// args[1]: 数据集名 string
func QueryDownloadRecordListByDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("QueryDownloadRecordListByDataset-参数数量错误")
	}

	datasetOwner := strings.ToLower(args[0])
	datasetName := strings.ToLower(args[1])
	res, err := utils.GetStateByPartialCompositeKeys2(stub, model.DownloadRecordDatasetKey, []string{datasetOwner, datasetName})
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDownloadRecordListByDataset-查询下载记录出错: %s", err))
	}

	var records []model.DownloadRecord
	for _, recordByte := range res {
		var record model.DownloadRecord
		err = json.Unmarshal(recordByte, &record)
		if err != nil {
			return shim.Error(fmt.Sprintf("QueryDownloadRecordListByDataset-反序列化出错: %s", err))
		}
		records = append(records, record)
	}

	recordsByte, err := json.Marshal(records)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDownloadRecordListByDataset-序列化出错: %s", err))
	}

	return shim.Success(recordsByte)
}
