package getteamcontests

import (
	"spider/config"
	"spider/internal/repository"
	cfurlgenerator "spider/pkg/cf-url-generator"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetTeamContests struct {
	log          *easecrawler.EaseLogger
	repo         *repository.CfRepository
	urlGenerator *cfurlgenerator.GenerateCFurl
	useAccount   string
	count        int
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
		Interval:         1 * time.Hour,
		StartImmediately: true,
	}
}
func (g *GetTeamContests) Run(c *easecrawler.Context) error {
	apiKeyPool, err := g.repo.GetCfApikeyPool()
	if err != nil {
		return err
	}
	g.urlGenerator = cfurlgenerator.NewGenerator(apiKeyPool)
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetCfTeamContests()
	return nil
}
