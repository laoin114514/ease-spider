package main

import (
	"fmt"
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

func main() {
	fmt.Println(User.Status(true, genUrl.User_status{
		Handle:         "233zhang",
		From:           1,
		Count:          50000,
		IncludeSources: false,
	}))
	db.Init()
	defer db.Pool.Close()
	action()
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
	luogu.Use()
	// cfBetterSub.Use()

	countTime.End()
	fmt.Printf("\n\n")
}
