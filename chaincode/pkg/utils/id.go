package utils

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GenerateID 生成ID
func GenerateID(stub shim.ChaincodeStubInterface, idKey string) int32 {
	id, _ := stub.GetState(idKey)
	if id == nil {
		stub.PutState(idKey, []byte("1"))
		return 1
	}

	res := Byte2Int(id) + 1
	stub.PutState(idKey, Int2Byte(res))
	return res
}
