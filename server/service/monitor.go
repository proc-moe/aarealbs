package service

import (
	"fmt"

	"github.com/shirou/gopsutil/load"

	"github.com/proc-moe/aarealbs/server/model"
	cron "github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/mem"
)

func MonitorStart() {
	cronTab := cron.New()
	task := func() {
		fmt.Println("start")
		var recordCount int64
		model.DB.Model(&model.RecordInfo{}).Count(&recordCount)

		var userCount int64
		model.DB.Model(&model.UserInfo{}).Count(&userCount)

		mem, _ := mem.VirtualMemory()

		cpu, _ := load.Avg()

		v := model.Monitor{
			CPULoad:     cpu.Load1,
			MemLoad:     mem.UsedPercent,
			RecordCount: int(recordCount),
			UserCount:   int(userCount),
		}
		model.DB.Create(&v)
	}

	cronTab.AddFunc("0/5 * * * * ? *", task)
	cronTab.Start()
}
