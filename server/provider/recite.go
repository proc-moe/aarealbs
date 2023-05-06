package provider

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/e"
	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/utils/auth"
	"github.com/proc-moe/aarealbs/server/utils/klog"
)

type AddQueueReq struct {
	UserID       int `json:"user_info_id"`
	RecordInfoID int `json:"record_info_id"`
}

type AddQueueRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func AddQueue(c *gin.Context) {
	token := c.Request.Header.Get("authorization")

	var req AddQueueReq
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

	err1 := auth.UserIsAdmin(token)
	err2 := auth.UserIsSelfAndUnbanned(token, int(req.UserID))
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

	// calculate pattern
	var record model.RecordInfo
	model.DB.First(&record, req.RecordInfoID)

	var recitePattern model.RecitePattern
	model.DB.Where("p_id = ? AND round = ?", record.PatternID, 0).First(&recitePattern)

	estimate := time.Now().Unix() + int64(recitePattern.TimeGapEst)
	dbQueue := model.ReciteQueue{
		UserInfoID:     uint(req.UserID),
		RecordInfoID:   uint(req.RecordInfoID),
		RemindTimeUnix: estimate,
		Round:          0,
		RoundMax:       15,
		Status:         0,
	}

	r := model.DB.Create(&dbQueue)
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
