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
	"spider/src/luogu"
	updatecookie "spider/src/luogu/updateCookie"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var User genUrl.User
var countTime component.CountTime
var tempDB component.TempDB
var former int64 = 0

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	func() {
		c := resty.New()
		for i := 0; i < 100; i++ {
			resp, err := c.R().Get("https://laoin.work/api/blogs/848044096939?userId=552390554345")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(resp.Body()))
		}
	}()
	return
	db.Init()
	action1()
	defer db.Pool.Close()
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
			luogu.Use(10)
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
