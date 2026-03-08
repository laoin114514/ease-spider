package getcookie

import (
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetCookie struct {
	log        *easecrawler.EaseLogger
	cookiePool map[string]string
}

func NewGetCookie() *GetCookie {
	return &GetCookie{
		cookiePool: make(map[string]string),
	}
}

func (l *GetCookie) Name() string {
	return "get_cookie"
}
func (l *GetCookie) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         24 * time.Hour,
		StartImmediately: true,
	}
}
func (l *GetCookie) Run(c *easecrawler.Context) error {
	l.log = easecrawler.GetCrawlerLogger(c)
	l.Update()
	return nil
}
