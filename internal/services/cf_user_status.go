package services

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/internal/repository"
	"spider/internal/utils"
	"time"
)

// ================================获取CF用户提交记录===============================================
type CfUserStatus struct {
	moduleDetail[models.CfUserStatusResponse]
}

// NewCfUserStatus 创建CF用户提交记录服务
func NewCfUserStatus() *CfUserStatus {
	return &CfUserStatus{
		moduleDetail: moduleDetail[models.CfUserStatusResponse]{
			LogService: NewLogService("logs/cfUserStatus.log", "logs/cfUserStatus.err.log"),
			repo:       repository.NewCfRepository(),
			debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		},
	}
}

// GetCfRecords 获取CF用户提交记录 - 对外提供的主要接口
func (r *CfUserStatus) GetCfRecords() error {
	r.debug.Debug("开始获取CF用户提交记录")

	// 构建并发工具类
	conCurrenter := utils.NewConCurrenter[models.CfUserData](config.AppConfig.Cf.CfRecordsConcurrency)

	// 获取CF用户数据
	cfUserDatas, err := r.repo.GetCfAccountData()
	if err != nil {
		r.logError("获取CF用户数据失败", err)
		return err
	}
	r.debug.Debug(fmt.Sprintf("成功获取到 %d 个CF用户", len(cfUserDatas)))

	// 并发获取CF提交记录
	conCurrenter.Run(cfUserDatas, func(cfUserData models.CfUserData) error {
		return r.processUserRecords(cfUserData)
	})

	r.debug.Debug("CF用户提交记录获取完成")
	return nil
}

// processUserRecords 处理单个用户的提交记录 - 私有方法
func (r *CfUserStatus) processUserRecords(cfUserData models.CfUserData) error {

	// 获取数据库中已有的提交记录
	err := r.loadExistingRecords(&cfUserData)
	if err != nil {
		r.logError(fmt.Sprintf("%s获取数据库中已有提交记录失败", cfUserData.RealName), err)
		return err
	}

	// 构建请求URL
	url, err := r.buildRequestURL(cfUserData.Account, true)
	if err != nil {
		r.logError(fmt.Sprintf("%s构建请求URL失败", cfUserData.RealName), err)
		return err
	}

	// 发起请求
	resp, err := r.fetchUserStatus(url)
	if err != nil {
		r.logError(fmt.Sprintf("%s请求数据失败，尝试无apiKey请求", cfUserData.RealName), err)
		url, err = r.buildRequestURL(cfUserData.Account, false)
		if err != nil {
			r.logError(fmt.Sprintf("%s构建请求URL失败", cfUserData.RealName), err)
			return err
		}
		resp, err = r.fetchUserStatus(url)
		if err != nil {
			r.logError(fmt.Sprintf("%s请求数据失败,URL:%s", cfUserData.RealName, url), err)
			return err
		}
	}

	// 处理CF提交记录
	err = r.processSubmissionRecords(&cfUserData, &resp)
	if err != nil {
		r.logError(fmt.Sprintf("%s处理提交记录失败", cfUserData.RealName), err)
		return err
	}

	// 如果插入记录为0，则不进行日志记录
	if cfUserData.InsertCount == 0 {
		return nil
	}
	r.logSuccess(fmt.Sprintf("%s处理提交记录成功", cfUserData.RealName), cfUserData.InsertCount)
	return nil
}

// loadExistingRecords 加载数据库中已有的提交记录 - 私有方法
func (r *CfUserStatus) loadExistingRecords(cfUserData *models.CfUserData) error {
	oldDataSet, err := r.repo.GetCfRecordsIdToset(cfUserData.Account)
	if err != nil {
		return err
	}
	cfUserData.OldDataSet = oldDataSet
	return nil
}

// buildRequestURL 构建请求URL - 私有方法
func (r *CfUserStatus) buildRequestURL(account string, isApiKey bool) (string, error) {
	// 尝试使用HTTPS
	url, err := utils.GenerateCFurlInstance.User.Status(
		isApiKey,
		&models.UserStatusParams{
			Handle: account,
			From:   1,
			Count:  constants.CfMaxRecords,
		},
	)
	if err != nil {
		// 回退到HTTP
		url, err = utils.GenerateCFurlInstance.User.Status(
			false,
			&models.UserStatusParams{
				Handle: account,
				From:   1,
				Count:  constants.CfMaxRecords,
			},
		)
	}
	return url, err
}

// fetchUserStatus 获取用户状态数据 - 私有方法
func (r *CfUserStatus) fetchUserStatus(url string) (models.CfUserStatusResponse, error) {
	req := utils.NewRequest[models.CfUserStatusResponse](true)
	return req.Get(url, map[string]string{})
}

// processSubmissionRecords 处理提交记录 - 私有方法
func (r *CfUserStatus) processSubmissionRecords(cfUserData *models.CfUserData, resp *models.CfUserStatusResponse) error {

	skippedCount := 0
	for _, cfRecord := range resp.Result {
		if r.shouldSkipRecord(cfRecord, cfUserData) {
			skippedCount++
			continue
		}

		table := r.buildSubmissionTable(&cfRecord, cfUserData)
		err := db.Insert_cf_all_sub(table)
		if err != nil {
			r.logError(fmt.Sprintf("%s插入提交记录失败", cfUserData.RealName), err)
			return err
		}
		cfUserData.InsertCount++
	}

	r.debug.Debug(fmt.Sprintf("用户 %s 提交记录处理完成，跳过 %d 条，新增 %d 条", cfUserData.RealName, skippedCount, cfUserData.InsertCount))
	return nil
}

// shouldSkipRecord 判断是否应该跳过该提交记录 - 私有方法
func (r *CfUserStatus) shouldSkipRecord(cfRecord models.CfSubmission, cfUserData *models.CfUserData) bool {
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
func (r *CfUserStatus) buildSubmissionTable(cfRecord *models.CfSubmission, cfUserData *models.CfUserData) db.Cf_all_submissions {
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

// logError 记录错误日志 - 私有方法
func (r *CfUserStatus) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *CfUserStatus) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}
