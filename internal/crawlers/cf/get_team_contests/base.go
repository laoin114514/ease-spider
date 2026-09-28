package getteamcontests

import (
	"spider/config"
	"spider/internal/crawlers/cf/cfclient"
	"spider/internal/repository"
	"time"

	cf "github.com/laoin114514/codeforcesClient"
	easecrawler "spider/pkg/crawler"
)

type GetTeamContests struct {
	log        *easecrawler.EaseLogger
	repo       *repository.CfRepository
	client     *cf.Client
	useAccount string
	count      int
}

func NewGetTeamContests() *GetTeamContests {
	return &GetTeamContests{
		repo:       repository.NewCfRepository(),
		count:      0,
		useAccount: config.AppConfig.Cf.ManagerAccount,
	}
}
func (g *GetTeamContests) Name() string {
	return "get_team_contests"
}
func (g *GetTeamContests) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         config.IntervalOr(config.TimerFrequency().CfTeamContests, 1*time.Hour),
		StartImmediately: true,
	}
}
func (g *GetTeamContests) Run(c *easecrawler.Context) error {
	keys, err := g.repo.GetCfApikeyPool()
	if err != nil {
		return err
	}
	// 单个请求，不需要限流
	g.client = cfclient.NewSigned(keys, 0)
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetCfTeamContests()
	return nil
}
