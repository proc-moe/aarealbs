package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type Monitor struct {
	gorm.Model
	CPULoad     float64
	MemLoad     float64
	RecordCount int
	UserCount   int
}
