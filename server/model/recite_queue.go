package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type ReciteInfoStatus int

const (
	ReciteRecordStatus_On      = 0
	ReciteRecordStatus_Off     = 1
	ReciteRecordStatus_Deleted = 2
)

type ReciteQueue struct {
	gorm.Model
	UserInfoID     uint       `gorm:"uniqueIndex:recite_queue_idx"`
	UserInfo       UserInfo   `gorm:"references:UserId"`
	RecordInfoID   uint       `gorm:"uniqueIndex:recite_queue_idx"`
	RecordInfo     RecordInfo `gorm:"references:ID"`
	RemindTimeUnix int64      // next remind time in Unix
	Round          uint8      // 轮数
	RoundMax       uint8      // 终止轮数
	Status         int        // 状态
}
