package model

import (
	"github.com/proc-moe/aarealbs/server/utils/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error

	dsn := "proc-moe:123@tcp(localhost:3306)/remember?parseTime=true"
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}), &gorm.Config{})

	if err != nil {
		klog.E("[InitDB] sql open failed, err=%v", err)
		panic(err)
	}

	klog.I("[InitDB] connected")

	err = DB.AutoMigrate(&BatchInfo{}, &EffiencyStats{}, &Monitor{}, &ReciteHistory{})
	if err != nil {
		klog.E("[InitDB] sql open failed, err=%v", err)
		panic(err)
	}
	err = DB.AutoMigrate(&RecitePattern{}, &ReciteProcess{}, &ReciteQueue{}, &RecordInfo{})
	if err != nil {
		klog.E("[InitDB] sql open failed, err=%v", err)
		panic(err)
	}
	err = DB.AutoMigrate(&UserInfo{}, &Token{})
	if err != nil {
		klog.E("[InitDB] sql open failed, err=%v", err)
		panic(err)
	}

	// init data
	// default batch info
	if DB.Where("id = ?", 1).Find(&BatchInfo{}).RowsAffected == 0 {
		klog.I("no default batch, adding")
		v := BatchInfo{
			Title:       "Default Batch",
			Description: "默认的数据集",
			Author:      "system",
		}
		r := DB.Create(&v)
		klog.I("no default batch, adding, rows=%v", r.RowsAffected)
	} else {
		klog.I("default batch found, skip")
	}
}
