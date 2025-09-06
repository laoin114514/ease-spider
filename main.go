package main

import (
	"fmt"
	"spider/cf/cfBetterSub"
	cfOfficial "spider/cf/cf_official"
	"spider/cf/genUrl"
	"spider/component"
	"spider/db"
	"spider/dingding"
	"spider/luogu"
	"time"
)

var User genUrl.User
var countTime component.CountTime
var tempDB component.TempDB
var former int64 = 0

func main() {
	db.Init()
	defer db.Pool.Close()

	now := time.Now()
	tempDB.Set("startTime", now.Unix())

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

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
	cfBetterSub.Use()
	cfOfficial.Use()
	dingding.Use()
	luogu.Use()

	countTime.End()
	fmt.Printf("\n\n")
}
