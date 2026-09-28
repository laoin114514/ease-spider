package getcheckrecords

import (
	"spider/internal/repository"
	"time"

	easecrawler "spider/pkg/crawler"
)

type GetCheckRecords struct {
	log         *easecrawler.EaseLogger
	repo        *repository.DingdingRepository
	token       string
	dingUserMap map[string]any
	insertCount int
}

func NewGetCheckRecords() *GetCheckRecords {
	return &GetCheckRecords{
		repo:        repository.NewDingdingRepository(),
		token:       "",
		dingUserMap: make(map[string]any),
		insertCount: 0,
	}
}
func (g *GetCheckRecords) Name() string {
	return "get_check_records"
}
func (g *GetCheckRecords) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         1 * time.Hour,
		StartImmediately: true,
	}
}
func (g *GetCheckRecords) Run(c *easecrawler.Context) error {
	g.log = easecrawler.GetCrawlerLogger(c)
	err := g.GetDingdingCheckUpData()
	if err != nil {
		return err
	}
	g.log.Printf("获取钉钉打卡数据完成 %d", g.insertCount)
	return nil
}
