package provider

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/e"
	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/utils/auth"
	"github.com/proc-moe/aarealbs/server/utils/klog"
)

type EffiencyRsp struct {
	BaseRsp
	Effiency []Effiency `json:"effiency"`
}

type Effiency struct {
	ID           int     `json:"id"`
	UID          int     `json:"uid"`
	ForgetRate   float32 `json:"forget_rate"`
	ResponseTime float32 `json:"response_time"`
	ReciteTry    int     `json:"recite_try"`
}

func GetEffiency(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	uidStr := c.Param("uid")
	var uid int64
	var err error
	if uidStr == "all" {
		uid = -1
	} else {
		uid, err = strconv.ParseInt(uidStr, 10, 64)
		if err != nil {
			klog.E("err parsing user_id,err=%v,userIDStr=%v", err, uidStr)
			rsp := RecordsRsp{
				Code: e.PARAM_ERR,
				Msg:  e.Str[e.PARAM_ERR],
			}
			c.JSON(200, rsp)
			return
		}

	}
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
	err2 := auth.UserIsSelfAndUnbanned(token, int(uid))
	if uid == -1 && err1 != nil {
		klog.E("admin=%v not admin", err1.Error())
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}
	if uid != -1 && err1 != nil && err2 != nil {
		klog.E("admin=%v, self&unbanned=%v", err1.Error(), err2.Error())
		klog.E("unauthorized")
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	var count int64
	var DBStats []model.EffiencyStats
	if uid != -1 {
		model.DB.Model(&model.EffiencyStats{}).
			Where("user_info_id = ?", uid).
			Count(&count)
		model.DB.Offset(int(offset)).Limit(int(limit)).
			Where("user_info_id = ?", uid).
			Find(&DBStats)
	} else {
		model.DB.Model(&model.ReciteQueue{}).Count(&count)
		model.DB.Offset(int(offset)).Limit(int(limit)).Find(&DBStats)
	}

	rspStats := make([]Effiency, 0)
	for _, v := range DBStats {
		klog.I("uid = %v", v.UserInfoID)
		stat := Effiency{
			ID:           int(v.ID),
			UID:          int(v.UserInfoID),
			ForgetRate:   v.ForgetRate,
			ResponseTime: v.ResponseTime,
			ReciteTry:    v.ReciteTry,
		}
		rspStats = append(rspStats, stat)
	}
	rsp := EffiencyRsp{
		BaseRsp:  SuccessBaseRsp,
		Effiency: rspStats,
	}
	c.JSON(200, rsp)
}

func GetEffiencyCount(c *gin.Context) {
	var count int64
	model.DB.Model(&model.EffiencyStats{}).Count(&count)
	rsp := MonitorCount{
		BaseRsp: SuccessBaseRsp,
		Count:   int(count),
	}
	c.JSON(200, rsp)
}
