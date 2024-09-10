package main

import (
	"chaincode/model"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func initTest() *shim.MockStub {
	scc := new(BlockChainGenshin)
	stub := shim.NewMockStub("ex01", scc)
	res := stub.MockInit("1", [][]byte{[]byte("init")})
	if res.Status != shim.OK {
		fmt.Printf("Init failed: %s", string(res.Message))
		return nil
	}
	return stub
}

func checkInvoke(t *testing.T, stub *shim.MockStub, success bool, args [][]byte) pb.Response {
	res := stub.MockInvoke("1", args)
	if success && res.Status != shim.OK || !success && res.Status == shim.OK {
		fmt.Println("\n\n! Test failed on invoking")
		for i, arg := range args {
			fmt.Printf("! %d: %s\n", i, string(arg))
		}
		fmt.Println("! Should success: ", success)
		fmt.Println("! Status: ", res.Status)
		fmt.Println("! Message: ", res.Message)
		fmt.Println("! Payload: ", res.Payload)
		t.FailNow()
	}
	return res
}

func ToJson(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes
}

var stub *shim.MockStub

func testHelloWorld(t *testing.T) {
	fmt.Printf("\nTest: HelloWorld\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("hello"),
		}).Payload))
}

func testUser(t *testing.T) {
	fmt.Printf("\n1: CreateUser x2 [success]\n%s%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createUser"),
			[]byte("test_user1"),
			[]byte("TestUser1"),
		}).Payload),
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createUser"),
			[]byte("test_user2"),
			[]byte("TestUser2"),
		}).Payload))

	fmt.Printf("\n2: CreateUser [failed] (user already exists)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createUser"),
			[]byte("test_user1"),
			[]byte("TestUser2"),
		}).Payload))

	fmt.Printf("\n3: CreateUser [failed] (invalid user ID)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createUser"),
			[]byte("test_user####"),
			[]byte("TestUser2"),
		}).Payload))

	fmt.Printf("\n4: CreateUser [failed] (invalid username)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createUser"),
			[]byte("test_user3"),
			[]byte("Test###User2"),
		}).Payload))

	fmt.Printf("\n5: ModifyUserName [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("modifyUserName"),
			[]byte("test_user1"),
			[]byte("TestUser1_Mod"),
		}).Payload))

	fmt.Printf("\n6: ModifyUserName [failed] (user not exist)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("modifyUserName"),
			[]byte("test_user114514"),
			[]byte("TestUser2_Mod"),
		}).Payload))

	fmt.Printf("\n7: QueryAllUsers [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryAllUsers"),
		}).Payload))

	fmt.Printf("\n8: QueryUser [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryUser"),
			[]byte("test_user1"),
		}).Payload))
}

const sha256_a = "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9"
const sha256_b = "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"
const sha256_c = "4e07408562be3a2f0f6e5d51a79a8c001f4b3eac9d9ad681c7f7e1b7f268d5f0"
const sha256_d = "ef2d127de37be1e72f7c744f5a326f4f9db08e27f5b4a4273f3b1a6a6f98ae2e"
const sha256_e = "d82c8d1619ad8176d665453cfb2e55f0f7f7b3f4b8f4b7f4b7f4b7f4b7f4b7f4"
const sha256_invalid = "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4"

const dataset_owner = "test_user1"
const dataset_name = "test_dataset"
const downloader = "test_user2"

var filelist1 = []model.DatasetFile{
	{Hash: sha256_a, FileName: "aba.txt"},
	{Hash: sha256_b, FileName: "aba2.txt"},
}
var filelist2 = []model.DatasetFile{
	{Hash: sha256_a, FileName: "aba.txt"},
	{Hash: sha256_b, FileName: "aba2.txt"},
	{Hash: sha256_c, FileName: "aba3.txt"},
	{Hash: sha256_d, FileName: "aba4.txt"},
}
var filelistInvalid = []model.DatasetFile{
	{Hash: sha256_a, FileName: "aba\\//.txt"},
}

func testFile(t *testing.T) {
	fmt.Printf("\n1: CreateFile x4 [sucesss]\n%s%s%s%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte(sha256_a),
			[]byte("1024"),
		}).Payload),
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte(sha256_b),
			[]byte("1025"),
		}).Payload),
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte(sha256_c),
			[]byte("1034"),
		}).Payload),
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte(sha256_d),
			[]byte("1424"),
		}).Payload),
	)

	fmt.Printf("\n2: CreateFile [failed] (file already exists)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createFile"),
			[]byte(sha256_a),
			[]byte("1024"),
		}).Payload))

	fmt.Printf("\n3: CreateFile [failed] (invalid file size)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createFile"),
			[]byte(sha256_e),
			[]byte("-1"),
		}).Payload))

	fmt.Printf("\n4: CreateFile [failed] (invalid file hash)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createFile"),
			[]byte(sha256_invalid),
			[]byte("1024"),
		}).Payload))

	fmt.Printf("\n5: QueryFile [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryFile"),
			[]byte(sha256_a),
		}).Payload))

	fmt.Printf("\n6: QueryFiles [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryFiles"),
			[]byte(ToJson([]string{sha256_a, sha256_b})),
		}).Payload))
}

func testDataset(t *testing.T) {
	dataset_version1 := model.Version{
		Files:        filelist1,
		Rows:         100,
		CreationTime: "2021-01-01T00:00:00Z",
		ChangeLog:    "Initial version",
	}
	dataset_version2 := model.Version{
		Files:        filelist2,
		Rows:         200,
		CreationTime: "2021-01-02T00:00:00Z",
		ChangeLog:    "Add file 3, 4",
	}

	dataset_version_time_invalid := model.Version{
		Files:        filelist1,
		Rows:         100,
		CreationTime: "2021-01-01T00:00:00",
		ChangeLog:    "Initial version",
	}
	dataset_version_rows_invalid := model.Version{
		Files:        filelist1,
		Rows:         -1,
		CreationTime: "2021-01-01T00:00:00Z",
		ChangeLog:    "Initial version",
	}
	dataset_version_files_invalid := model.Version{
		Files:        filelistInvalid,
		Rows:         100,
		CreationTime: "2021-01-01T00:00:00Z",
		ChangeLog:    "Initial version",
	}

	fmt.Printf("\n1: CreateDataset [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))

	fmt.Printf("\n2: CreateDataset [failed] (dataset already exists)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))

	fmt.Printf("\n3: CreateDataset [failed] (owner not exist)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte("test_user1145"),
			[]byte(dataset_name),
		}).Payload))

	fmt.Printf("\n4: CreateDataset [failed] (invalid dataset name)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte("test#dataset"),
		}).Payload))

	fmt.Printf("\n5: AddDatasetVersion [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("addDatasetVersion"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson(dataset_version1)),
		}).Payload))

	fmt.Printf("\n6: AddDatasetVersion [failed] (invalid dataset version time)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("addDatasetVersion"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson(dataset_version_time_invalid)),
		}).Payload))

	fmt.Printf("\n7: AddDatasetVersion [failed] (invalid dataset version rows)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("addDatasetVersion"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson(dataset_version_rows_invalid)),
		}).Payload))

	fmt.Printf("\n8: AddDatasetVersion [failed] (invalid dataset version files)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("addDatasetVersion"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson(dataset_version_files_invalid)),
		}).Payload))

	fmt.Printf("\n9: AddDatasetVersion [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("addDatasetVersion"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson(dataset_version2)),
		}).Payload))

	fmt.Printf("\n10: QueryDataset [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))

	fmt.Printf("\n11: QueryAllDatasets [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryAllDatasets"),
		}).Payload))

	fmt.Printf("\n12: QueryDatasetsByUser [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryDatasetsByUser"),
			[]byte(dataset_owner),
		}).Payload))

	fmt.Printf("\n13: DeleteDataset [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("deleteDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))

	fmt.Printf("\n14: AddDatasetVersion [failed] (dataset deleted)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("addDatasetVersion"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson(dataset_version2)),
		}).Payload))
}

func testRecord(t *testing.T) {

	fmt.Printf("\n1: CreateRecord [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createRecord"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(downloader),
			[]byte(ToJson(filelist1)),
			[]byte("2021-01-01T00:00:00Z"),
		}).Payload))

	fmt.Printf("\n2: CreateRecord [failed] (dataset not exist)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createRecord"),
			[]byte(dataset_owner),
			[]byte(dataset_name + "_"),
			[]byte(downloader),
			[]byte(ToJson(filelist1)),
			[]byte("2021-01-01T00:00:00Z"),
		}).Payload))

	fmt.Printf("\n3: CreateRecord [failed] (invalid file list)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createRecord"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(downloader),
			[]byte(ToJson(filelistInvalid)),
			[]byte("2021-01-01T00:00:00Z"),
		}).Payload))

	fmt.Printf("\n4: CreateRecord [failed] (invalid download time)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createRecord"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(downloader),
			[]byte(ToJson(filelist1)),
			[]byte("2021-01-01T00:00:00"),
		}).Payload))

	fmt.Printf("\n5: QueryRecordsByUser [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryRecordsByUser"),
			[]byte(downloader),
		}).Payload))

	fmt.Printf("\n6: QueryRecordsByDataset [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryRecordsByDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))
}

func TestGenshin(t *testing.T) {
	t.Run("HelloWorld", testHelloWorld)
	t.Run("User", testUser)
	t.Run("File", testFile)
	t.Run("Dataset", testDataset)
	t.Run("Record", testRecord)
}

func TestMain(m *testing.M) {
	stub = initTest()

	if stub == nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}
