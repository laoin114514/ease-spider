package getuserrecords

import (
	"spider/config"
	"spider/internal/repository"
	"spider/pkg/luogu2api"
	"time"

	easecrawler "spider/pkg/crawler"
)

// GetUserRecords 获取洛谷用户提交记录。
//
// 提交记录统一走 Luogu2Api 服务（pkg/luogu2api SDK），账号与登录态由服务端的号池维护，
// 插件自身不再登录洛谷、也不保存 cookie。
type GetUserRecords struct {
	log  *easecrawler.EaseLogger
	repo *repository.LuoguRepository
	// sdk 懒加载后复用；SDK 无可变状态，可并发使用（见 pkg/luogu2api/doc.go）
	sdk *luogu2api.Client
}

func NewGetUserRecords() *GetUserRecords {
	return &GetUserRecords{
		repo: repository.NewLuoguRepository(),
	}
}
func (g *GetUserRecords) Name() string {
	return "get_user_records"
}
func (g *GetUserRecords) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:         config.IntervalOr(config.TimerFrequency().LuoguRecords, 5*time.Hour),
		StartImmediately: true,
	}
}
func (g *GetUserRecords) Run(c *easecrawler.Context) error {
	g.log = easecrawler.GetCrawlerLogger(c)
	g.log.Println("开始获取用户记录")

	// 两个步骤互不依赖：取记录失败时，历史数据的私有题难度仍然照常修正
	recordsErr := g.GetAndStore()
	privateErr := g.ChangePrivateProblem()
	if recordsErr != nil {
		return recordsErr
	}
	return privateErr
}
