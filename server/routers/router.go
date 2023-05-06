package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/provider"
)

func Init(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/login", provider.Login)
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
	}
}
