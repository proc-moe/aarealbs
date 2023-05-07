package provider

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/e"
	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/service"
	"github.com/proc-moe/aarealbs/server/utils/auth"
	"github.com/proc-moe/aarealbs/server/utils/klog"
	"gorm.io/gorm"
)

type AddQueueReq struct {
	UserID       int `json:"uid"`
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

	service.ReciteProcessReciteStart(req.UserID)
	rsp := EditRecordRsp{
		Code: 0,
		Msg:  "success",
	}
	c.JSON(200, rsp)
}

type GetUserQueueRsp struct {
	Code   int         `json:"code"`
	Total  int         `json:"total"`
	Record []UserQueue `json:"record"`
	Msg    string      `json:"msg"`
}

type UserQueue struct {
	ID             int    `json:"id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	DeletedAt      string `json:"deleted_at"`
	UserId         int    `json:"uid"`
	RecordInfoID   int    `json:"record_info_id"`
	RemindTimeUnix int    `json:"remind_time_unix"`
	Round          int    `json:"round"`
	RoundMax       int    `json:"round_max"`
	Status         int    `json:"status"`
}

func GetUserQueue(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	userIdStr := c.Param("uid")
	var userId int64
	var err error
	if userIdStr == "all" {
		userId = -1
	} else {
		userId, err = strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			klog.E("err parsing user_id,err=%v,userIDStr=%v", err, userIdStr)
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
	err2 := auth.UserIsSelfAndUnbanned(token, int(userId))
	if userId == -1 && err1 != nil {
		klog.E("admin=%v not admin", err1.Error())
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}
	if userId != -1 && err1 != nil && err2 != nil {
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
	var DBQueue []model.ReciteQueue
	if userId != -1 {
		model.DB.Model(&model.ReciteQueue{}).Where("user_info_id = ?", userId).Count(&count)
		model.DB.Offset(int(offset)).Limit(int(limit)).Where("user_info_id = ?", userId).Find(&DBQueue)
	} else {
		model.DB.Model(&model.ReciteQueue{}).Count(&count)
		model.DB.Offset(int(offset)).Limit(int(limit)).Find(&DBQueue)
	}

	rspQueue := make([]UserQueue, 0)
	for _, v := range DBQueue {
		queue := UserQueue{
			ID:             int(v.ID),
			CreatedAt:      v.CreatedAt.Format(time.RFC822),
			UpdatedAt:      v.UpdatedAt.Format(time.RFC822),
			DeletedAt:      v.DeletedAt.Time.Format(time.RFC822),
			UserId:         int(v.UserInfoID),
			RecordInfoID:   int(v.RecordInfoID),
			RemindTimeUnix: int(v.RemindTimeUnix),
			Round:          int(v.Round),
			RoundMax:       int(v.RoundMax),
			Status:         v.Status,
		}
		rspQueue = append(rspQueue, queue)
	}
	rsp := GetUserQueueRsp{
		Code:   0,
		Total:  int(count),
		Record: rspQueue,
		Msg:    "success",
	}
	c.JSON(200, rsp)
}

// api 4.3
type ReciteHistoryRsp struct {
	Code    int             `json:"code"`
	Total   int             `json:"total"`
	History []ReciteHistory `json:"record"`
	Msg     string          `json:"msg"`
}
type ReciteHistory struct {
	ID           int    `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
	UserInfoId   int    `json:"uid"`
	RecordInfoId int    `json:"record_info_id"`
	TimeGap      int    `json:"time_gap"`
	TimeGapEst   int    `json:"time_gap_est"`
	Result       int    `json:"result"`
}

func GetReciteHistory(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	userIdStr := c.Param("uid")
	var userId int64
	var err error
	if userIdStr == "all" {
		userId = -1
	} else {
		userId, err = strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			klog.E("err parsing user_id,err=%v,userIDStr=%v", err, userIdStr)
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
	err2 := auth.UserIsSelfAndUnbanned(token, int(userId))
	if userId == -1 && err1 != nil {
		klog.E("admin=%v not admin", err1.Error())
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}
	if userId != -1 && err1 != nil && err2 != nil {
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
	var DBHistory []model.ReciteHistory
	if userId != -1 {
		model.DB.Model(&model.ReciteHistory{}).Where("user_info_id = ?", userId).Count(&count)
		model.DB.Offset(int(offset)).Limit(int(limit)).Where("user_info_id = ?", userId).Find(&DBHistory)
	} else {
		model.DB.Model(&model.ReciteHistory{}).Count(&count)
		model.DB.Offset(int(offset)).Limit(int(limit)).Find(&DBHistory)
	}

	rspHistory := make([]ReciteHistory, 0)
	for _, v := range DBHistory {
		queue := ReciteHistory{
			ID:           int(v.ID),
			CreatedAt:    v.CreatedAt.Format(time.RFC822),
			UpdatedAt:    v.UpdatedAt.Format(time.RFC822),
			DeletedAt:    v.DeletedAt.Time.Format(time.RFC822),
			UserInfoId:   int(v.UserInfoID),
			RecordInfoId: int(v.RecordInfoID),
			TimeGap:      int(v.TimeGap),
			TimeGapEst:   int(v.TimeGapEst),
			Result:       int(v.Result),
		}
		rspHistory = append(rspHistory, queue)
	}
	rsp := ReciteHistoryRsp{
		Code:    0,
		Total:   int(count),
		History: rspHistory,
		Msg:     "success",
	}
	c.JSON(200, rsp)
}

type ReciteReq struct {
	ReciteID int `json:"recite_id"`
	Result   int `json:"result"`
}

type ReciteRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Recite(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	var req ReciteReq
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

	reciteIdStr := c.Param("recite_id")
	reciteId, err := strconv.ParseInt(reciteIdStr, 10, 64)
	if err != nil {
		klog.E("err parsing reciteId:%v", reciteIdStr)
		rsp := RecordsRsp{
			Code: e.PARAM_ERR,
			Msg:  e.Str[e.PARAM_ERR],
		}
		c.JSON(200, rsp)
		return
	}
	klog.I("[recite.go/recite] reciteID:%v", reciteId)
	// -=================parse DONE======================

	reciteQueue := &model.ReciteQueue{}
	r := model.DB.Where("id = ?", reciteId).First(&reciteQueue)
	if r.RowsAffected == 0 || r.Error != nil {
		rsp := EditRecordRsp{
			Code: e.DATA_NOT_FOUND,
			Msg:  r.Error.Error() + fmt.Sprintf("rows affected %v", r.RowsAffected),
		}
		c.JSON(200, rsp)
		return
	}

	klog.I("[recite.go/recite] reciteQueue userID:%v, recordId:%v,round:%v", reciteQueue.UserInfoID, reciteQueue.RecordInfoID, reciteQueue.RecordInfo)
	err1 := auth.UserIsAdmin(token)
	err2 := auth.UserIsSelfAndUnbanned(token, int(reciteQueue.UserInfoID))
	if err1 != nil && err2 != nil {
		klog.E("admin=%v, self&unbanned=%v", err1.Error(), err2.Error())
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	// 计算下一次
	recordInfo := &model.RecordInfo{}
	r = model.DB.Where("id = ?", reciteQueue.RecordInfoID).First(&recordInfo)
	if r.RowsAffected == 0 || r.Error != nil {
		rsp := EditRecordRsp{
			Code: e.DATA_NOT_FOUND,
			Msg:  r.Error.Error() + fmt.Sprintf("rows affected %v", r.RowsAffected),
		}
		c.JSON(200, rsp)
		return
	}
	klog.I("[recite.go/recite] recordID:%v msg:%v", recordInfo.ID, recordInfo.Msg)

	oldPattern := &model.RecitePattern{}
	r = model.DB.Where("p_id = ? AND round = ?", recordInfo.PatternID, reciteQueue.Round).First(&oldPattern)

	if req.Result == 0 {
		reciteQueue.Round += 1
	}
	recitePattern := &model.RecitePattern{}
	r = model.DB.Where("p_id = ? AND round = ?", recordInfo.PatternID, reciteQueue.Round).First(&recitePattern)

	// 背完了
	klog.I("now: %v", time.Now().Unix())
	klog.I("queue remind time: %v", reciteQueue.RemindTimeUnix)
	if r.Error != nil && r.Error.Error() == "record not found" || reciteQueue.Round >= reciteQueue.RoundMax {
		model.DB.Transaction(func(tx *gorm.DB) error {
			queue := model.ReciteQueue{}
			queue.ID = uint(reciteId)
			r := model.DB.Unscoped().Delete(&queue)
			klog.I("delete old queue, row=%v", r.RowsAffected)
			r = model.DB.Create(&model.ReciteHistory{
				UserInfoID:   reciteQueue.UserInfoID,
				RecordInfoID: reciteQueue.RecordInfoID,
				TimeGap:      uint(time.Now().Unix() - reciteQueue.CreatedAt.Unix()),
				TimeGapEst:   oldPattern.TimeGapEst,
				Result:       uint(req.Result),
			})
			klog.I("create history, row=%v", r.RowsAffected)
			return nil
		})
		service.UpdateEffiency(int(reciteQueue.UserInfoID), req.Result, int(time.Now().Unix())-int(reciteQueue.CreatedAt.Unix())-int(oldPattern.TimeGapEst))
		service.ReciteProcessReciteDone(int(reciteQueue.UserInfoID))
		rsp := ReciteRsp{
			Code: 200,
			Msg:  fmt.Sprintf("rows affected %v", r.RowsAffected),
		}
		c.JSON(200, rsp)
		return
	}
	// 否则就报错
	if r.RowsAffected == 0 || r.Error != nil {
		rsp := ReciteRsp{
			Code: e.DB_ERR,
			Msg:  e.Str[e.DB_ERR] + fmt.Sprintf(",err:%v", r.Error.Error()),
		}
		c.JSON(200, rsp)
		return
	}
	klog.I("[recite.go/recite] pid:%v round:%v gap:%v", recitePattern.PID, recitePattern.Round, recitePattern.TimeGapEst)

	var timeEst int = int(recitePattern.TimeGapEst) + int(time.Now().Unix())
	klog.I("[recite.go/recite] nxtTime:%v", timeEst)

	str := "new queueID:"
	// 推进新一轮
	model.DB.Transaction(func(tx *gorm.DB) error {
		queue := model.ReciteQueue{}
		queue.ID = uint(reciteId)
		r := model.DB.Unscoped().Delete(&queue)
		klog.I("delete old queue, row=%v", r.RowsAffected)
		r = model.DB.Create(&model.ReciteHistory{
			UserInfoID:   reciteQueue.UserInfoID,
			RecordInfoID: reciteQueue.RecordInfoID,
			TimeGap:      uint(time.Now().Unix() - reciteQueue.CreatedAt.Unix()),
			TimeGapEst:   recitePattern.TimeGapEst,
			Result:       uint(req.Result),
		})
		klog.I("create history, row=%v", r.RowsAffected)
		ret := &model.ReciteQueue{
			UserInfoID:     reciteQueue.UserInfoID,
			RecordInfoID:   reciteQueue.RecordInfoID,
			RemindTimeUnix: int64(timeEst),
			Round:          reciteQueue.Round,
			RoundMax:       reciteQueue.RoundMax,
			Status:         reciteQueue.Status,
		}
		r = model.DB.Create(ret)
		klog.I("create new queue, row=%v", r.RowsAffected)
		str += fmt.Sprintf("%v", ret.ID)
		return nil
	})
	c.JSON(200, ReciteRsp{
		Code: 0,
		Msg:  "success," + str,
	})

}

func GetUserTimeUpQueue(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	userIdStr := c.Param("uid")
	var userId int64
	var err error
	if userIdStr == "all" {
		userId = -1
	} else {
		userId, err = strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			klog.E("err parsing user_id,err=%v,userIDStr=%v", err, userIdStr)
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
	err2 := auth.UserIsSelfAndUnbanned(token, int(userId))
	if userId == -1 && err1 != nil {
		klog.E("admin=%v not admin", err1.Error())
		rsp := RecordsRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}
	if userId != -1 && err1 != nil && err2 != nil {
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
	var DBQueue []model.ReciteQueue
	now := time.Now().Unix()
	if userId != -1 {
		model.DB.Model(&model.ReciteQueue{}).
			Where("user_info_id = ? AND remind_time_unix < ?", userId, now).
			Count(&count)
		model.DB.Offset(int(offset)).Limit(int(limit)).
			Where("user_info_id = ? AND remind_time_unix < ?", userId, now).
			Find(&DBQueue)
	} else {
		model.DB.Model(&model.ReciteQueue{}).Count(&count)
		model.DB.Offset(int(offset)).Limit(int(limit)).Find(&DBQueue)
	}

	rspQueue := make([]UserQueue, 0)
	for _, v := range DBQueue {
		queue := UserQueue{
			ID:             int(v.ID),
			CreatedAt:      v.CreatedAt.Format(time.RFC822),
			UpdatedAt:      v.UpdatedAt.Format(time.RFC822),
			DeletedAt:      v.DeletedAt.Time.Format(time.RFC822),
			UserId:         int(v.UserInfoID),
			RecordInfoID:   int(v.RecordInfoID),
			RemindTimeUnix: int(v.RemindTimeUnix),
			Round:          int(v.Round),
			RoundMax:       int(v.RoundMax),
			Status:         v.Status,
		}
		rspQueue = append(rspQueue, queue)
	}
	rsp := GetUserQueueRsp{
		Code:   0,
		Total:  int(count),
		Record: rspQueue,
		Msg:    "success",
	}
	c.JSON(200, rsp)
}
