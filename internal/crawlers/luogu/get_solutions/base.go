package getsolutions

import (
	"spider/internal/models"
	"spider/internal/repository"
	easecrawler "spider/pkg/ease-crawler"
	"time"
)

type GetSolutions struct {
	log       *easecrawler.EaseLogger
	repo      *repository.LuoguRepository
	totalPage int
	solutions []models.SolutionContent
	count     int
}

func NewGetSolutions() *GetSolutions {
	return &GetSolutions{
		repo: repository.NewLuoguRepository(),
	}
}
func (g *GetSolutions) Name() string {
	return "get_solutions"
}
func (g *GetSolutions) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         1 * time.Minute,
		StartImmediately: true,
	}
}
func (g *GetSolutions) Run(c *easecrawler.Context) error {
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetAndStore()
	return nil
}
