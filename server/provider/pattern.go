package provider

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/e"
	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/utils/klog"
)

type Pattern struct {
	ID      int `json:"id"`
	Round   int `json:"round"`
	EstTime int `json:"est_time"`
}
type PatternRsp struct {
	BaseRsp
	Patterns []Pattern `json:"patterns"`
}

func GetPatterns(c *gin.Context) {
	pidStr := c.Param("pattern_id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		klog.E("err parsing monitor_id,err=%v,userIDStr=%v", err, pid)
		rsp := BaseRsp{
			Code: e.PARAM_ERR,
			Msg:  e.Str[e.PARAM_ERR],
		}
		c.JSON(200, rsp)
		return
	}
	RecitePatternsDB := make([]model.RecitePattern, 0)
	model.DB.Where("pattern_id = ?", pid).Find(&RecitePatternsDB)
	patterns := make([]Pattern, 0)
	for _, v := range RecitePatternsDB {
		patterns = append(patterns, Pattern{
			ID:      int(v.ID),
			Round:   int(v.Round),
			EstTime: int(v.TimeGapEst),
		})
	}
	ret := PatternRsp{
		BaseRsp:  SuccessBaseRsp,
		Patterns: patterns,
	}
	c.JSON(200, ret)
}
