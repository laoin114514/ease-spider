package getuserrecords

import (
	"spider/config"
	"spider/internal/crawlers/cf/cfclient"
	"spider/internal/repository"
	"time"

	cf "github.com/laoin114514/codeforcesClient"
	easecrawler "github.com/laoin114514/ease-crawler"
)

type GetUserRecords struct {
	log  *easecrawler.EaseLogger
	repo *repository.CfRepository
	// signed 带号池签名（用 WithHandle 指定账号），plain 无签名作为兜底
	signed *cf.Client
	plain  *cf.Client
	// keys 有凭据的账号：没有凭据的直接走 plain，省一次必然失败的签名请求
	keys map[string]repository.CfApiKey
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
	keys, err := g.repo.GetCfApikeyPool()
	if err != nil {
		return err
	}
	// 限流取任务并发数：客户端自带限流与 429/5xx 重试，替代原先的裸请求
	g.keys = keys
	g.signed = cfclient.NewSigned(keys, config.AppConfig.Cf.CfRecordsConcurrency)
	g.plain = cfclient.NewPlain(config.AppConfig.Cf.CfRecordsConcurrency)
	g.log = easecrawler.GetCrawlerLogger(c)
	g.GetCfRecords()
	return nil
}
