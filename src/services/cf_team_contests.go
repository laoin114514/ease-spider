package services

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/src/constants"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"time"
)

// ================================获取CF团队比赛===============================================
type CfTeamContests struct {
	moduleDetail[models.CfTeamContestsResponse]
	useAccount string
	count      int
}

// NewCfTeamContests 创建CF团队比赛服务
func NewCfTeamContests() *CfTeamContests {
	return &CfTeamContests{
		moduleDetail: moduleDetail[models.CfTeamContestsResponse]{
			LogService: NewLogService("logs/cfTeamContests.log", "logs/cfTeamContests.err.log"),
			repo:       repository.NewCfRepository(),
			debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		},
		count:      0,
		useAccount: config.AppConfig.Cf.ManagerAccount,
	}
}

// GetCfTeamContests 获取CF团队比赛 - 对外提供的主要接口
func (r *CfTeamContests) GetCfTeamContests() error {
	r.count = 0

	// 构建请求URL
	url, err := r.buildTeamContestsURL()
	if err != nil {
		r.logError("构建请求URL失败", err)
		return err
	}

	// 获取数据
	resp, err := r.fetchTeamContestsData(url)
	if err != nil {
		r.logError("获取数据失败", err)
		return err
	}

	// 处理团队比赛数据
	err = r.processTeamContests(&resp)
	if err != nil {
		r.logError("处理数据失败", err)
		return err
	}

	r.logSuccess("插入团队比赛", r.count)
	return nil
}

// buildTeamContestsURL 构建团队比赛请求URL - 私有方法
func (r *CfTeamContests) buildTeamContestsURL() (string, error) {
	GroupCode := config.AppConfig.Cf.GroupCode
	return utils.GenerateCFurlInstance.Contest.List(r.useAccount, &models.ContestListParams{
		GroupCode: GroupCode,
	})
}

// fetchTeamContestsData 获取团队比赛数据 - 私有方法
func (r *CfTeamContests) fetchTeamContestsData(url string) (models.CfTeamContestsResponse, error) {
	req := utils.NewRequest[models.CfTeamContestsResponse](true)
	return req.Get(url, map[string]string{})
}

// processTeamContests 处理团队比赛数据 - 私有方法
func (r *CfTeamContests) processTeamContests(resp *models.CfTeamContestsResponse) error {
	for _, cfTeamContest := range resp.Result {
		table := r.buildTeamContestTable(&cfTeamContest)
		err := db.Insert_cf_team_contests(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}

// buildTeamContestTable 构建团队比赛表 - 私有方法
func (r *CfTeamContests) buildTeamContestTable(cfTeamContest *models.CfContest) db.Cf_team_contests {
	if cfTeamContest.StartTimeSeconds == 0 {
		cfTeamContest.StartTimeSeconds = time.Now().Unix()
	}
	return db.Cf_team_contests{
		Contest_id:   int(cfTeamContest.Id),
		Contest_name: cfTeamContest.Name,
		PrePare_by:   cfTeamContest.PreparedBy,
		Start_time:   time.Unix(cfTeamContest.StartTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}

// logError 记录错误日志 - 私有方法
func (r *CfTeamContests) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *CfTeamContests) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}
