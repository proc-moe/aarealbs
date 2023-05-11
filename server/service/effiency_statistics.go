package service

import (
	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/utils/klog"
)

func UpdateEffiency(uid int, isForget int, respTime int) {
	klog.I("Update Effiency")
	x := model.EffiencyStats{}
	r := model.DB.Where("user_info_id = ?", uid).Find(&x)

	if r.RowsAffected == 0 {
		x.UserInfoID = uint(uid)
		x.ForgetRate = float32(isForget)
		x.ResponseTime = float32(respTime)
		x.ReciteTry = 1
		model.DB.Create(&x)
	} else {
		x.ForgetRate = (x.ForgetRate*float32(x.ReciteTry) + float32(isForget)) / float32(x.ReciteTry+1)
		x.ResponseTime = (x.ResponseTime*float32(x.ReciteTry) + float32(respTime)) / float32(x.ReciteTry+1)
		x.ReciteTry += 1
		model.DB.Save(&x)
	}
}
