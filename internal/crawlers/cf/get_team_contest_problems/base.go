package getteamcontestproblems

import (
	"spider/internal/repository"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetTeamContestProblems struct {
	log   *easecrawler.EaseLogger
	repo  *repository.CfRepository
	count int
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
		Interval:         1 * time.Minute,
		StartImmediately: true,
	}
}
func (g *GetTeamContestProblems) Run(c *easecrawler.Context) error {
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetCfTeamContestsProblems()
	return nil
}
