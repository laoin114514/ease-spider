package main

import (
	"fmt"
	"spider/component"
	"spider/config/db"
	"spider/src/cf/cfBetterSub"
	"spider/src/cf/cf_official_problems"
	"spider/src/cf/genUrl"
	"spider/src/cf/team_contest_problems"
	teamtainning "spider/src/cf/team_contests"
	"spider/src/dingding"
	"spider/src/services"
	"spider/src/utils"
	updatecookie "spider/src/utils/updateCookie"
	"sync"
	"time"
)

var User genUrl.User
var countTime component.CountTime
var tempDB component.TempDB
var former int64 = 0

func main() {
	utils.InitGlobalJSONDB("tempDB.json")
	timer := utils.NewTimer()
	luoguService := services.NewLuoguService()
	dingdingService := services.NewDingdingService()
	wg := sync.WaitGroup{}
	db.Init()
	defer db.Pool.Close()
	dingdingService.GetDingdingCheckUpData()
	fmt.Println(dingdingService.GetLog())
	fmt.Println(dingdingService.GetErr())
	return
	// updatecookie.Use(false)
	str := timer.CountDurationStr(func() {
		luoguService.GetLuoguUsersRecords(10)
		luoguService.Log()
	})
	fmt.Println(str)
	return
	now := time.Now()
	tempDB.Set("startTime", now.Unix())
	go func() {
		//该计时器按天爬取打卡、官方题、团队赛等
		fmt.Println("计时器1启动")
		updatecookie.Use(false)
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			action1()
		}
	}()
	go func() {
		//该计时器爬取洛谷过题记录
		fmt.Println("计时器2启动")
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			luoguService.GetLuoguUsersRecords(1)
		}
	}()
	go func() {
		//该计时器爬取cf提交记录
		fmt.Println("计时器3启动")
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cfBetterSub.Use(5)
		}
	}()

	wg.Wait()
}
func action1() {
	var wg sync.WaitGroup
	countTime.Start()
	fmt.Println(component.NowDateTime())
	updatecookie.Use(false)
	component.SenEamil("3247428622@qq.com")
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
		team_contest_problems.Use(3)
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		cf_official_problems.Use()
		wg.Done()
	}()
	time.Sleep(2 * time.Second) //等待wg.Add生效
	wg.Wait()

	countTime.End()
}
