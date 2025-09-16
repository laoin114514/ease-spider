package main

import (
	"fmt"
	"spider/component"
	"spider/config/db"
	"spider/src/cf/cfBetterSub"
	cfOfficial "spider/src/cf/cf_official"
	"spider/src/cf/genUrl"
	teamtainning "spider/src/cf/teamTainning"
	"spider/src/cf/teamquestions"
	"spider/src/dingding"
	"spider/src/luogu"
	updatecookie "spider/src/luogu/updateCookie"
	"time"
)

var User genUrl.User
var countTime component.CountTime
var tempDB component.TempDB
var former int64 = 0

func main() {
	db.Init()
	defer db.Pool.Close()
	action()
	return
	now := time.Now()
	tempDB.Set("startTime", now.Unix())

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now = time.Now()
		current := now.Unix()
		if (now.Day() == 1 || now.Day() == 10 || now.Day() == 20) && now.Hour() == 23 && now.Minute() == 59 {
		}
		if (current - former) > 5*3000 {
			fmt.Println(component.NowDateTime())
			updatecookie.Use()
			component.SenEamil("2908451607@qq.com")
			former = current
			action()
		}
	}
}
func action() {
	countTime.Start()

	//功能区
	dingding.Use()
	teamtainning.Use()
	teamquestions.Use()
	cfOfficial.Use()
	cfBetterSub.Use()
	luogu.Use()

	countTime.End()
}
