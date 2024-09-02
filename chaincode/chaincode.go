package main

import (
	"chaincode/api"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type BlockChainGenshin struct {
}

// Init 链码初始化
func (t *BlockChainGenshin) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("链码初始化")
	// 初始化默认数据
	var userIds = [3]string{
		"user1",
		"user2",
		"admin",
	}
	var userNames = [3]string{"User1", "User2", "Admin"}
	// 初始化账号数据
	for i, val := range userIds {
		// 写入账本
		if res := api.CreateUser(stub, []string{val, userNames[i]}); res.Status != shim.OK {
			return res
		}
	}
	return shim.Success(nil)
}

// Invoke 实现Invoke接口调用智能合约
func (t *BlockChainGenshin) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case "hello":
		return api.Hello(stub, args)
	case "createUser":
		return api.CreateUser(stub, args)
	case "modifyUserName":
		return api.ModifyUserName(stub, args)
	case "queryUserList":
		return api.QueryUserList(stub, args)
	case "queryUser":
		return api.QueryUser(stub, args)
	case "createDataset":
		return api.CreateDataset(stub, args)
	case "queryDataset":
		return api.QueryDataset(stub, args)
	case "queryDatasetList":
		return api.QueryDatasetList(stub, args)
	case "appendDatasetVersion":
		return api.AppendDatasetVersion(stub, args)
	case "increaseDatasetStars":
		return api.IncreaseDatasetStars(stub, args)
	case "increaseDatasetDownloads":
		return api.IncreaseDatasetDownloads(stub, args)
	case "createFile":
		return api.CreateFile(stub, args)
	case "queryFile":
		return api.QueryFile(stub, args)
	case "queryMultipleFiles":
		return api.QueryMultipleFiles(stub, args)
	case "createDownloadRecord":
		return api.CreateDownloadRecord(stub, args)
	case "queryDownloadRecordListByUser":
		return api.QueryDownloadRecordListByUser(stub, args)
	case "queryDownloadRecordListByDataset":
		return api.QueryDownloadRecordListByDataset(stub, args)
	default:
		return shim.Error(fmt.Sprintf("没有该功能: %s", funcName))
	}
}

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = timeLocal
	err = shim.Start(new(BlockChainGenshin))
	if err != nil {
		fmt.Printf("Error starting GenShIn chaincode: %s", err)
	}
}
