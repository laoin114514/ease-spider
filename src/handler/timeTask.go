package handler

import (
	"log"
	"spider/config"
	"spider/src/services"
	"spider/src/utils"
)

func TimeTask() {

	timer := utils.NewTimer()
	cfService := services.NewCfService()
	dingService := services.NewDingdingService()
	luogu := services.NewLuogu()
	timer.RunWithTimer(
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.LuoguUpdateCookie),
		"更新洛谷Cookie",
		func() error {
			log.Println("开始更新洛谷Cookie")
			err := luogu.LuoguCookie.Update()
			if err != nil {
				return err
			}
			luogu.LuoguCookie.SaveLog()
			luogu.LuoguCookie.SaveErr()
			luogu.LuoguCookie.Clear()
			return nil
		},
	)
	timer.RunWithTimer(
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.CfOfficialProblems),
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
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.CfOfficialContests),
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
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.CfTeamContestsProblems),
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
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.CfTeamContests),
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
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.Dingding),
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
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.CfRecords),
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
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.LuoguRecords),
		"获取洛谷用户提交记录",
		func() error {
			log.Println("开始获取洛谷用户提交记录")
			err := luogu.LuoguRecords.GetAndStore()
			if err != nil {
				return err
			}
			err = luogu.LuoguRecords.ChangePrivateProblem()
			if err != nil {
				return err
			}
			luogu.LuoguRecords.SaveLog()
			luogu.LuoguRecords.SaveErr()
			luogu.LuoguRecords.Clear()
			log.Println("洛谷用户提交记录获取完成")
			return nil
		})
	//获取洛谷题解
	timer.RunWithTimer(
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.LuoguSolution),
		"获取洛谷题解",
		func() error {
			log.Println("开始获取洛谷题解")
			err := luogu.LuoguSolution.GetAndStore()
			if err != nil {
				return err
			}
			luogu.LuoguSolution.SaveLog()
			luogu.LuoguSolution.SaveErr()
			luogu.LuoguSolution.Clear()
			log.Println("洛谷题解获取完成")
			return nil
		})
	//获取洛谷源代码
	timer.RunWithTimer(
		utils.NewDateFormat().AnalysisTimerFrequency(config.AppConfig.TimerFrequency.LuoguSubmissionDetail),
		"获取洛谷源代码",
		func() error {
			log.Println("开始获取洛谷源代码")
			err := luogu.LuoguSubmissionDetail.GetAndStoreSourceCode()
			if err != nil {
				return err
			}
			luogu.LuoguSubmissionDetail.SaveLog()
			luogu.LuoguSubmissionDetail.SaveErr()
			luogu.LuoguSubmissionDetail.Clear()
			log.Println("洛谷源代码获取完成")
			return nil
		})
}
