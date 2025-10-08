package handler

import (
	"fmt"
	"spider/src/services"
	"spider/src/utils"
	"time"
)

func TimeTask() {
	timer := utils.NewTimer()
	cfService := services.NewCfService()
	luoguService := services.NewLuoguService()
	dingService := services.NewDingdingService()
	luoguUpdateCookieService := services.NewLuoguUpdateCookie()
	luoguUpdateCookieService.UpdateLuoguCookie()
	//每小时执行一次，更新洛谷Cookie
	timer.RunWithTimer(
		time.Hour,
		func() error {
			luoguUpdateCookieService.UpdateLuoguCookie()
			return nil
		},
	)
	//每两小时执行一次，非高速度要求任务
	timer.MultiRunWithTimer(
		time.Hour*2,
		//获取cf官方比赛
		func() error {
			err := cfService.CfOfficialContests.GetCfOfficialContests()
			if err != nil {
				return err
			}
			fmt.Println(cfService.CfOfficialContests.GetLog())
			fmt.Println(cfService.CfOfficialContests.GetErr())
			cfService.CfOfficialContests.Clear()
			return nil
		},
		//获取cf团队比赛
		func() error {
			err := cfService.CfTeamContests.GetCfTeamContests()
			if err != nil {
				return err
			}
			fmt.Println(cfService.CfTeamContests.GetLog())
			fmt.Println(cfService.CfTeamContests.GetErr())
			cfService.CfTeamContests.Clear()
			return nil
		},
		//获取cf团队比赛题目
		func() error {
			err := cfService.CfTeamContestsProblems.GetCfTeamContestsProblems(3)
			if err != nil {
				return err
			}
			fmt.Println(cfService.CfTeamContestsProblems.GetLog())
			fmt.Println(cfService.CfTeamContestsProblems.GetErr())
			cfService.CfTeamContestsProblems.Clear()
			return nil
		},
		//获取钉钉打卡数据
		func() error {
			err := dingService.GetDingdingCheckUpData()
			if err != nil {
				return err
			}
			fmt.Println(dingService.GetLog())
			fmt.Println(dingService.GetErr())
			dingService.Clear()
			return nil
		},
	)
	//获取Cf提交记录
	timer.RunWithTimer(
		time.Second*120,
		func() error {
			err := cfService.CfUserStatus.GetCfRecords(3)
			if err != nil {
				return err
			}
			fmt.Println(cfService.CfUserStatus.GetLog())
			fmt.Println(cfService.CfUserStatus.GetErr())
			cfService.CfUserStatus.Clear()
			return nil
		})
	//获取洛谷用户提交记录
	timer.RunWithTimer(
		time.Second*120,
		func() error {
			err := luoguService.GetLuoguUsersRecords(5)
			if err != nil {
				return err
			}
			fmt.Println(luoguService.GetLog())
			fmt.Println(luoguService.GetErr())
			luoguService.Clear()
			return nil
		})
}
