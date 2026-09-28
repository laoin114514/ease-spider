package getuserrecords

import (
	"spider/config"
	"spider/internal/repository"
	"time"

	easecrawler "spider/pkg/crawler"
)

// GetUserRecords 牛客用户记录插件
type GetUserRecords struct {
	repo *repository.NiukeRepository
}

// NewGetUserRecords 创建插件实例
func NewGetUserRecords() *GetUserRecords {
	return &GetUserRecords{
		repo: repository.NewNiukeRepository(),
	}
}

// Name 返回插件名称
func (g *GetUserRecords) Name() string {
	return "get_user_records"
}

// Meta 返回插件元数据
func (g *GetUserRecords) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         config.IntervalOr(config.TimerFrequency().NiukeRecords, 1*time.Hour),
		StartImmediately: true,
	}
}

// Run 执行插件
func (g *GetUserRecords) Run(c *easecrawler.Context) error {
	log := easecrawler.GetCrawlerLogger(c)
	service := NewService(log, g.repo)
	return service.GetAndStore()
}
