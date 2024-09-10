package v1

import (
	bc "application/blockchain"
	"application/model"
	"application/sql"
	"application/pkg/app"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QueryAllUsers(c *gin.Context) {
	appG := app.Gin{C: c}
	resp, err := bc.ChannelQuery("queryAllUsers", nil)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	var data []model.User
	if err = json.Unmarshal(resp.Payload, &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", data)
}

func QueryUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	userID := body.ID

	resp, err := bc.ChannelQuery("queryUser", [][]byte{[]byte(userID)})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	var user model.User
	if err = json.Unmarshal(resp.Payload, &user); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("反序列化出错: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", user)
}

func CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		ID   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	user := &sql.User{
		ID:		  body.ID,
		Password: body.Password,
	}

	if err := sql.CreateUser(user); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入数据库失败: %s", err.Error()))
		return
	}

	userID := body.ID
	userName := body.Name

	resp, err := bc.ChannelExecute("createUser", [][]byte{[]byte(userID), []byte(userName)})
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("调用智能合约出错: %s", err.Error()))
		return
	}

	appG.Response(http.StatusOK, "成功", resp)
}

func CheckUserLogin(c *gin.Context) {
	appG := app.Gin{C: c}
	var body struct {
		ID     string `json:"id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	userPassword, err := sql.GetUserPassword(body.ID); 
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("写入数据库失败: %s", err.Error()))
		return
	}

	if body.Password != userPassword {
		appG.Response(http.StatusUnauthorized, "失败", "密码错误")
		return
	}

	appG.Response(http.StatusOK, "成功", "登录成功")
}