package getofifcialcontest

import (
	"spider/internal/crawlers/cf/cfclient"
	"spider/internal/repository"
	"time"

	cf "github.com/laoin114514/codeforcesClient"
	easecrawler "github.com/laoin114514/ease-crawler"
)

type GetOfifcialContest struct {
	log    *easecrawler.EaseLogger
	repo   *repository.CfRepository
	client *cf.Client
	count  int
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
	keys, err := g.repo.GetCfApikeyPool()
	if err != nil {
		return err
	}
	// 单个请求，不需要限流
	g.client = cfclient.NewSigned(keys, 0)
	g.log = easecrawler.GetCrawlerLogger(c)
	err = g.GetCfOfficialContests()
	if err != nil {
		return err
	}
	g.log.Printf("获取官方比赛完成 %d", g.count)
	return nil
}
