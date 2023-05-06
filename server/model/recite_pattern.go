package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type RecitePattern struct {
	gorm.Model
	PID        uint `gorm:"unique"`
	Round      uint
	TimeGapEst uint
}
