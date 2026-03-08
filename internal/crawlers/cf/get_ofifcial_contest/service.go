package getofifcialcontest

import (
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/internal/utils"
	"time"
)

// GetCfOfficialContests 获取CF官方比赛 - 对外提供的主要接口
func (g *GetOfifcialContest) GetCfOfficialContests() error {
	g.count = 0

	// 构建请求URL
	url, err := g.buildOfficialContestsURL()
	if err != nil {
		g.log.Errorf("构建请求URL失败 %s", err.Error())
		return err
	}

	// 获取数据
	resp, err := g.fetchOfficialContestsData(url)
	if err != nil {
		g.log.Errorf("获取数据失败 %s", err.Error())
		return err
	}

	// 处理官方比赛数据
	err = g.processOfficialContests(&resp)
	if err != nil {
		g.log.Errorf("处理数据失败 %s", err.Error())
		return err
	}

	g.log.Printf("插入官方比赛 %d", g.count)
	return nil
}

// buildOfficialContestsURL 构建官方比赛请求URL - 私有方法
func (g *GetOfifcialContest) buildOfficialContestsURL() (string, error) {
	return utils.GenerateCFurlInstance.Contest.List(config.AppConfig.Cf.ManagerAccount, &models.ContestListParams{
		GroupCode: config.AppConfig.Cf.GroupCode,
	})
}

// fetchOfficialContestsData 获取官方比赛数据 - 私有方法
func (g *GetOfifcialContest) fetchOfficialContestsData(url string) (models.CfOfficialContestsResponse, error) {
	req := utils.NewRequest[models.CfOfficialContestsResponse](true)
	return req.Get(url, map[string]string{})
}

// processOfficialContests 处理官方比赛数据 - 私有方法
func (g *GetOfifcialContest) processOfficialContests(resp *models.CfOfficialContestsResponse) error {
	for _, cfTeamContest := range resp.Result {
		table := g.buildOfficialContestTable(&cfTeamContest)
		err := db.Insert_cf_official_contests(table)
		if err != nil {
			continue
		}
		g.count++
	}
	return nil
}

// buildOfficialContestTable 构建官方比赛表 - 私有方法
func (g *GetOfifcialContest) buildOfficialContestTable(cfTeamContest *models.CfContest) db.Cf_official_contests {
	return db.Cf_official_contests{
		Official_contest_id:   int(cfTeamContest.Id),
		Official_contest_name: cfTeamContest.Name,
		Phase:                 cfTeamContest.Phase,
		Start_time:            time.Unix(cfTeamContest.StartTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}
