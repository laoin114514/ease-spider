package crawlers

import (
	cfGetUserRecords "spider/internal/crawlers/cf/get_user_records"
	luoguGetCookie "spider/internal/crawlers/luogu/get_cookie"
	luoguGetRecordDetail "spider/internal/crawlers/luogu/get_record_detail"
	luoguGetSolutions "spider/internal/crawlers/luogu/get_solutions"
	luoguGetUserRecords "spider/internal/crawlers/luogu/get_user_records"
	easecrawler "spider/pkg/ease-crawler"
)

func Run() {
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
	}

	e.Run()
}
