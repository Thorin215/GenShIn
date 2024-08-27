package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getFile(stub shim.ChaincodeStubInterface, fileHash string) (model.DatasetFile, error) {
	fileByte, err := utils.GetStateByKey(stub, model.DatasetFileKey, fileHash)
	if err != nil {
		return model.DatasetFile{}, fmt.Errorf("getFile-查询文件出错: %s", err)
	}
	if fileByte == nil {
		return model.DatasetFile{}, fmt.Errorf("getFile-文件不存在")
	}
	var file model.DatasetFile
	err = json.Unmarshal(fileByte, &file)
	if err != nil {
		return model.DatasetFile{}, fmt.Errorf("getFile-反序列化出错: %s", err)
	}
	return file, nil
}
func checkFileExist(stub shim.ChaincodeStubInterface, fileHash string) (bool, error) {
	fileByte, err := utils.GetStateByKey(stub, model.DatasetFileKey, fileHash)
	if err != nil {
		return false, fmt.Errorf("checkFileExist-查询文件出错: %s", err)
	}
	return fileByte != nil, nil
}

// CreateFile 创建文件
// args[0]: 文件名 string
// args[1]: 文件大小 int64
// args[2]: 文件哈希 string
func CreateFile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("CreateFile-参数数量错误")
	}

	file := model.DatasetFile{
		FileName: args[0],
		Size:     utils.Str2Int64(args[1]),
		Hash:     args[2],
	}

	if err := model.ValidateFile(file); err != nil {
		return shim.Error(fmt.Sprintf("CreateFile-参数错误: %s", err))
	}

	if exist, err := checkFileExist(stub, file.Hash); err != nil {
		return shim.Error(fmt.Sprintf("CreateFile-查询文件出错: %s", err))
	} else if exist {
		return shim.Error("CreateFile-文件已存在")
	}

	if err := utils.WriteLedgerS(file, stub, model.DatasetFileKey, file.Hash); err != nil {
		return shim.Error(fmt.Sprintf("CreateFile-写入账本出错: %s", err))
	}

	return shim.Success(nil)
}

// QueryFile 查询文件信息
// args[0]: 文件Hash string
// return: DatasetFile as JSON
func QueryFile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("QueryFile-参数数量错误")
	}

	fileByte, err := utils.GetStateByKey(stub, model.DatasetFileKey, args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryFile-查询文件出错: %s", err))
	}
	if fileByte == nil {
		return shim.Error("QueryFile-文件不存在")
	}

	return shim.Success(fileByte)
}
