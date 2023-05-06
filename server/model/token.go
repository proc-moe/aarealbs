package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	UserInfoID uint     `gorm:"unique"`
	UserInfo   UserInfo `gorm:"references:UserId"`
	Token      string   `gorm:"unique"`
	Expire     int
}
