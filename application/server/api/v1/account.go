// package v1

// import (
// 	bc "application/blockchain"
// 	"application/pkg/app"
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type AccountIdBody struct {
// 	AccountId string `json:"accountId"`
// }

// type AccountRequestBody struct {
// 	Args []AccountIdBody `json:"args"`
// }

// func QueryAccountList(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	body := new(AccountRequestBody)
// 	//解析Body参数
// 	if err := c.ShouldBind(body); err != nil {
// 		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
// 		return
// 	}
// 	var bodyBytes [][]byte
// 	for _, val := range body.Args {
// 		bodyBytes = append(bodyBytes, []byte(val.AccountId))
// 	}
// 	//调用智能合约
// 	resp, err := bc.ChannelQuery("queryAccountList", bodyBytes)
// 	if err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", err.Error())
// 		return
// 	}
// 	// 反序列化json
// 	var data []map[string]interface{}
// 	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
// 		appG.Response(http.StatusInternalServerError, "失败", err.Error())
// 		return
// 	}
// 	appG.Response(http.StatusOK, "成功", data)
// }

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
	//resp, err := bc.ChannelQuery("queryAccountList", bodyBytes)
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

// package api

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
// 	"github.com/your_project/app"
// )

func CheckAccount(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		ID string `json:"id"`
	}

	// 解析 Body 参数
	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	// 调用智能合约
	resp, err := bc.ChannelQuery("queryUser", [][]byte{[]byte(body.ID)})
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

	appG.Response(http.StatusOK, "成功", user)
}


