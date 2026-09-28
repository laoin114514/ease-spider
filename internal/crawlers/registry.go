package crawlers

import (
	cfGetOfficialProblem "spider/internal/crawlers/cf/get_official_problem"
	cfGetOfficialContest "spider/internal/crawlers/cf/get_ofifcial_contest"
	cfGetTeamContestProblems "spider/internal/crawlers/cf/get_team_contest_problems"
	cfGetTeamContests "spider/internal/crawlers/cf/get_team_contests"
	cfGetUserRecords "spider/internal/crawlers/cf/get_user_records"
	getcheckrecords "spider/internal/crawlers/dingding/get_check_records"
	luoguGetUserRecords "spider/internal/crawlers/luogu/get_user_records"
	niukeGetUserRecords "spider/internal/crawlers/niuke/get_user_records"

	easecrawler "spider/pkg/crawler"
)

func Register() *easecrawler.Engine {
	e := easecrawler.New()

	// 洛谷只保留提交记录：数据经 Luogu2Api 服务（pkg/luogu2api SDK）获取，
	// 取 cookie、提交详情、题解三个定时服务已随 SDK 接入一并删除
	luogu := e.Group("luogu")
	{
		luogu.Register(luoguGetUserRecords.NewGetUserRecords())
	}

	cf := e.Group("cf")
	{
		cf.Register(cfGetUserRecords.NewGetUserRecords())
		cf.Register(cfGetTeamContests.NewGetTeamContests())
		cf.Register(cfGetTeamContestProblems.NewGetTeamContestProblems())
		cf.Register(cfGetOfficialContest.NewGetOfifcialContest())
		cf.Register(cfGetOfficialProblem.NewGetOfifcialProblems())
	}

	dingding := e.Group("dingding")
	{
		dingding.Register(getcheckrecords.NewGetCheckRecords())
	}

	niuke := e.Group("niuke")
	{
		niuke.Register(niukeGetUserRecords.NewGetUserRecords())
	}

	return e
}
