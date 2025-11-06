package services

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"strings"
)

// ================================获取CF官方题目===============================================
type CfOfficialProblems struct {
	moduleDetail[models.CfOfficialProblemsResponse]
	count int
}

// NewCfOfficialProblems 创建CF官方题目服务
func NewCfOfficialProblems() *CfOfficialProblems {
	return &CfOfficialProblems{
		moduleDetail: moduleDetail[models.CfOfficialProblemsResponse]{
			LogService: NewLogService("logs/cfOfficialProblems.log", "logs/cfOfficialProblems.err.log"),
			repo:       repository.NewCfRepository(),
			debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		},
		count: 0,
	}
}

// GetCfOfficialProblems 获取CF官方题目 - 对外提供的主要接口
func (r *CfOfficialProblems) GetCfOfficialProblems() error {
	r.count = 0

	// 构建请求URL
	url, err := r.buildProblemsURL()
	if err != nil {
		r.logError("构建请求URL失败", err)
		return err
	}

	// 获取数据
	resp, err := r.fetchProblemsData(url)
	if err != nil {
		r.logError("获取数据失败", err)
		return err
	}

	// 处理官方题目数据
	err = r.processOfficialProblems(&resp)
	if err != nil {
		r.logError("处理数据失败", err)
		return err
	}

	r.logSuccess("插入官方题目", r.count)
	return nil
}

// buildProblemsURL 构建题目请求URL - 私有方法
func (r *CfOfficialProblems) buildProblemsURL() (string, error) {
	return utils.GenerateCFurlInstance.ProblemSet.Problems(&models.ProblemsetProblemsParams{})
}

// fetchProblemsData 获取题目数据 - 私有方法
func (r *CfOfficialProblems) fetchProblemsData(url string) (models.CfOfficialProblemsResponse, error) {
	req := utils.NewRequest[models.CfOfficialProblemsResponse](true)
	return req.Get(url, map[string]string{})
}

// processOfficialProblems 处理官方题目数据 - 私有方法
func (r *CfOfficialProblems) processOfficialProblems(resp *models.CfOfficialProblemsResponse) error {
	for _, cfProblem := range resp.Result.Problems {
		table := r.buildProblemTable(&cfProblem)
		err := db.Insert_cf_official(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}

// buildProblemTable 构建题目表 - 私有方法
func (r *CfOfficialProblems) buildProblemTable(cfProblem *models.CfProblem) db.Cf_official_problems {
	tags := strings.Join(cfProblem.Tags, "\",\"")
	tagsArr := fmt.Sprintf("[\"%s\"]", tags)
	if cfProblem.Rating == 0 {
		cfProblem.Rating = -1
	}
	return db.Cf_official_problems{
		Problem_id: fmt.Sprintf("%d%s", cfProblem.ContestId, cfProblem.Index),
		Title:      cfProblem.Name,
		Points:     int(cfProblem.Points),
		Rating:     int(cfProblem.Rating),
		Tags:       tagsArr,
	}
}

// logError 记录错误日志 - 私有方法
func (r *CfOfficialProblems) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *CfOfficialProblems) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}
