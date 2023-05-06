package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	UserId   uint   `gorm:"unique"` // 用户 ID
	UserName string // 用户名
	Status   uint   // 状态
}
