// user_info
package provider

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/e"
	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/utils/auth"
	"github.com/proc-moe/aarealbs/server/utils/klog"
)

type GetAllUserInfoRsp struct {
	Code      int        `json:"code"`
	Msg       string     `json:"msg"`
	Total     int        `json:"total"`
	UserInfos []UserInfo `json:"user_infos"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Status   int    `json:"status"`
}

func GetUserInfos(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	offsetStr := c.Query("offset")
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		klog.E("err parsing offset")
		rsp := RecordsRsp{
			Code: e.PARAM_ERR,
			Msg:  e.Str[e.PARAM_ERR],
		}
		c.JSON(200, rsp)
		return
	}
	limitStr := c.Query("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		klog.E("err parsing limit")
		rsp := RecordsRsp{
			Code: e.PARAM_ERR,
			Msg:  e.Str[e.PARAM_ERR],
		}
		c.JSON(200, rsp)
		return
	}

	err1 := auth.UserIsAdmin(token)
	if err1 != nil {
		klog.E("admin=%v unauthorized", err1.Error())
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	var count int64
	model.DB.Model(&model.UserInfo{}).Count(&count)

	var dbUserInfos []model.UserInfo
	userInfos := make([]UserInfo, 0)
	model.DB.Offset(int(offset)).Limit(int(limit)).Find(&dbUserInfos)
	for _, v := range dbUserInfos {
		userInfo := UserInfo{
			ID:       int(v.ID),
			UserID:   int(v.UserId),
			UserName: v.UserName,
			Status:   int(v.Status),
		}
		userInfos = append(userInfos, userInfo)
	}
	rsp := GetAllUserInfoRsp{
		Code:      200,
		Total:     int(count),
		UserInfos: userInfos,
		Msg:       "success",
	}
	c.JSON(200, rsp)
}

type GetUserInfoRsp struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Status    int    `json:"status"`
}

func GetUserInfo(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	userIdStr := c.Param("uid")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		klog.E("err parsing user_id,err=%v,userIDStr=%v", err, userIdStr)
		rsp := RecordsRsp{
			Code: e.PARAM_ERR,
			Msg:  e.Str[e.PARAM_ERR],
		}
		c.JSON(200, rsp)
		return
	}

	err1 := auth.UserIsAdmin(token)
	err2 := auth.UserIsSelfAndUnbanned(token, int(userId))
	if err1 != nil && err2 != nil {
		klog.E("admin=%v, self&unbanned=%v", err1.Error(), err2.Error())
		klog.E("unauthorized")
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	userInfo := model.UserInfo{}

	r := model.DB.Where("user_id = ?", userId).First(&userInfo)
	if r.RowsAffected == 0 || r.Error != nil {
		rsp := EditRecordRsp{
			Code: e.DATA_NOT_FOUND,
			Msg:  r.Error.Error() + fmt.Sprintf("rows affected %v", r.RowsAffected),
		}
		c.JSON(200, rsp)
		return
	}
	rsp := GetUserInfoRsp{
		Code:      0,
		Msg:       "success",
		ID:        int(userInfo.ID),
		CreatedAt: userInfo.CreatedAt.Format(time.RFC822),
		UpdatedAt: userInfo.UpdatedAt.Format(time.RFC822),
		DeletedAt: userInfo.DeletedAt.Time.Format(time.RFC822),
		UserID:    int(userInfo.UserId),
		UserName:  userInfo.UserName,
		Status:    int(userInfo.Status),
	}
	c.JSON(200, rsp)
}

type EditUserInfoReq struct {
	UserId int `json:"user_id"`
	Status int `json:"status"`
}

type EditUserInfoRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func EditUserInfo(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	err1 := auth.UserIsAdmin(token)
	if err1 != nil {
		klog.E("admin=%v unauthorized", err1.Error())
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	var req EditUserInfoReq
	err := c.BindJSON(&req)
	if err != nil {
		klog.E("parse body failed err=%v", err.Error())
		rsp := BatchAddRecordRsp{
			Code: e.PARAM_ERR,
			Msg:  e.Str[e.PARAM_ERR] + ":" + err.Error(),
		}
		c.JSON(200, rsp)
		return
	}

	userInfoSave := model.UserInfo{}
	userInfoSave.UserId = uint(req.UserId)
	r := model.DB.Model(&userInfoSave).Where("user_id = ?", userInfoSave.UserId).
		Updates(map[string]interface{}{
			"status": req.Status,
		})
	if r.RowsAffected == 0 || r.Error != nil {
		rsp := EditRecordRsp{
			Code: e.DATA_NOT_FOUND,
			Msg:  r.Error.Error() + fmt.Sprintf("rows affected %v", r.RowsAffected),
		}
		c.JSON(200, rsp)
		return
	}
	rsp := EditRecordRsp{
		Code: 0,
		Msg:  "success",
	}
	c.JSON(200, rsp)
}
