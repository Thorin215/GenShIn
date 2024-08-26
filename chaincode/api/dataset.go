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

// Use composite key, [Owner, Name], all to lowercase
func getDatasetCompositeKey(owner, name string) []string {
	return []string{strings.ToLower(owner), strings.ToLower(name)}
}

func getDataset(stub shim.ChaincodeStubInterface, owner, name string) (model.Dataset, error) {
	datasetByte, err := utils.GetStateByCompositeKey(stub, model.DatasetKey, getDatasetCompositeKey(owner, name))
	if err != nil {
		return model.Dataset{}, fmt.Errorf("getDataset-查询数据集出错: %s", err)
	}
	if datasetByte == nil {
		return model.Dataset{}, fmt.Errorf("getDataset-数据集不存在")
	}
	var dataset model.Dataset
	err = json.Unmarshal(datasetByte, &dataset)
	if err != nil {
		return model.Dataset{}, fmt.Errorf("getDataset-反序列化出错: %s", err)
	}
	return dataset, nil
}
func checkDatasetExist(stub shim.ChaincodeStubInterface, owner, name string) (bool, error) {
	datasetByte, err := utils.GetStateByCompositeKey(stub, model.DatasetKey, getDatasetCompositeKey(owner, name))
	if err != nil {
		return false, fmt.Errorf("checkDatasetExist-查询数据集出错: %s", err)
	}
	return datasetByte != nil, nil
}

// CreateDataset 创建数据集
// args[0]: 所有者ID string
// args[1]: 数据集名字 string
// args[2]: 版本列表 string, []DatasetVersion as JSON
func CreateDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("CreateDataset-参数数量错误")
	}

	if exist, err := checkUserExist(stub, args[1]); err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-查询用户出错: %s", err))
	} else if !exist {
		return shim.Error("CreateDataset-参数错误: 用户不存在")
	}

	var versions []model.DatasetVersion
	if err := json.Unmarshal([]byte(args[2]), &versions); err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-反序列化出错: %s", err))
	}

	dataset := model.Dataset{
		Owner:    args[0],
		Name:     args[1],
		Versions: versions,
	}
	if err := model.ValidateDataset(dataset); err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-参数错误: %s", err))
	}

	if exist, err := checkDatasetExist(stub, dataset.Owner, dataset.Name); err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-查询数据集出错: %s", err))
	} else if exist {
		return shim.Error("CreateDataset-数据集已存在")
	}

	if err := utils.WriteLedger(dataset, stub, model.DatasetKey, []string{dataset.Owner, dataset.Name}); err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}

// QueryDataset 查询数据集信息
// args[0]: 数据集所有者ID string
// args[1]: 数据集名字 string
// return: Dataset as JSON
func QueryDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("QueryDataset-参数数量错误")
	}

	dataset, err := utils.GetStateByCompositeKey(stub, model.DatasetKey, getDatasetCompositeKey(args[0], args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDataset-查询数据集出错: %s", err))
	}

	return shim.Success(dataset)
}

// QueryDatasetList 查询某个用户的数据集列表
// args[0]: 用户ID string
// return: []Dataset as JSON
func QueryDatasetList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("QueryDatasetList-参数数量错误")
	}

	if exist, err := checkUserExist(stub, args[0]); err != nil {
		return shim.Error(fmt.Sprintf("QueryDatasetList-查询用户出错: %s", err))
	} else if !exist {
		return shim.Error("QueryDatasetList-参数错误: 用户不存在")
	}

	res, err := utils.GetStateByPartialCompositeKeys(stub, model.DatasetKey, []string{strings.ToLower(args[0])})
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDatasetList-查询数据集出错: %s", err))
	}

	var datasets []model.Dataset
	for _, datasetByte := range res {
		var dataset model.Dataset
		err = json.Unmarshal(datasetByte, &dataset)
		if err != nil {
			return shim.Error(fmt.Sprintf("QueryDatasetList-反序列化出错: %s", err))
		}
		datasets = append(datasets, dataset)
	}

	datasetsByte, err := json.Marshal(datasets)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDatasetList-序列化出错: %s", err))
	}

	return shim.Success(datasetsByte)
}
