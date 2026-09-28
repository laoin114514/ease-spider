package getteamcontestproblems

import (
	"spider/config"
	"spider/internal/crawlers/cf/cfclient"
	"spider/internal/repository"
	"sync/atomic"
	"time"

	cf "github.com/laoin114514/codeforcesClient"
	easecrawler "github.com/laoin114514/ease-crawler"
)

type GetTeamContestProblems struct {
	log    *easecrawler.EaseLogger
	repo   *repository.CfRepository
	client *cf.Client
	// count 在并发器里累加，用原子量避免数据竞争
	count atomic.Int64
}

func NewGetTeamContestProblems() *GetTeamContestProblems {
	return &GetTeamContestProblems{
		repo: repository.NewCfRepository(),
	}
}
func (g *GetTeamContestProblems) Name() string {
	return "get_team_contest_problems"
}
func (g *GetTeamContestProblems) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         1 * time.Hour,
		StartImmediately: true,
	}
}
func (g *GetTeamContestProblems) Run(c *easecrawler.Context) error {
	keys, err := g.repo.GetCfApikeyPool()
	if err != nil {
		return err
	}
	g.client = cfclient.NewSigned(keys, config.AppConfig.Cf.CfTeamContestProblemsConcurrency)
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetCfTeamContestsProblems()
	return nil
}
