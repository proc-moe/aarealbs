package model

import (
	"fmt"

	"github.com/proc-moe/aarealbs/server/utils/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error

	dsn := "cyq:123@tcp(localhost:3306)/remember?parseTime=true"
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

	// INIT data
	klog.I(InitTokensDebug(114514, "YJSNPI", 2, "test"))
	klog.I(InitTokensDebug(114515, "user", 1, "normal"))

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

	// default batch info
	if DB.Where("p_id = ?", 0).Find(&RecitePattern{}).RowsAffected == 0 {
		klog.I("no default pattern, adding")
		x := []uint{5 * 60, 30 * 60, 12 * 3600, 86400, 2 * 86400, 4 * 86400, 7 * 86400, 15 * 86400}
		for k, v := range x {
			v := RecitePattern{
				PID:        0,
				Round:      uint(k),
				TimeGapEst: v,
			}
			r := DB.Create(&v)
			klog.I("no default pattern, adding, rows=%v", r.RowsAffected)
		}
	} else {
		klog.I("default batch found, skip")
	}
}

func InitTokensDebug(USERID uint, USERNAME string, STATUS uint, TOKEN string) string {
	userInfo := UserInfo{
		UserId:   USERID,
		UserName: USERNAME,
		Status:   STATUS,
	}
	token := Token{
		Token:      TOKEN,
		UserInfoID: USERID,
		Expire:     2147483647,
	}
	returnmsg := ""
	fmt.Println()
	if DB.Where("user_id = ?", USERID).First(&UserInfo{}).RowsAffected == 0 {
		r := DB.Create(&userInfo)
		fmt.Printf("[AddUser] ID, row affected = %v\n", r.RowsAffected)
	} else {
		returnmsg += "user already exists..."
	}

	if DB.Where("token = ?", TOKEN).Find(&Token{}).RowsAffected == 0 {
		r := DB.Create(&token)
		fmt.Printf("[AddToken] token, row affected = %v\n", r.RowsAffected)
	} else {
		returnmsg += "token already exists"
	}
	return returnmsg
}
