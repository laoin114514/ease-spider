package main

import (
	"fmt"
	"spider/cf/cfBetterSub"
	"spider/component"
	"spider/dingding"
	"spider/luogu"
	"spider/newLuogu"
	"time"
)

var countTime component.CountTime
var tempDB component.TempDB

func main() {
	newLuogu.Use()
	return
	now := time.Now()
	tempDB.Set("startTime", now.Unix())
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	fmt.Println("开始计时")
	for range ticker.C {
		now = time.Now()
		//每天十二点收集数据
		if now.Hour() == 23 && now.Minute() == 59 {
			action()
		}
	}
}
func action() {
	fmt.Println(component.NowDateTime())
	countTime.Start()
	//功能区

	dingding.Use()
	luogu.Use()
	cfBetterSub.Use()

	//功能区
	countTime.End()
	fmt.Printf("\n\n")
}
