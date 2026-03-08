package getrecorddetail

import (
	"spider/internal/repository"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetRecordDetail struct {
	log    *easecrawler.EaseLogger
	repo   *repository.LuoguRepository
	cookie string
	count  int
}

func NewGetRecordDetail() *GetRecordDetail {
	return &GetRecordDetail{
		repo:  repository.NewLuoguRepository(),
		count: 0,
	}
}
func (g *GetRecordDetail) Name() string {
	return "get_record_detail"
}
func (g *GetRecordDetail) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         1 * time.Minute,
		StartImmediately: true,
	}
}
func (g *GetRecordDetail) Run(c *easecrawler.Context) error {
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetAndStoreSourceCode()
	return nil
}
