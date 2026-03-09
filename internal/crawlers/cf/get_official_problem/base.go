package getofficialproblem

import (
	"spider/internal/repository"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetOfifcialProblems struct {
	log   *easecrawler.EaseLogger
	repo  *repository.CfRepository
	count int
}

func NewGetOfifcialProblems() *GetOfifcialProblems {
	return &GetOfifcialProblems{
		repo:  repository.NewCfRepository(),
		count: 0,
	}
}

func (g *GetOfifcialProblems) Name() string {
	return "get_ofifcial_problems"
}
func (g *GetOfifcialProblems) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         1 * time.Hour,
		StartImmediately: true,
	}
}
func (g *GetOfifcialProblems) Run(c *easecrawler.Context) error {
	g.log = easecrawler.GetCrawlerLogger(c)
	err := g.GetCfOfficialProblems()
	if err != nil {
		return err
	}
	g.log.Printf("获取官方题目完成 %d", g.count)
	return nil
}
