package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getUser(stub shim.ChaincodeStubInterface, userID string) (model.User, error) {
	userByte, err := utils.GetStateByKey_Single(stub, model.UserKey, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("getUser-查询用户出错: %s", err)
	}
	if userByte == nil {
		return model.User{}, fmt.Errorf("getUser-用户不存在")
	}
	var user model.User
	err = json.Unmarshal(userByte, &user)
	if err != nil {
		return model.User{}, fmt.Errorf("getUser-反序列化出错: %s", err)
	}
	return user, nil
}
func checkUserExist(stub shim.ChaincodeStubInterface, userID string) (bool, error) {
	userByte, err := utils.GetStateByKey_Single(stub, model.UserKey, userID)
	if err != nil {
		return false, fmt.Errorf("checkUserExist-查询用户出错: %s", err)
	}
	return userByte != nil, nil
}

// [QueryAllUsers] 查询用户列表
// args: nil
// return: []User | string (JSON)
func QueryAllUsers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("QueryAllUsers-参数数量错误")
	}

	res, err := utils.GetStateByObjectType(stub, model.UserKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAllUsers-查询用户出错: %s", err))
	}

	var users []model.User
	for _, userByte := range res {
		var user model.User
		err = json.Unmarshal(userByte, &user)
		if err != nil {
			return shim.Error(fmt.Sprintf("QueryAllUsers-反序列化出错: %s", err))
		}
		users = append(users, user)
	}

	usersByte, err := json.Marshal(users)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAllUsers-序列化出错: %s", err))
	}

	return shim.Success(usersByte)
}

// [QueryUser] 查询用户
// args[0]: 用户ID | string
// return: User | string (JSON)
func QueryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("QueryUser-参数数量错误")
	}
	userID := args[0]

	user, err := getUser(stub, userID)
	if err != nil {
		return shim.Error(err.Error())
	}

	userByte, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryUser-序列化出错: %s", err))
	}

	return shim.Success(userByte)
}

// [CreateUser] 创建用户
// args[0]: 用户ID | string
// args[1]: 用户名 | string
// return: nil
func CreateUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("CreateUser-参数数量错误")
	}
	userID := args[0]
	userName := args[1]

	user := model.User{
		ID:   userID,
		Name: userName,
	}

	if err := model.ValidateUser(user); err != nil {
		return shim.Error(fmt.Sprintf("CreateUser-参数错误: %s", err))
	}
	if existing, err := checkUserExist(stub, userID); err != nil {
		return shim.Error(err.Error())
	} else if existing {
		return shim.Error("CreateUser-用户已存在")
	}

	if err := utils.WriteLedger_Single(user, stub, model.UserKey, userID); err != nil {
		return shim.Error(fmt.Sprintf("CreateUser-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}

// [ModifyUserName] 修改用户名字
// args[0]: 用户ID | string
// args[1]: 新用户名 | string
// return: nil
func ModifyUserName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("ModifyUserName-参数数量错误")
	}
	userID := args[0]
	newName := args[1]

	user, err := getUser(stub, userID)
	if err != nil {
		return shim.Error(err.Error())
	}
	user.Name = newName
	err = model.ValidateUser(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("ModifyUserName-参数错误: %s", err))
	}

	err = utils.WriteLedger_Single(user, stub, model.UserKey, userID)
	if err != nil {
		return shim.Error(fmt.Sprintf("ModifyUserName-写入账本出错: %s", err))
	}
	return shim.Success(nil)
}
