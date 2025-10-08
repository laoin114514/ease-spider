package handler

import (
	"log"
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
		//获取cf官方题目
		func() error {
			err := cfService.CfOfficialProblems.GetCfOfficialProblems()
			if err != nil {
				return err
			}
			cfService.CfOfficialProblems.SaveLog()
			cfService.CfOfficialProblems.SaveErr()
			cfService.CfOfficialProblems.Clear()
			log.Println("cf官方题目获取完成")
			return nil
		},
		//获取cf官方比赛
		func() error {
			err := cfService.CfOfficialContests.GetCfOfficialContests()
			if err != nil {
				return err
			}
			cfService.CfOfficialContests.SaveLog()
			cfService.CfOfficialContests.SaveErr()
			cfService.CfOfficialContests.Clear()
			log.Println("cf官方比赛获取完成")
			return nil
		},
		//获取cf团队比赛
		func() error {
			err := cfService.CfTeamContests.GetCfTeamContests()
			if err != nil {
				return err
			}
			cfService.CfTeamContests.SaveLog()
			cfService.CfTeamContests.SaveErr()
			cfService.CfTeamContests.Clear()
			log.Println("cf团队比赛获取完成")
			return nil
		},
		//获取cf团队比赛题目
		func() error {
			err := cfService.CfTeamContestsProblems.GetCfTeamContestsProblems(3)
			if err != nil {
				return err
			}
			cfService.CfTeamContestsProblems.SaveLog()
			cfService.CfTeamContestsProblems.SaveErr()
			cfService.CfTeamContestsProblems.Clear()
			log.Println("cf团队比赛题目获取完成")
			return nil
		},
		//获取钉钉打卡数据
		func() error {
			err := dingService.GetDingdingCheckUpData()
			if err != nil {
				return err
			}
			dingService.SaveLog()
			dingService.SaveErr()
			dingService.Clear()
			log.Println("钉钉打卡数据获取完成")
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
			cfService.CfUserStatus.SaveLog()
			cfService.CfUserStatus.SaveErr()
			cfService.CfUserStatus.Clear()
			log.Println("cf提交记录获取完成")
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
			luoguService.SaveLog()
			luoguService.SaveErr()
			luoguService.Clear()
			log.Println("洛谷用户提交记录获取完成")
			return nil
		})
}
