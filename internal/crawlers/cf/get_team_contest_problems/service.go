package getteamcontestproblems

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/internal/utils"
	cfurlgenerator "spider/pkg/cf-url-generator"
	easecrawler "spider/pkg/ease-crawler"
)

// GetCfTeamContestsProblems 获取CF团队比赛题目 - 对外提供的主要接口
func (g *GetTeamContestProblems) GetCfTeamContestsProblems() error {
	g.count = 0

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

	g.log.Printf("插入团队比赛题目 %d", g.count)
	return nil
}

// processTeamContestProblems 处理单个团队比赛的题目 - 私有方法
func (g *GetTeamContestProblems) processTeamContestProblems(teamContest db.Cf_team_contests) error {
	// 构建请求URL
	url, err := g.buildTeamContestProblemsURL(teamContest)
	if err != nil {
		g.log.Errorf("构建团队比赛题目请求URL失败 %s", err.Error())
		return err
	}

	// 获取题目数据
	problems, err := g.fetchTeamContestProblemsData(url)
	if err != nil {
		g.log.Errorf("获取团队比赛题目失败 %s", err.Error())
		return err
	}

	// 处理团队比赛题目
	err = g.processTeamContestProblemsData(&problems)
	if err != nil {
		g.log.Errorf("处理团队比赛题目失败 %s", err.Error())
		return err
	}
	g.log.Printf("团队比赛 %s 题目处理完成", teamContest.Contest_name)
	return nil
}

// buildTeamContestProblemsURL 构建团队比赛题目请求URL - 私有方法
func (g *GetTeamContestProblems) buildTeamContestProblemsURL(teamContest db.Cf_team_contests) (string, error) {
	return g.urlGenerator.Contest.Standings(teamContest.PrePare_by, &cfurlgenerator.ContestStandingsParams{
		ContestID:      teamContest.Contest_id,
		AsManager:      true,
		From:           1,
		Count:          constants.CfMaxRecords,
		ShowUnofficial: true,
	})
}

// fetchTeamContestProblemsData 获取团队比赛题目数据 - 私有方法
func (g *GetTeamContestProblems) fetchTeamContestProblemsData(url string) (models.CfTeamContestProblemsResponse, error) {
	req := utils.NewRequest[models.CfTeamContestProblemsResponse](true)
	return req.Get(url, map[string]string{})
}

// processTeamContestProblemsData 处理团队比赛题目数据 - 私有方法
func (g *GetTeamContestProblems) processTeamContestProblemsData(resp *models.CfTeamContestProblemsResponse) error {
	for _, cfProblem := range resp.Result.Problems {
		table := g.buildTeamProblemTable(&cfProblem)
		err := db.Insert_team_questions(table)
		if err != nil {
			continue
		}
		g.count++
	}
	return nil
}

// buildTeamProblemTable 构建团队题目表 - 私有方法
func (g *GetTeamContestProblems) buildTeamProblemTable(cfProblem *models.CfProblem) db.Cf_team_problems {
	if cfProblem.Rating == 0 {
		cfProblem.Rating = -1
	}
	return db.Cf_team_problems{
		Team_contest_id:   int(cfProblem.ContestId),
		Team_contest_name: cfProblem.Name,
		Problem_name:      cfProblem.Name,
		Rating:            int(cfProblem.Rating),
	}
}

// logError 记录错误日志 - 私有方法
func (g *GetTeamContestProblems) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	g.log.Errorf(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (g *GetTeamContestProblems) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	g.log.Printf(successMsg)
}
