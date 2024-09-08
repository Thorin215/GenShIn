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

		// user api
	case "queryAllUsers":
		return api.QueryAllUsers(stub, args)
	case "queryUser":
		return api.QueryUser(stub, args)
	case "createUser":
		return api.CreateUser(stub, args)
	case "modifyUserName":
		return api.ModifyUserName(stub, args)

		// file api
	case "createFile":
		return api.CreateFile(stub, args)
	case "queryFile":
		return api.QueryFile(stub, args)
	case "queryFiles":
		return api.QueryFiles(stub, args)

		// dataset api
	case "createDataset":
		return api.CreateDataset(stub, args)
	case "addDatasetVersion":
		return api.AddDatasetVersion(stub, args)
	case "queryAllDatasets":
		return api.QueryAllDatasets(stub, args)
	case "queryDatasetsByUser":
		return api.QueryDatasetsByUser(stub, args)
	case "queryDataset":
		return api.QueryDataset(stub, args)
	case "incrementDatasetStars":
		return api.IncrementDatasetStars(stub, args)
	case "incrementDatasetDownloads":
		return api.IncrementDatasetDownloads(stub, args)
	case "deleteDataset":
		return api.DeleteDataset(stub, args)

		// record api
	case "createRecord":
		return api.CreateRecord(stub, args)
	case "queryRecordsByUser":
		return api.QueryRecordsByUser(stub, args)
	case "queryRecordsByDataset":
		return api.QueryRecordsByDataset(stub, args)

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
