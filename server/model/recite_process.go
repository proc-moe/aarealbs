package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type ReciteProcess struct {
	gorm.Model
	UserInfoID       uint
	UserInfo         UserInfo `gorm:"references:UserId"`
	ReciteDone       uint
	ReciteProcessing uint
	ReciteTODO       uint
}
