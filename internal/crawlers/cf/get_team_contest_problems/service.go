package getteamcontestproblems

import (
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"

	cf "github.com/laoin114514/codeforcesClient"
	easecrawler "github.com/laoin114514/ease-crawler"
)

// GetCfTeamContestsProblems 获取CF团队比赛题目 - 对外提供的主要接口
func (g *GetTeamContestProblems) GetCfTeamContestsProblems() error {
	g.count.Store(0)

	// 获取团队比赛数据
	teamContests, err := g.repo.GetTeamContests()
	if err != nil {
		g.log.Errorf("获取团队比赛失败 %s", err.Error())
		return err
	}

	// 构建并发工具类
	conCurrenter := easecrawler.NewConCurrenter[db.Cf_team_contests](config.AppConfig.Cf.CfTeamContestProblemsConcurrency)
	conCurrenter.SetLogger(g.log)
	// 并发处理团队比赛题目
	conCurrenter.Run(teamContests, func(teamContest db.Cf_team_contests) error {
		return g.processTeamContestProblems(teamContest)
	})

	g.log.Printf("插入团队比赛题目 %d", g.count.Load())
	return nil
}

// processTeamContestProblems 处理单个团队比赛的题目 - 私有方法
func (g *GetTeamContestProblems) processTeamContestProblems(teamContest db.Cf_team_contests) error {
	// 以出题人身份带签名请求 standings
	problems, err := g.client.WithHandle(teamContest.PrePare_by).ContestStandings(&cf.ContestStandingsParams{
		ContestID:      teamContest.Contest_id,
		AsManager:      true,
		From:           1,
		Count:          constants.CfMaxRecords,
		ShowUnofficial: true,
	})
	if err != nil {
		g.log.Errorf("获取团队比赛题目失败 %s", err.Error())
		return err
	}

	// 处理团队比赛题目
	err = g.processTeamContestProblemsData(problems)
	if err != nil {
		g.log.Errorf("处理团队比赛题目失败 %s", err.Error())
		return err
	}
	g.log.Printf("团队比赛 %s 题目处理完成", teamContest.Contest_name)
	return nil
}

// processTeamContestProblemsData 处理团队比赛题目数据 - 私有方法
func (g *GetTeamContestProblems) processTeamContestProblemsData(resp *cf.ContestStandingsResponse) error {
	if resp.Result == nil {
		return nil
	}
	for _, cfProblem := range resp.Result.Problems {
		if cfProblem == nil {
			continue
		}
		table := g.buildTeamProblemTable(cfProblem)
		err := db.Insert_team_questions(table)
		if err != nil {
			continue
		}
		g.count.Add(1)
	}
	return nil
}

// buildTeamProblemTable 构建团队题目表 - 私有方法
//
// 字段映射保持迁移前的行为不变：Team_contest_name 取的是题目标题，
// Official_contest_ID 没有赋值（沿用从前的写法，是否修正另行处理）。
func (g *GetTeamContestProblems) buildTeamProblemTable(cfProblem *cf.Problem) db.Cf_team_problems {
	rating := cfProblem.Rating
	if rating == 0 {
		rating = -1
	}
	return db.Cf_team_problems{
		Team_contest_id:   cfProblem.ContestID,
		Team_contest_name: cfProblem.Name,
		Problem_name:      cfProblem.Name,
		Rating:            rating,
	}
}
