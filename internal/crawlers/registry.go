package crawlers

import (
	cfGetOfficialProblem "spider/internal/crawlers/cf/get_official_problem"
	cfGetOfficialContest "spider/internal/crawlers/cf/get_ofifcial_contest"
	cfGetTeamContestProblems "spider/internal/crawlers/cf/get_team_contest_problems"
	cfGetTeamContests "spider/internal/crawlers/cf/get_team_contests"
	cfGetUserRecords "spider/internal/crawlers/cf/get_user_records"
	getcheckrecords "spider/internal/crawlers/dingding/get_check_records"
	luoguGetCookie "spider/internal/crawlers/luogu/get_cookie"
	luoguGetRecordDetail "spider/internal/crawlers/luogu/get_record_detail"
	luoguGetSolutions "spider/internal/crawlers/luogu/get_solutions"
	luoguGetUserRecords "spider/internal/crawlers/luogu/get_user_records"
	niukeGetUserRecords "spider/internal/crawlers/niuke/get_user_records"

	easecrawler "github.com/laoin114514/ease-crawler"
)

func Register() *easecrawler.Engine {
	e := easecrawler.New()

	luogu := e.Group("luogu")
	{
		luogu.Register(luoguGetUserRecords.NewGetUserRecords())
		luogu.Register(luoguGetCookie.NewGetCookie())
		luogu.Register(luoguGetSolutions.NewGetSolutions())
		luogu.Register(luoguGetRecordDetail.NewGetRecordDetail())
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

	a := e.Group("cf")
	a.Group("")
	return e
}
