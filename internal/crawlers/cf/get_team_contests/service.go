package getteamcontests

import (
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"time"

	cf "github.com/laoin114514/codeforcesClient"
)

// GetCfTeamContests 获取CF团队比赛 - 对外提供的主要接口
func (g *GetTeamContests) GetCfTeamContests() error {
	g.count = 0

	// 带 manager 账号的签名请求 contest.list
	resp, err := g.client.WithHandle(g.useAccount).ContestList(&cf.ContestListParams{
		GroupCode: config.AppConfig.Cf.GroupCode,
	})
	if err != nil {
		g.log.Errorf("获取数据失败 %s", err.Error())
		return err
	}

	// 处理团队比赛数据
	err = g.processTeamContests(resp)
	if err != nil {
		g.log.Errorf("处理数据失败 %s", err.Error())
		return err
	}

	g.log.Printf("插入团队比赛 %d", g.count)
	return nil
}

// processTeamContests 处理团队比赛数据 - 私有方法
func (g *GetTeamContests) processTeamContests(resp *cf.ContestListResponse) error {
	for _, cfTeamContest := range resp.Result {
		if cfTeamContest == nil {
			continue
		}
		table := g.buildTeamContestTable(cfTeamContest)
		err := db.Insert_cf_team_contests(table)
		if err != nil {
			continue
		}
		g.count++
	}
	return nil
}

// buildTeamContestTable 构建团队比赛表 - 私有方法
func (g *GetTeamContests) buildTeamContestTable(cfTeamContest *cf.Contest) db.Cf_team_contests {
	startTime := cfTeamContest.StartTimeSeconds
	if startTime == 0 {
		startTime = time.Now().Unix()
	}
	return db.Cf_team_contests{
		Contest_id:   cfTeamContest.ID,
		Contest_name: cfTeamContest.Name,
		PrePare_by:   cfTeamContest.PreparedBy,
		Start_time:   time.Unix(startTime, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}
