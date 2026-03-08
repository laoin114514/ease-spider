package crawlers

import (
	getcookie "spider/internal/crawlers/luogu/get_cookie"
	easecrawler "spider/pkg/ease-crawler"
)

func Run() {
	e := easecrawler.New()
	luogu := e.Group("luogu")
	{
		luogu.Register(getcookie.NewLuoguGetCookie())
	}

	e.Run()
}
