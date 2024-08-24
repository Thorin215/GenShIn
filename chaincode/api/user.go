package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// QueryUserList 查询用户列表
func QueryUserList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userList []model.User
	results, err := utils.GetStateByPartialCompositeKeys(stub, model.UserKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var user model.User
			err := json.Unmarshal(v, &user)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryUserList-反序列化出错: %s", err))
			}
			userList = append(userList, user)
		}
	}
	userListByte, err := json.Marshal(userList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryUserList-序列化出错: %s", err))
	}
	return shim.Success(userListByte)
}

// CreateUser 创建用户
// args[0]: 用户名 string
func CreateUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("CreateUser-参数数量错误")
	}
	user := model.User{
		ID:   utils.GenerateID(stub, model.UserIDKey),
		Name: args[0],
	}

	err := utils.WriteLedger(user, stub, model.UserKey, []string{utils.Int2Str(user.ID)})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateUser-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}

// ModifyUser 修改用户名字
// args[0]: 用户ID int
// args[1]: 用户名 string
func ModifyUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("ModifyUser-参数数量错误")
	}
	userID := args[0]
	userName := args[1]
	userByte, err := utils.GetStateByPartialCompositeKeys2(stub, model.UserKey, []string{userID})
	if err != nil {
		return shim.Error(fmt.Sprintf("ModifyUser-查询用户出错: %s", err))
	}
	if userByte == nil {
		return shim.Error("ModifyUser-用户不存在")
	}
	var user model.User
	err = json.Unmarshal(userByte[0], &user)
	if err != nil {
		return shim.Error(fmt.Sprintf("ModifyUser-反序列化出错: %s", err))
	}
	user.Name = userName
	err = utils.WriteLedger(user, stub, model.UserKey, []string{userID})
	if err != nil {
		return shim.Error(fmt.Sprintf("ModifyUser-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}
