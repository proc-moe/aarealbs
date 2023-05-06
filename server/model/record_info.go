package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type RecordInfo struct {
	gorm.Model
	BatchInfoID uint // 条目集
	BatchInfo   BatchInfo
	Msg         string
	RespMsg     string
	UserInfoID  uint     // 作者ID
	UserInfo    UserInfo `gorm:"references:UserId"`
	PatternID   uint
}
