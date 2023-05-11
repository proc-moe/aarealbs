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
)

type Record struct {
	ID          int    `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
	BatchInfoID int    `json:"batch_info_id"`
	Msg         string `json:"msg"`
	RespMsg     string `json:"resp_msg"`
	UserInfoID  int    `json:"uid"`
}
type RecordsRsp struct {
	Code   int      `json:"code"`
	Total  int      `json:"total"`
	Record []Record `json:"record"`
	Msg    string   `json:"msg"`
}

// API 2.1
func GetRecords(c *gin.Context) {
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

	var count int64
	model.DB.Model(&model.RecordInfo{}).Where("user_info_id = ?", userId).Count(&count)
	var DBrecords []model.RecordInfo
	records := make([]Record, 0)
	model.DB.Offset(int(offset)).Limit(int(limit)).Where("user_info_id = ?", userId).Find(&DBrecords)
	for _, v := range DBrecords {
		record := Record{
			ID:          int(v.ID),
			CreatedAt:   v.CreatedAt.Format(time.RFC822),
			UpdatedAt:   v.UpdatedAt.Format(time.RFC822),
			DeletedAt:   v.DeletedAt.Time.Format(time.RFC822),
			BatchInfoID: int(v.BatchInfoID),
			Msg:         v.Msg,
			RespMsg:     v.RespMsg,
			UserInfoID:  int(v.UserInfoID),
		}
		records = append(records, record)
	}

	rsp := RecordsRsp{
		Code:   0,
		Total:  int(count),
		Record: records,
		Msg:    "success",
	}
	c.JSON(200, rsp)
}

// api 2.2 修改词条
type EditRecordReq struct {
	ID        int    `json:"id"`
	Msg       string `json:"msg"`
	RespMsg   string `json:"resp_msg"`
	PatternID int    `json:"pattern_id"`
}

type EditRecordRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func EditRecord(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	var req EditRecordReq
	c.BindJSON(&req)

	record := model.RecordInfo{}
	r := model.DB.Where("id = ?", req.ID).First(&record)
	if r.RowsAffected == 0 {
		rsp := EditRecordRsp{
			Code: e.DATA_NOT_FOUND,
			Msg:  e.Str[e.DATA_NOT_FOUND] + "record id not found",
		}
		c.JSON(200, rsp)
		return
	}
	klog.I("record=%#v", record)

	err1 := auth.UserIsSelfAndUnbanned(token, int(record.UserInfoID))
	err2 := auth.UserIsAdmin(token)
	if err1 != nil && err2 != nil {
		klog.E("admin=%v, self&unbanned=%v", err1.Error(), err2.Error())
		klog.E("unauthorized")
		rsp := AddRecordRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	recordSave := model.RecordInfo{}
	recordSave.ID = record.ID
	klog.I("recordSave=%#v", recordSave)
	r = model.DB.Model(&recordSave).Updates(map[string]interface{}{
		"msg":        req.Msg,
		"resp_msg":   req.RespMsg,
		"pattern_id": req.PatternID,
	})
	klog.I("recordSaved rows=%v", r.RowsAffected)
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

// api 2.3 删除词条
type DeleteRecordReq struct {
	ID int `json:"id"`
}

type DeleteRecordRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func DelRecord(c *gin.Context) {

	token := c.Request.Header.Get("authorization")
	var req DeleteRecordReq
	c.BindJSON(&req)
	klog.I("DelRecord REQ=%v", req)
	// valid check
	record := model.RecordInfo{}
	r := model.DB.Where("id = ?", req.ID).First(&record)
	if r.RowsAffected == 0 {
		rsp := EditRecordRsp{
			Code: e.DATA_NOT_FOUND,
			Msg:  e.Str[e.DATA_NOT_FOUND] + "record id not found",
		}
		c.JSON(200, rsp)
		return
	}
	// auth check
	err1 := auth.UserIsSelfAndUnbanned(token, int(record.UserInfoID))
	err2 := auth.UserIsAdmin(token)
	if err1 != nil && err2 != nil {
		klog.E("admin=%v, self&unbanned=%v", err1.Error(), err2.Error())
		klog.E("unauthorized")
		rsp := AddRecordRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	// delete record
	record = model.RecordInfo{}
	record.ID = uint(req.ID)
	r = model.DB.Delete(&record)
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

// api 2.4 批量添加词条
type BatchAddRecordReq struct {
	Batches []BatchAddSingleRecordReq `json:"data"`
}

type BatchAddSingleRecordReq struct {
	BatchName string `json:"batch_name"`
	Msg       string `json:"msg"`
	RespMsg   string `json:"resp_msg"`
	UserID    int    `json:"uid"`
}

type BatchAddRecordRsp struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Count int    `json:"count"`
}

func BatchAddRecord(c *gin.Context) {
	token := c.Request.Header.Get("authorization")

	// param check
	var req BatchAddRecordReq
	err := c.BindJSON(&req)
	klog.I("BatchAddRecord REQ=%v", req)
	if err != nil {
		klog.E("params err")
		rsp := BatchAddRecordRsp{
			Code: e.PARAM_ERR,
			Msg:  e.Str[e.PARAM_ERR] + ":" + err.Error(),
		}
		c.JSON(200, rsp)
		return
	}

	// auth
	userID := req.Batches[0].UserID
	err1 := auth.UserIsSelfAndUnbanned(token, userID)
	err2 := auth.UserIsAdmin(token)
	if err1 != nil && err2 != nil {
		klog.E("admin=%v, self&unbanned=%v", err1.Error(), err2.Error())
		klog.E("unauthorized")
		rsp := BatchAddRecordRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	// create_batch
	batch := model.BatchInfo{
		Title:       req.Batches[0].BatchName,
		Description: req.Batches[0].BatchName,
		Author:      strconv.FormatInt(int64(req.Batches[0].UserID), 10),
	}

	r := model.DB.Create(&batch)
	// database create failed
	if r.RowsAffected == 0 || r.Error != nil {
		rsp := AddRecordRsp{
			Code: e.DATA_NOT_FOUND,
			Msg:  r.Error.Error() + fmt.Sprintf("rows affected %v", r.RowsAffected),
		}
		c.JSON(200, rsp)
		return
	}
	str := "上传失败:["
	failedCount := 0
	for k, v := range req.Batches {
		dbrecord := model.RecordInfo{
			BatchInfoID: batch.ID,
			Msg:         v.Msg,
			RespMsg:     v.RespMsg,
			UserInfoID:  uint(userID),
			PatternID:   0,
		}
		r := model.DB.Create(&dbrecord)
		if r.RowsAffected == 0 || r.Error != nil {
			klog.E("batch upd, k=%v,err=%v, record=%#v", k, err, dbrecord)
			str += " id:"
			str += strconv.FormatInt(int64(k), 10)
			failedCount++
		} else {
			service.ReciteProcessAdd(userID)
		}
	}
	if failedCount == 0 {
		rsp := BatchAddRecordRsp{
			Code: 0,
			Msg:  "success",
		}
		c.JSON(200, rsp)
		return
	}

	str += "]"
	rsp := BatchAddRecordRsp{
		Code: e.PARTIAL_FAIL,
		Msg:  e.Str[e.PARTIAL_FAIL] + str,
	}
	c.JSON(200, rsp)
}

// api 2.5 添加词条
type AddRecordReq struct {
	Msg       string `json:"msg"`
	RespMsg   string `json:"resp_msg"`
	UserID    int    `json:"uid"`
	PatternID int    `json:"pattern_id"`
}

type AddRecordRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func AddRecord(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	var req AddRecordReq
	c.BindJSON(&req)
	klog.I("AddRecord REQ=%v", req)
	err1 := auth.UserIsSelfAndUnbanned(token, req.UserID)
	err2 := auth.UserIsAdmin(token)
	if err1 != nil && err2 != nil {
		klog.E("admin=%v, self&unbanned=%v", err1.Error(), err2.Error())
		klog.E("unauthorized")
		rsp := AddRecordRsp{
			Code: e.UNAUTHORIZED,
			Msg:  e.Str[e.UNAUTHORIZED],
		}
		c.JSON(200, rsp)
		return
	}

	dbrecord := model.RecordInfo{
		BatchInfoID: 1,
		Msg:         req.Msg,
		RespMsg:     req.RespMsg,
		UserInfoID:  uint(req.UserID),
		PatternID:   uint(req.PatternID),
	}
	klog.I("/record/add %v", dbrecord)
	r := model.DB.Create(&dbrecord)
	klog.I("create record, row = %v", r.RowsAffected)
	// database create failed
	if r.RowsAffected == 0 || r.Error != nil {
		rsp := AddRecordRsp{
			Code: e.DATA_NOT_FOUND,
			Msg:  r.Error.Error() + fmt.Sprintf("rows affected %v", r.RowsAffected),
		}
		c.JSON(200, rsp)
		return
	}
	service.ReciteProcessAdd(req.UserID)
	rsp := AddRecordRsp{
		Code: 0,
		Msg:  "success",
	}
	c.JSON(200, rsp)
}
