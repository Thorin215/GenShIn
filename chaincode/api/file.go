package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getFile(stub shim.ChaincodeStubInterface, fileHash string) (model.File, error) {
	fileByte, err := utils.GetStateByKey_Single(stub, model.FileKey, fileHash)
	if err != nil {
		return model.File{}, fmt.Errorf("getFile-查询文件出错: %s", err)
	}
	if fileByte == nil {
		return model.File{}, fmt.Errorf("getFile-文件不存在: %s", fileHash)
	}
	var file model.File
	err = json.Unmarshal(fileByte, &file)
	if err != nil {
		return model.File{}, fmt.Errorf("getFile-反序列化出错: %s", err)
	}
	return file, nil
}
func checkFileExist(stub shim.ChaincodeStubInterface, fileHash string) (bool, error) {
	fileByte, err := utils.GetStateByKey_Single(stub, model.FileKey, fileHash)
	if err != nil {
		return false, fmt.Errorf("checkFileExist-查询文件出错: %s", err)
	}
	return fileByte != nil, nil
}
func incrementFileReferenceCount(stub shim.ChaincodeStubInterface, fileHash string) error {
	file, err := getFile(stub, fileHash)
	if err != nil {
		return err
	}
	file.ReferenceCount++
	return utils.WriteLedger_Single(file, stub, model.FileKey, file.Hash)
}
func decrementFileReferenceCount(stub shim.ChaincodeStubInterface, fileHash string) error {
	file, err := getFile(stub, fileHash)
	if err != nil {
		return err
	}
	if file.ReferenceCount <= 0 {
		return fmt.Errorf("decrementFileReferenceCount-引用计数小于零: %s", fileHash)
	}
	file.ReferenceCount--
	return utils.WriteLedger_Single(file, stub, model.FileKey, file.Hash)
}

// [CreateFile] 创建文件
// args[0]: 文件哈希 | string
// args[1]: 文件大小 | string (int64)
// return: nil
func CreateFile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("CreateFile-参数数量错误")
	}

	file := model.File{
		Hash:           args[0],
		Size:           utils.Str2Int64(args[1]),
		ReferenceCount: 0,
	}

	if err := model.ValidateFile(file); err != nil {
		return shim.Error(fmt.Sprintf("CreateFile-参数错误: %s", err))
	}

	if exist, err := checkFileExist(stub, file.Hash); err != nil {
		return shim.Error(fmt.Sprintf("CreateFile-查询文件出错: %s", err))
	} else if exist {
		return shim.Error("CreateFile-文件已存在")
	}

	if err := utils.WriteLedger_Single(file, stub, model.FileKey, file.Hash); err != nil {
		return shim.Error(fmt.Sprintf("CreateFile-写入账本出错: %s", err))
	}

	return shim.Success(nil)
}

// [QueryFile] 查询文件信息
// args[0]: 文件哈希 | string (SHA-256)
// return: File | string (JSON)
func QueryFile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("QueryFile-参数数量错误")
	}

	fileByte, err := utils.GetStateByKey_Single(stub, model.FileKey, args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryFile-查询文件出错: %s", err))
	}
	if fileByte == nil {
		return shim.Error("QueryFile-文件不存在")
	}

	return shim.Success(fileByte)
}

// [QueryFiles] 查询多个文件信息
// args[0]: 文件Hash []string | string (JSON)
// return: []File | string (JSON)
func QueryFiles(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("QueryFiles-参数数量错误")
	}

	var fileHashes []string
	err := json.Unmarshal([]byte(args[0]), &fileHashes)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryFiles-反序列化出错: %s", err))
	}

	var files []model.File
	for _, fileHash := range fileHashes {
		file, err := getFile(stub, fileHash)
		if err != nil {
			return shim.Error(err.Error())
		}
		files = append(files, file)
	}

	filesByte, err := json.Marshal(files)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryFiles-序列化出错: %s", err))
	}

	return shim.Success(filesByte)
}
