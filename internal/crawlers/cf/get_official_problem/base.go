package getofficialproblem

import (
	"spider/config"
	"spider/internal/crawlers/cf/cfclient"
	"time"

	cf "github.com/laoin114514/codeforcesClient"
	easecrawler "spider/pkg/crawler"
)

type GetOfifcialProblems struct {
	log    *easecrawler.EaseLogger
	client *cf.Client
	count  int
}

func NewGetOfifcialProblems() *GetOfifcialProblems {
	return &GetOfifcialProblems{
		count: 0,
	}
}

func (g *GetOfifcialProblems) Name() string {
	return "get_ofifcial_problems"
}
func (g *GetOfifcialProblems) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         config.IntervalOr(config.TimerFrequency().CfOfficialProblems, 1*time.Hour),
		StartImmediately: true,
	}
}
func (g *GetOfifcialProblems) Run(c *easecrawler.Context) error {
	// problemset.problems 是公共接口，不需要凭据
	g.client = cfclient.NewPlain(0)
	g.log = easecrawler.GetCrawlerLogger(c)

	err := g.GetCfOfficialProblems()
	if err != nil {
		return err
	}
	g.log.Printf("获取官方题目完成 %d", g.count)
	return nil
}
