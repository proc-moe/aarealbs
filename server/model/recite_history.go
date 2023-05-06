// id
// user_id - 外键
// record_id - 外键
// round
// time_gap
// time_gap_est
// result_status

package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type ReciteHistory struct {
	gorm.Model
	UserInfoID   uint
	UserInfo     UserInfo `gorm:"references:UserId"`
	RecordInfoID uint
	RecordInfo   RecordInfo
	TimeGap      uint
	TimeGapEst   uint
	Result       uint
}
