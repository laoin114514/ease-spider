package getuserrecords

import (
	"spider/internal/repository"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetUserRecords struct {
	log  *easecrawler.EaseLogger
	repo *repository.CfRepository
}

func NewGetUserRecords() *GetUserRecords {
	return &GetUserRecords{
		repo: repository.NewCfRepository(),
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
	g.GetCfRecords()
	return nil
}
