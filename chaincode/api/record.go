package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// [CreateRecord] 创建下载记录
// args[0]: 所有者ID | string
// args[1]: 数据集名 | string
// args[2]: 下载者ID | string
// args[3]: 文件列表 []DatasetFile | string (JSON)
// args[4]: 下载时间 | string
// return: nil
func CreateRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("CreateRecord-参数数量错误")
	}

	var files []model.DatasetFile
	if err := json.Unmarshal([]byte(args[3]), &files); err != nil {
		return shim.Error(fmt.Sprintf("CreateRecord-反序列化出错: %s", err))
	}

	record := model.Record{
		DatasetOwner: args[0],
		DatasetName:  args[1],
		User:         args[2],
		Files:        files,
		Time:         args[4],
	}

	if err := model.ValidateRecord(record); err != nil {
		return shim.Error(fmt.Sprintf("CreateRecord-参数错误: %s", err))
	}

	if exist, err := checkUserExist(stub, record.User); err != nil {
		return shim.Error(fmt.Sprintf("CreateRecord-查询用户出错: %s", err))
	} else if !exist {
		return shim.Error(fmt.Sprintf("CreateRecord-参数错误: 用户不存在: %s", record.User))
	}

	if exist, err := checkDatasetExist(stub, record.DatasetOwner, record.DatasetName); err != nil {
		return shim.Error(fmt.Sprintf("CreateRecord-查询数据集出错: %s", err))
	} else if !exist {
		return shim.Error(fmt.Sprintf("CreateRecord-参数错误: 数据集不存在: %s/%s",
			record.DatasetOwner,
			record.DatasetName,
		))
	}

	for _, file := range record.Files {
		if exist, err := checkFileExist(stub, file.Hash); err != nil {
			return shim.Error(fmt.Sprintf("CreateRecord-查询文件出错: %s", err))
		} else if !exist {
			return shim.Error(fmt.Sprintf("CreateRecord-参数错误: 文件不存在: %s %s",
				file.Hash,
				file.FileName,
			))
		}
	}

	keyDataset := []string{
		record.DatasetOwner,
		record.DatasetName,
		record.User,
	}
	keyUser := []string{
		record.User,
		record.DatasetOwner,
		record.DatasetName,
	}

	if err := utils.WriteLedger(record, stub, model.RecordDatasetKey, keyDataset); err != nil {
		return shim.Error(fmt.Sprintf("CreateRecord-写入账本出错: %s", err))
	}
	if err := utils.WriteLedger(record, stub, model.RecordUserKey, keyUser); err != nil {
		return shim.Error(fmt.Sprintf("CreateRecord-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}

// [QueryRecordsByUser] 查询下载记录列表
// args[0]: 下载者ID | string
// return: []Record | string (JSON)
func QueryRecordsByUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("QueryRecordsByUser-参数数量错误")
	}

	res, err := utils.GetStateByPartialKey(stub, model.RecordUserKey, []string{args[0]})
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryRecordsByUser-查询记录出错: %s", err))
	}

	var records []model.Record
	for _, recordByte := range res {
		var record model.Record
		err = json.Unmarshal(recordByte, &record)
		if err != nil {
			return shim.Error(fmt.Sprintf("QueryRecordsByUser-反序列化出错: %s", err))
		}
		records = append(records, record)
	}

	recordsByte, err := json.Marshal(records)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryRecordsByUser-序列化出错: %s", err))
	}

	return shim.Success(recordsByte)
}

// [QueryRecordsByDataset] 查询下载记录列表
// args[0]: 所有者ID | string
// args[1]: 数据集名字 | string
func QueryRecordsByDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("QueryRecordsByDataset-参数数量错误")
	}

	res, err := utils.GetStateByPartialKey(stub, model.RecordDatasetKey, []string{args[0], args[1]})
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryRecordsByDataset-查询记录出错: %s", err))
	}

	var records []model.Record
	for _, recordByte := range res {
		var record model.Record
		err = json.Unmarshal(recordByte, &record)
		if err != nil {
			return shim.Error(fmt.Sprintf("QueryRecordsByDataset-反序列化出错: %s", err))
		}
		records = append(records, record)
	}

	recordsByte, err := json.Marshal(records)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryRecordsByDataset-序列化出错: %s", err))
	}

	return shim.Success(recordsByte)
}
