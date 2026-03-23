package getofifcialcontest

import (
	"spider/internal/repository"
	cfurlgenerator "spider/pkg/cf-url-generator"
	"time"

	easecrawler "github.com/laoin114514/ease-crawler"
)

type GetOfifcialContest struct {
	log          *easecrawler.EaseLogger
	repo         *repository.CfRepository
	urlGenerator *cfurlgenerator.GenerateCFurl
	count        int
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
	apiKeyPool, err := g.repo.GetCfApikeyPool()
	if err != nil {
		return err
	}
	g.urlGenerator = cfurlgenerator.NewGenerator(apiKeyPool)

	g.log = easecrawler.GetCrawlerLogger(c)
	err = g.GetCfOfficialContests()
	if err != nil {
		return err
	}
	g.log.Printf("获取官方比赛完成 %d", g.count)
	return nil
}
