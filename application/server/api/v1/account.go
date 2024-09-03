package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountIdBody struct {
	//AccountId string `json:"accountId"`
	ID  string `json:"id"`  // 用户ID
	Name string `json:"name"` // 用户名
}

type AccountRequestBody struct {
	Args []AccountIdBody `json:"args"`
}

func QueryAccountList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(AccountRequestBody)
	// 解析 Body 参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	var bodyBytes [][]byte
	for _, val := range body.Args {
		bodyBytes = append(bodyBytes, []byte(val.ID))
	}
	
	// 调用智能合约
	resp, err := bc.ChannelQuery("queryUserList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 反序列化 JSON
	var data []map[string]interface{}
	if err = json.Unmarshal(resp.Payload, &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", data)
}

func CheckAccount(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		Args []struct {
			ID string `json:"id"`
		} `json:"args"`
	}

	// 解析 Body 参数
	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	// 确保 args 数组不为空
	if len(body.Args) == 0 {
		appG.Response(http.StatusBadRequest, "失败", "请求体格式不正确")
		return
	}

	// 获取第一个元素的 ID
	userID := body.Args[0].ID

	// 打印提供给区块链的数据
	fmt.Printf("提供给区块链的数据: %s\n", userID)

	// 调用智能合约
	resp, err := bc.ChannelQuery("queryUser", [][]byte{[]byte(userID)})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	// 检查返回的数据
	var user map[string]interface{}
	if err = json.Unmarshal(resp.Payload, &user); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	if len(user) == 0 {
		appG.Response(http.StatusNotFound, "失败", "未找到该用户")
		return
	}

	// 将 user 包装在一个数组中
	data := []map[string]interface{}{user}

	appG.Response(http.StatusOK, "成功", data)
}


