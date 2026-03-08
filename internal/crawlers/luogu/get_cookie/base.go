package getcookie

import (
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type LuoguGetCookie struct {
	log        *easecrawler.EaseLogger
	cookiePool map[string]string
}

func NewLuoguGetCookie() *LuoguGetCookie {
	return &LuoguGetCookie{
		cookiePool: make(map[string]string),
	}
}

func (l *LuoguGetCookie) Name() string {
	return "luogu_get_cookie"
}
func (l *LuoguGetCookie) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         1 * time.Minute,
		StartImmediately: true,
	}
}
func (l *LuoguGetCookie) Run(c *easecrawler.Context) error {
	l.log = easecrawler.GetCrawlerLogger(c)
	l.Update()
	return nil
}
