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

	/* ensure all files exist */
	for _, version := range dataset.Versions {
		for _, fileHash := range version.Files {
			if exist, err := checkFileExist(stub, fileHash); err != nil {
				return shim.Error(fmt.Sprintf("CreateDataset-查询文件出错: %s", err))
			} else if !exist {
				return shim.Error(fmt.Sprintf("CreateDataset-参数错误: 文件不存在: %s", fileHash))
			}
		}
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

// UpdateDatasetVersions 更新数据集版本
// args[0]: 所有者ID string
// args[1]: 数据集名字 string
// args[2]: 版本列表 string, []DatasetVersion as JSON
func UpdateDatasetVersions(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("UpdateDatasetVersions-参数数量错误")
	}

	if exist, err := checkDatasetExist(stub, args[0], args[1]); err != nil {
		return shim.Error(fmt.Sprintf("UpdateDatasetVersions-查询数据集出错: %s", err))
	} else if !exist {
		return shim.Error("UpdateDatasetVersions-参数错误: 数据集不存在")
	}

	dataset, err := getDataset(stub, args[0], args[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("UpdateDatasetVersions-查询数据集出错: %s", err))
	}

	var versions []model.DatasetVersion
	if err := json.Unmarshal([]byte(args[2]), &versions); err != nil {
		return shim.Error(fmt.Sprintf("UpdateDatasetVersions-反序列化出错: %s", err))
	}

	dataset.Versions = versions
	if err := model.ValidateDataset(dataset); err != nil {
		return shim.Error(fmt.Sprintf("UpdateDatasetVersions-参数错误: %s", err))
	}

	if err := utils.WriteLedger(dataset, stub, model.DatasetKey, []string{dataset.Owner, dataset.Name}); err != nil {
		return shim.Error(fmt.Sprintf("UpdateDatasetVersions-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}

// IncreaseDatasetStars 更新数据集点赞数
// args[0]: 所有者ID string
// args[1]: 数据集名字 string
// args[2]: 点赞数 int (default = 1)
func IncreaseDatasetStars(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 && len(args) != 3 {
		return shim.Error("IncreaseDatasetStars-参数数量错误")
	}

	if exist, err := checkDatasetExist(stub, args[0], args[1]); err != nil {
		return shim.Error(fmt.Sprintf("IncreaseDatasetStars-查询数据集出错: %s", err))
	} else if !exist {
		return shim.Error("IncreaseDatasetStars-参数错误: 数据集不存在")
	}

	dataset, err := getDataset(stub, args[0], args[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("IncreaseDatasetStars-查询数据集出错: %s", err))
	}

	if len(args) == 2 {
		dataset.Stars++
	} else {
		stars := utils.Str2Int32(args[2])
		if stars <= 0 {
			return shim.Error("IncreaseDatasetStars-参数错误: 必须为正整数")
		}

		dataset.Stars += stars
	}

	if err := utils.WriteLedger(dataset, stub, model.DatasetKey, []string{dataset.Owner, dataset.Name}); err != nil {
		return shim.Error(fmt.Sprintf("IncreaseDatasetStars-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}

// IncreaseDatasetDownloads 更新数据集下载数
// args[0]: 所有者ID string
// args[1]: 数据集名字 string
// args[2]: 下载数 int (default = 1)
func IncreaseDatasetDownloads(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 && len(args) != 3 {
		return shim.Error("IncreaseDatasetDownloads-参数数量错误")
	}

	if exist, err := checkDatasetExist(stub, args[0], args[1]); err != nil {
		return shim.Error(fmt.Sprintf("IncreaseDatasetDownloads-查询数据集出错: %s", err))
	} else if !exist {
		return shim.Error("IncreaseDatasetDownloads-参数错误: 数据集不存在")
	}

	dataset, err := getDataset(stub, args[0], args[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("IncreaseDatasetDownloads-查询数据集出错: %s", err))
	}

	if len(args) == 2 {
		dataset.Downloads++
	} else {
		downloads := utils.Str2Int32(args[2])
		if downloads <= 0 {
			return shim.Error("IncreaseDatasetDownloads-参数错误: 必须为正整数")
		}
		dataset.Downloads += downloads
	}

	if err := utils.WriteLedger(dataset, stub, model.DatasetKey, []string{dataset.Owner, dataset.Name}); err != nil {
		return shim.Error(fmt.Sprintf("IncreaseDatasetDownloads-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}
