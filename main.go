package main

import (
	"fmt"
	"spider/cf/cfBetterSub"
	"spider/component"
	"spider/dingding"
	"spider/luogu"
	"time"
)

func main() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		if now.Hour() == 23 && now.Minute() == 59 {
			action()
		}
	}
}
func action() {
	var countTime component.CountTime
	fmt.Println(component.NowDateTime())
	countTime.Start()
	dingding.Use()
	luogu.Use()
	cfBetterSub.Use()
	countTime.End()
	fmt.Printf("\n\n")
}
func countDuration(start int64) string {
	now := time.Now()
	nowStamp := now.Unix()
	fmt.Println(now)
	duration := nowStamp - start
	day := duration / (3600 * 24)
	hour := (duration % (3600 * 24)) / 3600
	minute := (duration % 3600) / 60
	return fmt.Sprintf("%d天%d小时%d分钟", day, hour, minute)
}
