package provider

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/e"
	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/utils/auth"
	"github.com/proc-moe/aarealbs/server/utils/klog"
)

type GetMonitorRecordRsp struct {
	BaseRsp
	Monitor `json:"monitor"`
}

type Monitor struct {
	CpuLoad     float64 `json:"cpu_load"`
	MemLoad     float64 `json:"mem_load"`
	RecordCount uint    `json:"record_count"`
	UserCount   uint    `json:"user_count"`
}

func GetMonitorRecord(c *gin.Context) {
	monitorIDStr := c.Param("id")
	monitorID, err := strconv.ParseInt(monitorIDStr, 10, 64)
	if err != nil {
		klog.E("err parsing monitor_id,err=%v,userIDStr=%v", err, monitorIDStr)
		rsp := BaseRsp{
			Code: e.PARAM_ERR,
			Msg:  e.Str[e.PARAM_ERR],
		}
		c.JSON(200, rsp)
		return
	}
	monitor := model.Monitor{}
	model.DB.Where("id = ?", monitorID).First(&monitor)
	ret := GetMonitorRecordRsp{
		BaseRsp: SuccessBaseRsp,
		Monitor: Monitor{
			CpuLoad:     monitor.CPULoad,
			MemLoad:     monitor.CPULoad,
			RecordCount: uint(monitor.RecordCount),
			UserCount:   uint(monitor.UserCount),
		},
	}
	c.JSON(200, ret)
}

type MonitorCount struct {
	BaseRsp
	Count int `json:"Count"`
}

func GetMonitorCount(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	err1 := auth.UserIsAdmin(token)
	if err1 != nil {
		klog.E("admin=%v not admin", err1.Error())
		rsp := BaseRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}
	var count int64
	model.DB.Model(&model.Monitor{}).Count(&count)
	rsp := MonitorCount{
		BaseRsp: SuccessBaseRsp,
		Count:   int(count),
	}
	c.JSON(200, rsp)
}
