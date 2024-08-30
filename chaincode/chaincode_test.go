package main

import (
	"chaincode/model"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func initTest(t *testing.T) *shim.MockStub {
	scc := new(BlockChainGenshin)
	stub := shim.NewMockStub("ex01", scc)
	checkInit(t, stub, [][]byte{[]byte("init")})
	return stub
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
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

// 测试链码初始化
func TestBlockChainGenshin_Init(t *testing.T) {
	initTest(t)
}

func Test_HelloWorld(t *testing.T) {
	stub := initTest(t)
	fmt.Printf("Test: HelloWorld\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("hello"),
		}).Payload))
}

func Test_DatasetFile(t *testing.T) {
	stub := initTest(t)

	fmt.Println("\nTest: DatasetFile")

	const sha256_a = "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9"
	const sha256_b = "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"
	const sha256_invalid = "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4"

	fmt.Printf("\n1: CreateFile [sucesss]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte("test_file"),
			[]byte("1024"),
			[]byte(sha256_a),
		}).Payload))

	fmt.Printf("\n2: CreateFile [failed] (file already exists)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createFile"),
			[]byte("test_file"),
			[]byte("1024"),
			[]byte(sha256_a),
		}).Payload))

	fmt.Printf("\n3: CreateFile [failed] (invalid filename)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createFile"),
			[]byte("t"),
			[]byte("1024"),
			[]byte(sha256_b),
		}).Payload))

	fmt.Printf("\n4: CreateFile [failed] (invalid file size)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createFile"),
			[]byte("test_file2"),
			[]byte("-1"),
			[]byte(sha256_b),
		}).Payload))

	fmt.Printf("\n5: CreateFile [failed] (invalid file hash)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createFile"),
			[]byte("test_file2"),
			[]byte("1024"),
			[]byte(sha256_invalid),
		}).Payload))

	fmt.Printf("\n6: QueryFile [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryFile"),
			[]byte(sha256_a),
		}).Payload))
}

func Test_User(t *testing.T) {
	stub := initTest(t)

	fmt.Println("\nTest: User")

	fmt.Printf("\n1: CreateUser [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createUser"),
			[]byte("test_user1"),
			[]byte("TestUser1"),
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
			[]byte("test_user2"),
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

	fmt.Printf("\n7: QueryUserList [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryUserList"),
		}).Payload))
}

func ToJson(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes
}
func Test_Dataset(t *testing.T) {
	stub := initTest(t)

	fmt.Println("\nTest: Dataset")

	fmt.Printf("\n0a: CreateUser as dataset owner [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createUser"),
			[]byte("test_user1"),
			[]byte("TestUser1"),
		}).Payload))

	const sha256_a = "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9"
	const sha256_b = "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"
	const sha256_c = "4e07408562be3a2f0f6e5d51a79a8c001f4b3eac9d9ad681c7f7e1b7f268d5f0"
	const sha256_d = "ef2d127de37be1e72f7c744f5a326f4f9db08e27f5b4a4273f3b1a6a6f98ae2e"
	const sha256_e = "d82c8d1619ad8176d665453cfb2e55f0f7f7b3f4b8f4b7f4b7f4b7f4b7f4b7f4"

	fmt.Printf("\n0b: CreateFile x4 as dataset files [success]\n%s%s%s%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte("test_file1"),
			[]byte("1024"),
			[]byte(sha256_a),
		}).Payload),
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte("test_file2"),
			[]byte("1025"),
			[]byte(sha256_b),
		}).Payload),
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte("test_file3"),
			[]byte("1034"),
			[]byte(sha256_c),
		}).Payload),
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte("test_file4"),
			[]byte("1424"),
			[]byte(sha256_d),
		}).Payload),
	)

	const dataset_owner = "test_user1"
	const dataset_name = "test_dataset"

	dataset_version1 := model.DatasetVersion{
		CreationTime: "2021-01-01T00:00:00Z",
		ChangeLog:    "Initial version",
		Rows:         100,
		Files:        []string{sha256_a, sha256_b},
	}
	dataset_version2 := model.DatasetVersion{
		CreationTime: "2021-01-02T00:00:00Z",
		ChangeLog:    "Add file 3, 4",
		Rows:         200,
		Files:        []string{sha256_a, sha256_b, sha256_c, sha256_d},
	}

	dataset_version_time_invalid := model.DatasetVersion{
		CreationTime: "2021-01-01T00:00:00",
		ChangeLog:    "Initial version",
		Rows:         100,
		Files:        []string{sha256_a, sha256_b},
	}
	dataset_version_rows_invalid := model.DatasetVersion{
		CreationTime: "2021-01-01T00:00:00Z",
		ChangeLog:    "Initial version",
		Rows:         -1,
		Files:        []string{sha256_a, sha256_b},
	}
	dataset_version_files_invalid := model.DatasetVersion{
		CreationTime: "2021-01-01T00:00:00Z",
		ChangeLog:    "Initial version",
		Rows:         100,
		Files:        []string{sha256_a, sha256_e},
	}

	res := ToJson([]model.DatasetVersion{dataset_version1})

	fmt.Printf("\n1: CreateDataset [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(res),
		}).Payload))

	fmt.Printf("\n2: CreateDataset [failed] (dataset already exists)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(res),
		}).Payload))

	fmt.Printf("\n3: CreateDataset [failed] (owner not exist)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte("test_user2"),
			[]byte(dataset_name),
			[]byte(res),
		}).Payload))

	fmt.Printf("\n4: CreateDataset [failed] (invalid dataset name)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte("test#dataset"),
			[]byte(res),
		}).Payload))

	fmt.Printf("\n5: CreateDataset [failed] (invalid dataset version time)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte("test_dataset2"),
			[]byte(ToJson([]model.DatasetVersion{dataset_version_time_invalid})),
		}).Payload))

	fmt.Printf("\n6: CreateDataset [failed] (invalid dataset version rows)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte("test_dataset3"),
			[]byte(ToJson([]model.DatasetVersion{dataset_version_rows_invalid})),
		}).Payload))

	fmt.Printf("\n7: CreateDataset [failed] (invalid dataset version files)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte("test_dataset4"),
			[]byte(ToJson([]model.DatasetVersion{dataset_version_files_invalid})),
		}).Payload))

	fmt.Printf("\n8: QueryDataset [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))

	fmt.Printf("\n9: AppendDatasetVersion [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("appendDatasetVersion"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson(dataset_version2)),
		}).Payload))

	fmt.Printf("\n10: AppendDatasetVersion [failed] (invalid dataset version files)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("appendDatasetVersion"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson(dataset_version_files_invalid)),
		}).Payload))

	fmt.Printf("\n11: IncreaseDatasetStars (default +stars) [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("increaseDatasetStars"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))

	fmt.Printf("\n12: IncreaseDatasetStars (specify +stars=3) [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("increaseDatasetStars"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte("3"),
		}).Payload))

	fmt.Printf("\n13: IncreaseDatasetStars [failed] (dataset not exist)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("increaseDatasetStars"),
			[]byte(dataset_owner),
			[]byte("test_dataset_not_exist"),
		}).Payload))

	fmt.Printf("\n14: IncreaseDatasetStars [failed] (invalid +stars)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("increaseDatasetStars"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte("-1"),
		}).Payload))

	fmt.Printf("\n15: IncreaseDatasetDownloads (default +downloads) [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("increaseDatasetDownloads"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))

	fmt.Printf("\n16: QueryDatasetList [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryDatasetList"),
			[]byte(dataset_owner),
		}).Payload))
}

func Test_DownloadRecord(t *testing.T) {
	stub := initTest(t)

	fmt.Println("\nTest: DownloadRecord")

	const dataset_owner = "test_user1"
	const dataset_name = "test_dataset"

	fmt.Printf("\n0a: CreateUser as dataset owner [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createUser"),
			[]byte("test_user1"),
			[]byte("TestUser1"),
		}).Payload))

	const sha256_a = "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9"
	const sha256_b = "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"

	fmt.Printf("\n0b: CreateFile x2 as dataset files [success]\n%s%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte("test_file1"),
			[]byte("1024"),
			[]byte(sha256_a),
		}).Payload),
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createFile"),
			[]byte("test_file2"),
			[]byte("1025"),
			[]byte(sha256_b),
		}).Payload),
	)

	dataset_version1 := model.DatasetVersion{
		CreationTime: "2021-01-01T00:00:00Z",
		ChangeLog:    "Initial version",
		Rows:         100,
		Files:        []string{sha256_a, sha256_b},
	}

	fmt.Printf("\n0c: CreateDataset [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte(ToJson([]model.DatasetVersion{dataset_version1})),
		}).Payload))

	fmt.Printf("\n0d: CreateUser as downloader [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createUser"),
			[]byte("test_user2"),
			[]byte("TestUser2"),
		}).Payload))

	fmt.Printf("\n1: CreateDownloadRecord [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("createDownloadRecord"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte("test_user2"),
			[]byte(ToJson([]string{sha256_a, sha256_b})),
			[]byte("2021-01-01T00:00:00Z"),
		}).Payload))

	fmt.Printf("\n2: CreateDownloadRecord [failed] (dataset not exist)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDownloadRecord"),
			[]byte(dataset_owner),
			[]byte("test_dataset_not_exist"),
			[]byte("test_user2"),
			[]byte(ToJson([]string{sha256_a, sha256_b})),
			[]byte("2021-01-01T00:00:00Z"),
		}).Payload))

	fmt.Printf("\n3: CreateDownloadRecord [failed] (invalid file hash)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDownloadRecord"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte("test_user2"),
			[]byte(ToJson([]string{sha256_a, "invalid_hash"})),
			[]byte("2021-01-01T00:00:00Z"),
		}).Payload))

	fmt.Printf("\n4: CreateDownloadRecord [failed] (invalid download time)\n%s",
		string(checkInvoke(t, stub, false, [][]byte{
			[]byte("createDownloadRecord"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
			[]byte("test_user2"),
			[]byte(ToJson([]string{sha256_a, sha256_b})),
			[]byte("2021-01-01T00:00:00"),
		}).Payload))

	fmt.Printf("\n5: QueryDownloadRecordListByUser [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryDownloadRecordListByUser"),
			[]byte("test_user2"),
		}).Payload))

	fmt.Printf("\n6: QueryDownloadRecordListByDataset [success]\n%s",
		string(checkInvoke(t, stub, true, [][]byte{
			[]byte("queryDownloadRecordListByDataset"),
			[]byte(dataset_owner),
			[]byte(dataset_name),
		}).Payload))
}
