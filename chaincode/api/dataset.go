package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getDataset(stub shim.ChaincodeStubInterface, owner, name string) (model.Dataset, error) {
	datasetByte, err := utils.GetStateByKey(stub, model.DatasetKey, []string{owner, name})
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
	datasetByte, err := utils.GetStateByKey(stub, model.DatasetKey, []string{owner, name})
	if err != nil {
		return false, fmt.Errorf("checkDatasetExist-查询数据集出错: %s", err)
	}
	return datasetByte != nil, nil
}

// [CreateDataset] 创建数据集
// args[0]: 所有者ID | string
// args[1]: 数据集名字 | string
// return: nil
func CreateDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("CreateDataset-参数数量错误")
	}

	if exist, err := checkUserExist(stub, args[0]); err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-查询用户出错: %s", err))
	} else if !exist {
		return shim.Error("CreateDataset-参数错误: 用户不存在")
	}

	dataset := model.Dataset{
		Owner:    args[0],
		Name:     args[1],
		Versions: []model.Version{},
		Deleted:  false,
	}

	if err := model.ValidateDataset(dataset); err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-参数错误: %s", err))
	}

	if exist, err := checkDatasetExist(stub, dataset.Owner, dataset.Name); err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-查询数据集出错: %s", err))
	} else if exist {
		return shim.Error("CreateDataset-数据集已存在")
	}

	err := utils.WriteLedger(dataset, stub, model.DatasetKey, []string{dataset.Owner, dataset.Name})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateDataset-写入账本出错: %s", err))
	}

	return shim.Success(nil)
}

// AddDatasetVersion 添加数据集版本
// args[0]: 所有者ID | string
// args[1]: 数据集名字 | string
// args[2]: 版本 Version | string (JSON)
// return: nil
func AddDatasetVersion(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("AddDatasetVersions-参数数量错误")
	}

	if exist, err := checkDatasetExist(stub, args[0], args[1]); err != nil {
		return shim.Error(fmt.Sprintf("AddDatasetVersions-查询数据集出错: %s", err))
	} else if !exist {
		return shim.Error("AddDatasetVersions-参数错误: 数据集不存在")
	}

	dataset, err := getDataset(stub, args[0], args[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("AddDatasetVersions-查询数据集出错: %s", err))
	}

	if dataset.Deleted {
		return shim.Error("AddDatasetVersions-参数错误: 数据集已删除")
	}

	var version model.Version
	if err := json.Unmarshal([]byte(args[2]), &version); err != nil {
		return shim.Error(fmt.Sprintf("AddDatasetVersions-反序列化出错: %s", err))
	}

	dataset.Versions = append(dataset.Versions, version)
	if err := model.ValidateDataset(dataset); err != nil {
		return shim.Error(fmt.Sprintf("AddDatasetVersions-参数错误: %s", err))
	}

	fileNameMp := make(map[string]bool)

	for _, file := range version.Files {
		if exist, err := checkFileExist(stub, file.Hash); err != nil {
			return shim.Error(fmt.Sprintf("AddDatasetVersions-查询文件出错: %s", err))
		} else if !exist {
			return shim.Error(fmt.Sprintf(
				"AddDatasetVersions-参数错误: 文件不存在: %s %s",
				file.Hash,
				file.FileName,
			))
		}
		if _, ok := fileNameMp[file.FileName]; ok {
			return shim.Error(fmt.Sprintf(
				"AddDatasetVersions-参数错误: 文件名重复: %s",
				file.FileName,
			))
		}
		fileNameMp[file.FileName] = true
	}

	for _, file := range version.Files {
		if err := incrementFileReferenceCount(stub, file.Hash); err != nil {
			return shim.Error(fmt.Sprintf("AddDatasetVersions-增加引用计数出错: %s", err))
		}
	}

	err = utils.WriteLedger(dataset, stub, model.DatasetKey, []string{dataset.Owner, dataset.Name})
	if err != nil {
		return shim.Error(fmt.Sprintf("AddDatasetVersions-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}

// [QueryAllDatasets] 查询全部数据集列表
// args: nil
// return: []Dataset | string (JSON)
func QueryAllDatasets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("QueryAllDatasets-参数数量错误")
	}

	res, err := utils.GetStateByObjectType(stub, model.DatasetKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAllDatasets-查询数据集出错: %s", err))
	}

	var datasets []model.Dataset
	for _, datasetByte := range res {
		var dataset model.Dataset
		err = json.Unmarshal(datasetByte, &dataset)
		if err != nil {
			return shim.Error(fmt.Sprintf("QueryAllDatasets-反序列化出错: %s", err))
		}
		datasets = append(datasets, dataset)
	}

	datasetsByte, err := json.Marshal(datasets)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAllDatasets-序列化出错: %s", err))
	}

	return shim.Success(datasetsByte)
}

// [QueryDatasetsByUser] 查询某个用户的数据集列表
// args[0]: 用户ID | string
// return: []Dataset | string (JSON)
func QueryDatasetsByUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("QueryDatasetsByUser-参数数量错误")
	}

	if exist, err := checkUserExist(stub, args[0]); err != nil {
		return shim.Error(fmt.Sprintf("QueryDatasetsByUser-查询用户出错: %s", err))
	} else if !exist {
		return shim.Error("QueryDatasetsByUser-参数错误: 用户不存在")
	}

	res, err := utils.GetStateByPartialKey(stub, model.DatasetKey, []string{args[0]})
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDatasetsByUser-查询数据集出错: %s", err))
	}

	var datasets []model.Dataset
	for _, datasetByte := range res {
		var dataset model.Dataset
		err = json.Unmarshal(datasetByte, &dataset)
		if err != nil {
			return shim.Error(fmt.Sprintf("QueryDatasetsByUser-反序列化出错: %s", err))
		}
		datasets = append(datasets, dataset)
	}

	datasetsByte, err := json.Marshal(datasets)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDatasetsByUser-序列化出错: %s", err))
	}

	return shim.Success(datasetsByte)
}

// [QueryDataset] 查询数据集
// args[0]: 所有者ID | string
// args[1]: 数据集名字 | string
// return: Dataset | string (JSON)
func QueryDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("QueryDataset-参数数量错误")
	}

	dataset, err := utils.GetStateByKey(stub, model.DatasetKey, []string{args[0], args[1]})
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDataset-查询数据集出错: %s", err))
	}

	return shim.Success(dataset)
}

// DeleteDataset 删除数据集
// args[0]: 所有者ID | string
// args[1]: 数据集名字 | string
// return: nil
func DeleteDataset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("DeleteDataset-参数数量错误")
	}

	if exist, err := checkDatasetExist(stub, args[0], args[1]); err != nil {
		return shim.Error(fmt.Sprintf("DeleteDataset-查询数据集出错: %s", err))
	} else if !exist {
		return shim.Error("DeleteDataset-参数错误: 数据集不存在")
	}

	dataset, err := getDataset(stub, args[0], args[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("DeleteDataset-查询数据集出错: %s", err))
	}

	dataset.Deleted = true

	for _, version := range dataset.Versions {
		for _, file := range version.Files {
			if err := decrementFileReferenceCount(stub, file.Hash); err != nil {
				return shim.Error(fmt.Sprintf("DeleteDataset-减少引用计数出错: %s", err))
			}
		}
	}

	err = utils.WriteLedger(dataset, stub, model.DatasetKey, []string{dataset.Owner, dataset.Name})
	if err != nil {
		return shim.Error(fmt.Sprintf("DeleteDataset-写入账本出错: %s", err))
	}

	return shim.Success(nil)
}
