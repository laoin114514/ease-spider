package handler

import (
	"log"
	"spider/config"
	"spider/src/services"
	"spider/src/utils"
	"strconv"
	"time"
)

func analysisTimerFrequency(s string) time.Duration {
	idx := 0
	for ; idx < len(s); idx++ {
		if !(s[idx] >= '0' && s[idx] <= '9') {
			break
		}
	}
	duration, err := strconv.Atoi(s[:idx])
	if err != nil {
		return 0
	}
	if s[idx:] == "s" {
		return time.Second * time.Duration(duration)
	} else if s[idx:] == "m" {
		return time.Minute * time.Duration(duration)
	} else if s[idx:] == "h" {
		return time.Hour * time.Duration(duration)
	} else if s[idx:] == "d" {
		return time.Hour * 24 * time.Duration(duration)
	} else {
		return time.Minute * time.Duration(duration)
	}
}
func TimeTask() {
	timer := utils.NewTimer()
	cfService := services.NewCfService()
	luoguService := services.NewLuoguService()
	dingService := services.NewDingdingService()
	luoguUpdateCookieService := services.NewLuoguUpdateCookie()

	//每小时执行一次，更新洛谷Cookie
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.LuoguUpdateCookie),
		"更新洛谷Cookie",
		func() error {
			luoguUpdateCookieService.UpdateLuoguCookie()
			luoguUpdateCookieService.SaveLog()
			luoguUpdateCookieService.Clear()
			return nil
		},
	)
	//获取cf官方题目
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.CfOfficialProblems),
		"获取cf官方题目",
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
		})
	//获取cf官方比赛
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.CfOfficialContests),
		"获取cf官方比赛",
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
		})
	//获取cf团队题目
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.CfTeamContestsProblems),
		"获取cf团队题目",
		func() error {
			err := cfService.CfTeamContestsProblems.GetCfTeamContestsProblems()
			if err != nil {
				return err
			}
			cfService.CfTeamContestsProblems.SaveLog()
			cfService.CfTeamContestsProblems.SaveErr()
			cfService.CfTeamContestsProblems.Clear()
			log.Println("cf团队题目获取完成")
			return nil
		})
	//获取cf团队比赛
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.CfTeamContests),
		"获取cf团队比赛",
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
		})
	//获取钉钉打卡数据
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.Dingding),
		"获取钉钉打卡数据",
		func() error {
			err := dingService.GetDingdingCheckUpData()
			if err != nil {
				return err
			}
			dingService.SaveLog()
			dingService.SaveErr()
			dingService.Clear()
			log.Println("钉钉打卡获取完成")
			return nil
		})
	//获取Cf提交记录
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.CfRecords),
		"获取cf提交记录",
		func() error {
			err := cfService.CfUserStatus.GetCfRecords()
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
		analysisTimerFrequency(config.AppConfig.TimerFrequency.LuoguRecords),
		"获取洛谷用户提交记录",
		func() error {
			err := luoguService.GetLuoguUsersRecords()
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
