package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type EffiencyStats struct {
	gorm.Model
	UserInfoID   uint
	UserInfo     UserInfo `gorm:"references:UserId"`
	ForgetRate   float32
	ResponseTime float32
	ReciteTry    int
}
