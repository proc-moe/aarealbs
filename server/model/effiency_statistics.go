package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type EffiencyStats struct {
	gorm.Model
	UserInfoID          uint
	UserInfo            UserInfo `gorm:"references:UserId"`
	AverageRememberTime string
	ForgetRate          string
	ResponseTime        string
}
