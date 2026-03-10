package getuserrecords

import (
	"spider/internal/repository"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetUserRecords struct {
	log  *easecrawler.EaseLogger
	repo *repository.LuoguRepository
}

func NewGetUserRecords() *GetUserRecords {
	return &GetUserRecords{
		repo: repository.NewLuoguRepository(),
	}
}
func (g *GetUserRecords) Name() string {
	return "get_user_records"
}
func (g *GetUserRecords) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         5 * time.Minute,
		StartImmediately: true,
	}
}
func (g *GetUserRecords) Run(c *easecrawler.Context) error {
	g.log = easecrawler.GetCrawlerLogger(c)
	g.log.Println("开始获取用户记录")
	g.GetAndStore()
	g.ChangePrivateProblem()
	return nil
}
