package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type Monitor struct {
	gorm.Model
	CPULoad         uint
	MemLoad         uint
	RecordInput     uint
	RecordOutput    uint
	WordCountInput  uint
	WordCountOutput uint
}
