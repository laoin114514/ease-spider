package getuserrecords

import (
	"spider/internal/repository"
	cfurlgenerator "spider/pkg/cf-url-generator"
	"time"

	easecrawler "github.com/laoin114514/ease-crawler"
)

type GetUserRecords struct {
	log          *easecrawler.EaseLogger
	repo         *repository.CfRepository
	urlGenerator *cfurlgenerator.GenerateCFurl
}

func NewGetUserRecords() *GetUserRecords {
	return &GetUserRecords{
		repo: repository.NewCfRepository(),
	}
}
func (g *GetUserRecords) Name() string {
	return "get_user_records"
}
func (g *GetUserRecords) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         5 * time.Minute,
		StartImmediately: true,
	}
}
func (g *GetUserRecords) Run(c *easecrawler.Context) error {
	apiKeyPool, err := g.repo.GetCfApikeyPool()
	if err != nil {
		return err
	}
	g.urlGenerator = cfurlgenerator.NewGenerator(apiKeyPool)
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetCfRecords()
	return nil
}
