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
	dingService := services.NewDingdingService()
	luoguRecordsService := services.NewLuoguRecordsService()
	luoguUpdateCookieService := services.NewLuoguUpdateCookie()
	// luoguSubmissionDetailService := services.NewLuoguSubmissionDetail()
	//每小时执行一次，更新洛谷Cookie
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.LuoguUpdateCookie),
		"更新洛谷Cookie",
		func() error {
			log.Println("开始更新洛谷Cookie")
			err := luoguUpdateCookieService.UpdateLuoguCookie()
			if err != nil {
				return err
			}
			luoguUpdateCookieService.SaveLog()
			luoguUpdateCookieService.SaveLog()
			luoguUpdateCookieService.Clear()
			return nil
		},
	)
	//获取洛谷提交记录源代码
	// timer.RunWithTimer(
	// 	analysisTimerFrequency(config.AppConfig.TimerFrequency.LuoguSubmissionDetail),
	// 	"获取洛谷提交记录源代码",
	// 	func() error {
	// 		log.Println("开始获取洛谷提交记录源代码")
	// 		err := luoguSubmissionDetailService.GetRecordSourceCode()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		luoguSubmissionDetailService.SaveLog()
	// 		luoguSubmissionDetailService.SaveErr()
	// 		luoguSubmissionDetailService.Clear()
	// 		log.Println("洛谷提交记录源代码获取完成")
	// 		return nil
	// 	})
	//获取cf官方题目
	timer.RunWithTimer(
		analysisTimerFrequency(config.AppConfig.TimerFrequency.CfOfficialProblems),
		"获取cf官方题目",
		func() error {
			log.Println("开始获取cf官方题目")
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
			log.Println("开始获取cf官方比赛")
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
			log.Println("开始获取cf团队题目")
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
			log.Println("开始获取cf团队比赛")
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
			log.Println("开始获取钉钉打卡数据")
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
			log.Println("开始获取cf提交记录")
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
			log.Println("开始获取洛谷用户提交记录")
			err := luoguRecordsService.GetLuoguUsersRecords()
			if err != nil {
				return err
			}
			err = luoguRecordsService.ChangePrivateProblem()
			if err != nil {
				return err
			}
			luoguRecordsService.SaveLog()
			luoguRecordsService.SaveErr()
			luoguRecordsService.Clear()
			log.Println("洛谷用户提交记录获取完成")
			return nil
		})
}
