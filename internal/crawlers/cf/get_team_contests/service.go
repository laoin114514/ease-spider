package getteamcontests

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/internal/utils"
	"time"
)

// GetCfTeamContests 获取CF团队比赛 - 对外提供的主要接口
func (g *GetTeamContests) GetCfTeamContests() error {
	g.count = 0

	// 构建请求URL
	url, err := g.buildTeamContestsURL()
	if err != nil {
		g.log.Errorf("构建请求URL失败 %s", err.Error())
		return err
	}

	// 获取数据
	resp, err := g.fetchTeamContestsData(url)
	if err != nil {
		g.log.Errorf("获取数据失败 %s", err.Error())
		return err
	}

	// 处理团队比赛数据
	err = g.processTeamContests(&resp)
	if err != nil {
		g.log.Errorf("处理数据失败 %s", err.Error())
		return err
	}

	g.log.Printf("插入团队比赛 %d", g.count)
	return nil
}

// buildTeamContestsURL 构建团队比赛请求URL - 私有方法
func (g *GetTeamContests) buildTeamContestsURL() (string, error) {
	GroupCode := config.AppConfig.Cf.GroupCode
	return utils.GenerateCFurlInstance.Contest.List(g.useAccount, &models.ContestListParams{
		GroupCode: GroupCode,
	})
}

// fetchTeamContestsData 获取团队比赛数据 - 私有方法
func (g *GetTeamContests) fetchTeamContestsData(url string) (models.CfTeamContestsResponse, error) {
	req := utils.NewRequest[models.CfTeamContestsResponse](true)
	return req.Get(url, map[string]string{})
}

// processTeamContests 处理团队比赛数据 - 私有方法
func (g *GetTeamContests) processTeamContests(resp *models.CfTeamContestsResponse) error {
	for _, cfTeamContest := range resp.Result {
		table := g.buildTeamContestTable(&cfTeamContest)
		err := db.Insert_cf_team_contests(table)
		if err != nil {
			continue
		}
		g.count++
	}
	return nil
}

// buildTeamContestTable 构建团队比赛表 - 私有方法
func (g *GetTeamContests) buildTeamContestTable(cfTeamContest *models.CfContest) db.Cf_team_contests {
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
func (g *GetTeamContests) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	g.log.Errorf(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (g *GetTeamContests) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	g.log.Printf(successMsg)
}
