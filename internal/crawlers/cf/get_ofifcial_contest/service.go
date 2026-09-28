package getofifcialcontest

import (
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"time"

	cf "github.com/laoin114514/codeforcesClient"
)

// GetCfOfficialContests 获取CF官方比赛 - 对外提供的主要接口
func (g *GetOfifcialContest) GetCfOfficialContests() error {
	g.count = 0

	// 带 manager 账号的签名请求 contest.list
	resp, err := g.client.WithHandle(config.AppConfig.Cf.ManagerAccount).ContestList(&cf.ContestListParams{
		GroupCode: config.AppConfig.Cf.GroupCode,
	})
	if err != nil {
		g.log.Errorf("获取数据失败 %s", err.Error())
		return err
	}

	// 处理官方比赛数据
	err = g.processOfficialContests(resp)
	if err != nil {
		g.log.Errorf("处理数据失败 %s", err.Error())
		return err
	}

	g.log.Printf("插入官方比赛 %d", g.count)
	return nil
}

// processOfficialContests 处理官方比赛数据 - 私有方法
func (g *GetOfifcialContest) processOfficialContests(resp *cf.ContestListResponse) error {
	for _, cfTeamContest := range resp.Result {
		if cfTeamContest == nil {
			continue
		}
		table := g.buildOfficialContestTable(cfTeamContest)
		err := db.Insert_cf_official_contests(table)
		if err != nil {
			continue
		}
		g.count++
	}
	return nil
}

// buildOfficialContestTable 构建官方比赛表 - 私有方法
func (g *GetOfifcialContest) buildOfficialContestTable(cfTeamContest *cf.Contest) db.Cf_official_contests {
	return db.Cf_official_contests{
		Official_contest_id:   cfTeamContest.ID,
		Official_contest_name: cfTeamContest.Name,
		Phase:                 cfTeamContest.Phase,
		Start_time:            time.Unix(cfTeamContest.StartTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}
