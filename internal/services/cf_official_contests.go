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

// ================================获取CF官方比赛===============================================
type CfOfficialContests struct {
	moduleDetail[models.CfOfficialContestsResponse]
	count int
}

// NewCfOfficialContests 创建CF官方比赛服务
func NewCfOfficialContests() *CfOfficialContests {
	return &CfOfficialContests{
		moduleDetail: moduleDetail[models.CfOfficialContestsResponse]{
			LogService: NewLogService("logs/cfOfficialContests.log", "logs/cfOfficialContests.err.log"),
			repo:       repository.NewCfRepository(),
			debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		},
		count: 0,
	}
}

// GetCfOfficialContests 获取CF官方比赛 - 对外提供的主要接口
func (r *CfOfficialContests) GetCfOfficialContests() error {
	r.count = 0

	// 构建请求URL
	url, err := r.buildOfficialContestsURL()
	if err != nil {
		r.logError("构建请求URL失败", err)
		return err
	}

	// 获取数据
	resp, err := r.fetchOfficialContestsData(url)
	if err != nil {
		r.logError("获取数据失败", err)
		return err
	}

	// 处理官方比赛数据
	err = r.processOfficialContests(&resp)
	if err != nil {
		r.logError("处理数据失败", err)
		return err
	}

	r.logSuccess("插入官方比赛", r.count)
	return nil
}

// buildOfficialContestsURL 构建官方比赛请求URL - 私有方法
func (r *CfOfficialContests) buildOfficialContestsURL() (string, error) {
	return utils.GenerateCFurlInstance.Contest.List(config.AppConfig.Cf.ManagerAccount, &models.ContestListParams{
		GroupCode: config.AppConfig.Cf.GroupCode,
	})
}

// fetchOfficialContestsData 获取官方比赛数据 - 私有方法
func (r *CfOfficialContests) fetchOfficialContestsData(url string) (models.CfOfficialContestsResponse, error) {
	req := utils.NewRequest[models.CfOfficialContestsResponse](true)
	return req.Get(url, map[string]string{})
}

// processOfficialContests 处理官方比赛数据 - 私有方法
func (r *CfOfficialContests) processOfficialContests(resp *models.CfOfficialContestsResponse) error {
	for _, cfTeamContest := range resp.Result {
		table := r.buildOfficialContestTable(&cfTeamContest)
		err := db.Insert_cf_official_contests(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}

// buildOfficialContestTable 构建官方比赛表 - 私有方法
func (r *CfOfficialContests) buildOfficialContestTable(cfTeamContest *models.CfContest) db.Cf_official_contests {
	return db.Cf_official_contests{
		Official_contest_id:   int(cfTeamContest.Id),
		Official_contest_name: cfTeamContest.Name,
		Phase:                 cfTeamContest.Phase,
		Start_time:            time.Unix(cfTeamContest.StartTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}

// logError 记录错误日志 - 私有方法
func (r *CfOfficialContests) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *CfOfficialContests) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}
