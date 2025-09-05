package main

import (
	"fmt"
	"spider/cf/cfBetterSub"
	cfofficial "spider/cf/cf_official"
	"spider/component"
	"spider/dingding"
	"spider/newLuogu"
	"time"
)

var countTime component.CountTime
var tempDB component.TempDB

func main() {
	cfofficial.Use()
	return
	now := time.Now()
	tempDB.Set("startTime", now.Unix())
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	fmt.Println("开始计时")
	var former int64 = 0
	for range ticker.C {
		now = time.Now()
		current := now.Unix()
		if (current - former) > 5*3000 {
			component.SenEamil("2908451607@qq.com")
			former = current
			action()
		}
	}
}
func action() {
	fmt.Println(component.NowDateTime())
	countTime.Start()

	//功能区
	dingding.Use()
	newLuogu.Use()
	cfBetterSub.Use()

	countTime.End()
	fmt.Printf("\n\n")
}
