package getuserrecords

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/internal/utils"
	cfurlgenerator "spider/pkg/cf-url-generator"
	easecrawler "spider/pkg/ease-crawler"
	"time"
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

	// 构建请求URL
	url, err := g.buildRequestURL(cfUserData.Account, true)
	if err != nil {
		g.log.Errorf("%s构建请求URL失败 %s", cfUserData.RealName, err.Error())
		return err
	}

	// 发起请求
	resp, err := g.fetchUserStatus(url)
	if err != nil {
		g.log.Errorf(fmt.Sprintf("%s请求数据失败，尝试无apiKey请求", cfUserData.RealName), err)
		url, err = g.buildRequestURL(cfUserData.Account, false)
		if err != nil {
			g.log.Errorf("%s构建请求URL失败 %s", cfUserData.RealName, err.Error())
			return err
		}
		resp, err = g.fetchUserStatus(url)
		if err != nil {
			g.log.Errorf("%s请求数据失败,URL:%s %s", cfUserData.RealName, url, err.Error())
			return err
		}
	}

	// 处理CF提交记录
	err = g.processSubmissionRecords(&cfUserData, &resp)
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

// buildRequestURL 构建请求URL - 私有方法
func (g *GetUserRecords) buildRequestURL(account string, isApiKey bool) (string, error) {
	// 尝试使用HTTPS
	url, err := cfurlgenerator.GenerateCFurlInstance.User.Status(
		isApiKey,
		&cfurlgenerator.UserStatusParams{
			Handle: account,
			From:   1,
			Count:  constants.CfMaxRecords,
		},
	)
	if err != nil {
		// 回退到HTTP
		url, err = cfurlgenerator.GenerateCFurlInstance.User.Status(
			false,
			&cfurlgenerator.UserStatusParams{
				Handle: account,
				From:   1,
				Count:  constants.CfMaxRecords,
			},
		)
	}
	return url, err
}

// fetchUserStatus 获取用户状态数据 - 私有方法
func (g *GetUserRecords) fetchUserStatus(url string) (models.CfUserStatusResponse, error) {
	req := utils.NewRequest[models.CfUserStatusResponse](true)
	return req.Get(url, map[string]string{})
}

// processSubmissionRecords 处理提交记录 - 私有方法
func (g *GetUserRecords) processSubmissionRecords(cfUserData *models.CfUserData, resp *models.CfUserStatusResponse) error {

	skippedCount := 0
	for _, cfRecord := range resp.Result {
		if g.shouldSkipRecord(cfRecord, cfUserData) {
			skippedCount++
			continue
		}

		table := g.buildSubmissionTable(&cfRecord, cfUserData)
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
func (g *GetUserRecords) shouldSkipRecord(cfRecord models.CfSubmission, cfUserData *models.CfUserData) bool {
	// 如果提交记录已存在，则跳过
	if cfUserData.OldDataSet[int(cfRecord.Id)] {
		return true
	}
	// 如果提交记录还在测试中，则跳过
	if cfRecord.Verdict == "TESTING" {
		return true
	}
	return false
}

// buildSubmissionTable 构建提交记录表 - 私有方法
func (g *GetUserRecords) buildSubmissionTable(cfRecord *models.CfSubmission, cfUserData *models.CfUserData) db.Cf_all_submissions {
	if cfRecord.Problem.Rating == 0 {
		cfRecord.Problem.Rating = -1
	}
	return db.Cf_all_submissions{
		Sub_id:        int(cfRecord.Id),
		Account:       cfUserData.Account,
		Problem_id:    fmt.Sprintf("%d%s", cfRecord.Problem.ContestId, cfRecord.Problem.Index),
		Problem_name:  cfRecord.Problem.Name,
		Verdict:       cfRecord.Verdict,
		Rating:        int(cfRecord.Problem.Rating),
		Creation_time: time.Unix(cfRecord.CreationTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}
