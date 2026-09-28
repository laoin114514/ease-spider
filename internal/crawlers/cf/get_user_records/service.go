package getuserrecords

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"time"

	cf "github.com/laoin114514/codeforcesClient"
	easecrawler "spider/pkg/crawler"
)

// GetCfRecords 获取CF用户提交记录 - 对外提供的主要接口
func (g *GetUserRecords) GetCfRecords() error {
	g.log.Printf("开始获取CF用户提交记录")

	// 构建并发工具类
	conCurrenter := easecrawler.NewConCurrenter[models.CfUserData](config.AppConfig.Cf.CfRecordsConcurrency)
	conCurrenter.SetLogger(g.log)

	// 获取CF用户数据
	cfUserDatas, err := g.repo.GetCfAccountData()
	if err != nil {
		g.log.Errorf("获取CF用户数据失败 %s", err.Error())
		return err
	}
	g.log.Printf("成功获取到 %d 个CF用户", len(cfUserDatas))

	// 并发获取CF提交记录
	conCurrenter.Run(cfUserDatas, func(cfUserData models.CfUserData) error {
		return g.processUserRecords(cfUserData)
	})

	g.log.Printf("CF用户提交记录获取完成")
	return nil
}

// processUserRecords 处理单个用户的提交记录 - 私有方法
func (g *GetUserRecords) processUserRecords(cfUserData models.CfUserData) error {

	// 获取数据库中已有的提交记录
	err := g.loadExistingRecords(&cfUserData)
	if err != nil {
		g.log.Errorf("%s获取数据库中已有提交记录失败 %s", cfUserData.RealName, err.Error())
		return err
	}

	// 拉取提交记录
	resp, err := g.fetchUserStatus(cfUserData.Account)
	if err != nil {
		g.log.Errorf("%s请求数据失败 %s", cfUserData.RealName, err.Error())
		return err
	}

	// 处理CF提交记录
	err = g.processSubmissionRecords(&cfUserData, resp)
	if err != nil {
		g.log.Errorf("%s处理提交记录失败 %s", cfUserData.RealName, err.Error())
		return err
	}

	// 如果插入记录为0，则不进行日志记录
	if cfUserData.InsertCount == 0 {
		return nil
	}
	g.log.Printf("%s处理提交记录成功 %d", cfUserData.RealName, cfUserData.InsertCount)
	return nil
}

// loadExistingRecords 加载数据库中已有的提交记录 - 私有方法
func (g *GetUserRecords) loadExistingRecords(cfUserData *models.CfUserData) error {
	oldDataSet, err := g.repo.GetCfRecordsIdToset(cfUserData.Account)
	if err != nil {
		return err
	}
	cfUserData.OldDataSet = oldDataSet
	return nil
}

// fetchUserStatus 拉取某个账号的全部提交记录 - 私有方法
//
// 有凭据时优先带 apiKey 签名请求，失败再退回无签名的公共接口（与迁移前一致）。
func (g *GetUserRecords) fetchUserStatus(account string) (*cf.UserStatusResponse, error) {
	params := &cf.UserStatusParams{
		Handle: account,
		From:   1,
		Count:  constants.CfMaxRecords,
	}
	if _, ok := g.keys[account]; !ok {
		return g.plain.UserStatus(params)
	}
	resp, err := g.signed.WithHandle(account).UserStatus(params)
	if err == nil {
		return resp, nil
	}
	g.log.Warnf("%s带apiKey请求失败，改用无签名请求 %s", account, err.Error())
	return g.plain.UserStatus(params)
}

// processSubmissionRecords 处理提交记录 - 私有方法
func (g *GetUserRecords) processSubmissionRecords(cfUserData *models.CfUserData, resp *cf.UserStatusResponse) error {

	skippedCount := 0
	for _, cfRecord := range resp.Result {
		if g.shouldSkipRecord(cfRecord, cfUserData) {
			skippedCount++
			continue
		}

		table := g.buildSubmissionTable(cfRecord, cfUserData)
		err := db.Insert_cf_all_sub(table)
		if err != nil {
			g.log.Errorf("%s插入提交记录失败 %s", cfUserData.RealName, err.Error())
			return err
		}
		cfUserData.InsertCount++
	}

	g.log.Printf("用户 %s 提交记录处理完成，跳过 %d 条，新增 %d 条", cfUserData.RealName, skippedCount, cfUserData.InsertCount)
	return nil
}

// shouldSkipRecord 判断是否应该跳过该提交记录 - 私有方法
func (g *GetUserRecords) shouldSkipRecord(cfRecord *cf.Submission, cfUserData *models.CfUserData) bool {
	// 客户端用指针返回记录，这里必须防 nil，否则取 Problem 会 panic
	if cfRecord == nil || cfRecord.Problem == nil {
		g.log.Warnf("%s收到缺少题目信息的提交记录，已跳过", cfUserData.RealName)
		return true
	}
	// 如果提交记录已存在，则跳过
	if cfUserData.OldDataSet[int(cfRecord.ID)] {
		return true
	}
	// 如果提交记录还在测试中，则跳过
	if cfRecord.Verdict == "TESTING" {
		return true
	}
	return false
}

// buildSubmissionTable 构建提交记录表 - 私有方法
func (g *GetUserRecords) buildSubmissionTable(cfRecord *cf.Submission, cfUserData *models.CfUserData) db.Cf_all_submissions {
	rating := cfRecord.Problem.Rating
	if rating == 0 {
		rating = -1
	}
	return db.Cf_all_submissions{
		Sub_id:        int(cfRecord.ID),
		Account:       cfUserData.Account,
		Problem_id:    fmt.Sprintf("%d%s", cfRecord.Problem.ContestID, cfRecord.Problem.Index),
		Problem_name:  cfRecord.Problem.Name,
		Verdict:       cfRecord.Verdict,
		Rating:        rating,
		Creation_time: time.Unix(cfRecord.CreationTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}
