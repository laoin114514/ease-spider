package getteamcontestproblems

import (
	"spider/internal/repository"
	cfurlgenerator "spider/pkg/cf-url-generator"
	"time"

	easecrawler "github.com/laoin114514/ease-crawler"
)

type GetTeamContestProblems struct {
	log          *easecrawler.EaseLogger
	repo         *repository.CfRepository
	urlGenerator *cfurlgenerator.GenerateCFurl
	count        int
}

func NewGetTeamContestProblems() *GetTeamContestProblems {
	return &GetTeamContestProblems{
		repo:  repository.NewCfRepository(),
		count: 0,
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
	apiKeyPool, err := g.repo.GetCfApikeyPool()
	if err != nil {
		return err
	}
	g.urlGenerator = cfurlgenerator.NewGenerator(apiKeyPool)
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetCfTeamContestsProblems()
	return nil
}
