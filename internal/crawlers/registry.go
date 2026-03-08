package crawlers

import (
	luoguGetCookie "spider/internal/crawlers/luogu/get_cookie"
	luoguGetUserRecords "spider/internal/crawlers/luogu/get_user_records"
	easecrawler "spider/pkg/ease-crawler"
)

func Run() {
	e := easecrawler.New()
	luogu := e.Group("luogu")
	{
		luogu.Register(luoguGetUserRecords.NewGetUserRecords())
		luogu.Register(luoguGetCookie.NewGetCookie())
	}

	e.Run()
}
