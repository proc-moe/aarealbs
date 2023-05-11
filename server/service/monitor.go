package service

import (
	"fmt"
	"time"

	"github.com/proc-moe/aarealbs/server/model"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func MonitorStart() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				{
					fmt.Println("system status recorded")
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
			}
		}
	}()
}
