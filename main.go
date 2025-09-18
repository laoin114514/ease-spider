package main

import (
	"fmt"
	"spider/component"
	"spider/config/db"
	"spider/src/cf/cfBetterSub"
	cfOfficial "spider/src/cf/cf_official"
	"spider/src/cf/genUrl"
	"spider/src/cf/team_contest_problems"
	teamtainning "spider/src/cf/team_contests"
	"spider/src/dingding"
	"spider/src/luogu"
	updatecookie "spider/src/luogu/updateCookie"
	"sync"
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

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now = time.Now()
		current := now.Unix()
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
	var wg sync.WaitGroup
	countTime.Start()

	//功能区
	go func() {
		wg.Add(1)
		dingding.Use()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		teamtainning.Use()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		team_contest_problems.Use()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		cfOfficial.Use()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		cfBetterSub.Use()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		luogu.Use()
		wg.Done()
	}()
	time.Sleep(2 * time.Second)
	wg.Wait()
	countTime.End()
}
