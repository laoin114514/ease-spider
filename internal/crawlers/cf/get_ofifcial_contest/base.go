package getofifcialcontest

import (
	"spider/internal/repository"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetOfifcialContest struct {
	log   *easecrawler.EaseLogger
	repo  *repository.CfRepository
	count int
}

func NewGetOfifcialContest() *GetOfifcialContest {
	return &GetOfifcialContest{
		repo:  repository.NewCfRepository(),
		count: 0,
	}
}
func (g *GetOfifcialContest) Name() string {
	return "get_ofifcial_contest"
}
func (g *GetOfifcialContest) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         1 * time.Hour,
		StartImmediately: true,
	}
}
func (g *GetOfifcialContest) Run(c *easecrawler.Context) error {
	g.log = easecrawler.GetCrawlerLogger(c)
	err := g.GetCfOfficialContests()
	if err != nil {
		return err
	}
	g.log.Printf("获取官方比赛完成 %d", g.count)
	return nil
}
