package service

import "github.com/proc-moe/aarealbs/server/model"

func ReciteProcessAdd(uid int) {
	var reciteProcess model.ReciteProcess
	r := model.DB.Where("user_info_id = ?", uid).First(&reciteProcess)
	if r.RowsAffected == 0 {
		reciteProcess.UserInfoID = uint(uid)
		reciteProcess.ReciteDone = 0
		reciteProcess.ReciteProcessing = 0
		reciteProcess.ReciteTODO = 1
		model.DB.Create(&reciteProcess)
	} else {
		model.DB.Model(&reciteProcess).Where("user_info_id = ?", uid).Updates(map[string]interface{}{
			"recite_todo": reciteProcess.ReciteTODO + 1,
		})
	}
}

func ReciteProcessReciteStart(uid int) {
	var reciteProcess model.ReciteProcess
	model.DB.Where("user_info_id = ?", uid).First(&reciteProcess)
	reciteProcess.ReciteTODO -= 1
	reciteProcess.ReciteProcessing += 1
	model.DB.Model(&model.ReciteProcess{}).Where("user_info_id = ?", uid).Updates(map[string]interface{}{
		"recite_todo":       reciteProcess.ReciteTODO,
		"recite_processing": reciteProcess.ReciteProcessing,
	})
}

func ReciteProcessReciteDone(uid int) {
	var reciteProcess model.ReciteProcess
	model.DB.Where("user_info_id = ?", uid).First(&reciteProcess)
	reciteProcess.ReciteProcessing -= 1
	reciteProcess.ReciteDone += 1
	model.DB.Model(&model.ReciteProcess{}).Where("user_info_id = ?", uid).Updates(map[string]interface{}{
		"recite_todo":       reciteProcess.ReciteTODO,
		"recite_processing": reciteProcess.ReciteProcessing,
	})
}
