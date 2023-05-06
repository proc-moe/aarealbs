package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type RecitePattern struct {
	gorm.Model
	PID        uint `gorm:"uniqueIndex:pattern_idx"`
	Round      uint `gorm:"uniqueIndex:pattern_idx"`
	TimeGapEst uint
}
