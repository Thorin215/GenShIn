package utils

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// WriteLedger 写入账本，复合主键
func WriteLedger(obj interface{}, stub shim.ChaincodeStubInterface, objectType string, keys []string) error {
	// 创建复合主键
	var key string
	if val, err := stub.CreateCompositeKey(objectType, keys); err != nil {
		return fmt.Errorf("%s-创建复合主键出错: %s", objectType, err)
	} else {
		key = val
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("%s-序列化json数据失败出错: %s", objectType, err)
	}
	// 写入区块链账本
	if err := stub.PutState(key, bytes); err != nil {
		return fmt.Errorf("%s-写入区块链账本出错: %s", objectType, err)
	}
	return nil
}

// DelLedger 删除账本，复合主键
func DelLedger(stub shim.ChaincodeStubInterface, objectType string, keys []string) error {
	// 创建复合主键
	var key string
	if val, err := stub.CreateCompositeKey(objectType, keys); err != nil {
		return fmt.Errorf("%s-创建复合主键出错: %s", objectType, err)
	} else {
		key = val
	}
	// 写入区块链账本
	if err := stub.DelState(key); err != nil {
		return fmt.Errorf("%s-删除区块链账本出错: %s", objectType, err)
	}
	return nil
}

// WriteLedger_Single 写入账本，单主键
func WriteLedger_Single(obj interface{}, stub shim.ChaincodeStubInterface, objectType string, key string) error {
	return WriteLedger(obj, stub, objectType, []string{key})
}

// DelLedger_Single 删除账本，单主键
func DelLedger_Single(stub shim.ChaincodeStubInterface, objectType string, key string) error {
	return DelLedger(stub, objectType, []string{key})
}

// GetStateByMultiplePartialKeys 根据复合主键查询数据(适合获取全部，多个，单个数据)
// 将 keys 拆分查询
func GetStateByMultiplePartialKeys(stub shim.ChaincodeStubInterface, objectType string, keys []string) (results [][]byte, err error) {
	if len(keys) == 0 {
		// 传入的keys长度为0，则查找并返回所有数据
		// 通过主键从区块链查找相关的数据，相当于对主键的模糊查询
		resultIterator, err := stub.GetStateByPartialCompositeKey(objectType, keys)
		if err != nil {
			return nil, fmt.Errorf("%s-获取全部数据出错: %s", objectType, err)
		}
		defer resultIterator.Close()

		//检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
		for resultIterator.HasNext() {
			val, err := resultIterator.Next()
			if err != nil {
				return nil, fmt.Errorf("%s-返回的数据出错: %s", objectType, err)
			}

			results = append(results, val.GetValue())
		}
	} else {
		// 传入的keys长度不为0，查找相应的数据并返回
		for _, v := range keys {
			// 创建组合键
			key, err := stub.CreateCompositeKey(objectType, []string{v})
			if err != nil {
				return nil, fmt.Errorf("%s-创建组合键出错: %s", objectType, err)
			}
			// 从账本中获取数据
			bytes, err := stub.GetState(key)
			if err != nil {
				return nil, fmt.Errorf("%s-获取数据出错: %s", objectType, err)
			}

			if bytes != nil {
				results = append(results, bytes)
			}
		}
	}

	return results, nil
}

// GetStateByPartialeKey 根据复合主键查询数据
// 将 keys 拼接查询
func GetStateByPartialKey(stub shim.ChaincodeStubInterface, objectType string, keys []string) (results [][]byte, err error) {
	// 通过主键从区块链查找相关的数据，相当于对主键的模糊查询
	resultIterator, err := stub.GetStateByPartialCompositeKey(objectType, keys)
	if err != nil {
		return nil, fmt.Errorf("%s-获取全部数据出错: %s", objectType, err)
	}
	defer resultIterator.Close()

	//检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
	for resultIterator.HasNext() {
		val, err := resultIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("%s-返回的数据出错: %s", objectType, err)
		}

		results = append(results, val.GetValue())
	}
	return results, nil
}

// GetStateByKey 根据复合主键查询数据
func GetStateByKey(stub shim.ChaincodeStubInterface, objectType string, keys []string) ([]byte, error) {
	key, err := stub.CreateCompositeKey(objectType, keys)
	if err != nil {
		return nil, fmt.Errorf("%s-创建组合键出错: %s", objectType, err)
	}
	bytes, err := stub.GetState(key)
	if err != nil {
		return nil, fmt.Errorf("%s-获取数据出错: %s", objectType, err)
	}
	return bytes, nil
}

// GetStateByKey_Single 根据单键查询数据
func GetStateByKey_Single(stub shim.ChaincodeStubInterface, objectType string, key string) ([]byte, error) {
	return GetStateByKey(stub, objectType, []string{key})
}

// GetStateByObjectType 根据对象类型查询所有数据
func GetStateByObjectType(stub shim.ChaincodeStubInterface, objectType string) (results [][]byte, err error) {
	return GetStateByPartialKey(stub, objectType, []string{})
}
