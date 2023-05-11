package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/provider"
)

func Init(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/login", provider.Login)
		api.GET("/is_admin", provider.IsAdmin)
		api.POST("/token/add", provider.AddToken)
		api.GET("/record/upload/:uid", provider.GetRecords)
		api.POST("/record/add", provider.AddRecord)
		api.POST("/record/edit", provider.EditRecord)
		api.POST("/record/delete", provider.DelRecord)
		api.POST("/record/batch/add", provider.BatchAddRecord)
		// user_info
		api.POST("/user/all", provider.GetUserInfos)
		api.POST("/user/single/:uid", provider.GetUserInfo)
		api.POST("/user/edit", provider.EditUserInfo)

		// recite
		api.POST("/queue/add", provider.AddQueue)
		api.GET("/queue/user/:uid", provider.GetUserQueue)
		api.GET("/record_history/:uid", provider.GetReciteHistory)
		api.POST("/recite/:recite_id", provider.Recite)
		api.GET("/timeup_queue/:uid", provider.GetUserTimeUpQueue)

		// // monitor
		api.GET("/monitor/single/:id", provider.GetMonitorRecord)
		api.GET("/monitor/count", provider.GetMonitorCount)

		// 6 effiency
		api.GET("/effiency/single/:uid", provider.GetEffiency)
		api.GET("/effiency/count", provider.GetEffiencyCount)

		// 7 pattern
		api.GET("/pattern/get/:pattern_id", provider.GetPatterns)
		// api.POST("/pattern/add", provider.AddPattern)
		// api.POST("/pattern/edit/:id", provider.EditPattern)
	}
}
