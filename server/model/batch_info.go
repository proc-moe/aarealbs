package model

// use gorm, better than raw query
import (
	"gorm.io/gorm"
)

type BatchInfo struct {
	gorm.Model
	Title       string
	Description string
	Author      string
}
