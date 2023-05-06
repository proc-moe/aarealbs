package provider

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/e"
	"github.com/proc-moe/aarealbs/server/model"
)

type LoginReq struct {
	Token string `form:"token"`
}

type LoginRsp struct {
	Code   int
	Msg    string `json:"msg,omitempty"`
	Token  string `json:"token,omitempty"`
	Name   string `json:"name,omitempty"`
	UserId int    `json:"user_id,omitempty"`
}

// API 1.1
func Login(c *gin.Context) {
	var req LoginReq
	var rsp LoginRsp
	c.BindJSON(&req)

	token := model.Token{}
	r := model.DB.Where("token = ?", req.Token).First(&token)

	if r.RowsAffected == 0 {
		rsp.Code = e.DATA_NOT_FOUND
		rsp.Msg = e.Str[rsp.Code] + ": can't find such token in database"

	} else {
		rsp.Code = 0
		rsp.Msg = "success"
		userInfo := model.UserInfo{}
		model.DB.Where("user_id = ?", token.UserInfoID).First(&userInfo)
		rsp.UserId = int(userInfo.UserId)
		rsp.Name = userInfo.UserName
		rsp.Token = token.Token
	}

	c.JSON(200, rsp)
}

// API 1.2
func AddToken(c *gin.Context) {
	userInfo := model.UserInfo{
		UserId:   114514,
		UserName: "YJSNPITEST",
		Status:   0,
	}
	token := model.Token{
		Token:      "test",
		UserInfoID: 114514,
		Expire:     2147483647,
	}
	returnmsg := ""
	fmt.Println()
	if model.DB.Where("user_id = ?", 114514).First(&model.UserInfo{}).RowsAffected == 0 {
		r := model.DB.Create(&userInfo)
		fmt.Printf("[AddUser] ID=114514, row affected = %v\n", r.RowsAffected)
	} else {
		returnmsg += "user already exists..."
	}

	if model.DB.Where("token = ?", "test").Find(&model.Token{}).RowsAffected == 0 {
		r := model.DB.Create(&token)
		fmt.Printf("[AddToken] token=test, row affected = %v\n", r.RowsAffected)
	} else {
		returnmsg += "token already exists"
	}
	if len(returnmsg) == 0 {
		returnmsg = "added"
	}
	c.JSON(200, returnmsg)
}
