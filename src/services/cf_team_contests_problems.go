package services

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/src/constants"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
)

// ================================获取CF团队比赛题目===============================================
type CfTeamContestsProblems struct {
	moduleDetail[models.CfTeamContestProblemsResponse]
	count int
}

// NewCfTeamContestsProblems 创建CF团队比赛题目服务
func NewCfTeamContestsProblems() *CfTeamContestsProblems {
	return &CfTeamContestsProblems{
		moduleDetail: moduleDetail[models.CfTeamContestProblemsResponse]{
			LogService: NewLogService("logs/cfTeamContestsProblems.log", "logs/cfTeamContestsProblems.err.log"),
			repo:       repository.NewCfRepository(),
			debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		},
		count: 0,
	}
}

// GetCfTeamContestsProblems 获取CF团队比赛题目 - 对外提供的主要接口
func (r *CfTeamContestsProblems) GetCfTeamContestsProblems() error {
	r.count = 0

	// 获取团队比赛数据
	teamContests, err := r.repo.GetTeamContests()
	if err != nil {
		r.logError("获取团队比赛失败", err)
		return err
	}

	// 构建并发工具类
	conCurrenter := utils.NewConCurrenter[db.Cf_team_contests](config.AppConfig.Cf.CfTeamContestProblemsConcurrency)

	// 并发处理团队比赛题目
	conCurrenter.Run(teamContests, func(teamContest db.Cf_team_contests) error {
		return r.processTeamContestProblems(teamContest)
	})

	r.logSuccess("插入团队比赛题目", r.count)
	return nil
}

// processTeamContestProblems 处理单个团队比赛的题目 - 私有方法
func (r *CfTeamContestsProblems) processTeamContestProblems(teamContest db.Cf_team_contests) error {
	// 构建请求URL
	url, err := r.buildTeamContestProblemsURL(teamContest)
	if err != nil {
		r.logError("构建团队比赛题目请求URL失败", err)
		return err
	}

	// 获取题目数据
	problems, err := r.fetchTeamContestProblemsData(url)
	if err != nil {
		r.logError("获取团队比赛题目失败", err)
		return err
	}

	// 处理团队比赛题目
	err = r.processTeamContestProblemsData(&problems)
	if err != nil {
		r.logError("处理团队比赛题目失败", err)
		return err
	}
	r.debug.Debug(fmt.Sprintf("团队比赛 %s 题目处理完成", teamContest.Contest_name))
	return nil
}

// buildTeamContestProblemsURL 构建团队比赛题目请求URL - 私有方法
func (r *CfTeamContestsProblems) buildTeamContestProblemsURL(teamContest db.Cf_team_contests) (string, error) {
	return utils.GenerateCFurlInstance.Contest.Standings(teamContest.PrePare_by, &models.ContestStandingsParams{
		ContestID:      teamContest.Contest_id,
		AsManager:      true,
		From:           1,
		Count:          constants.CfMaxRecords,
		ShowUnofficial: true,
	})
}

// fetchTeamContestProblemsData 获取团队比赛题目数据 - 私有方法
func (r *CfTeamContestsProblems) fetchTeamContestProblemsData(url string) (models.CfTeamContestProblemsResponse, error) {
	req := utils.NewRequest[models.CfTeamContestProblemsResponse](true)
	return req.Get(url, map[string]string{})
}

// processTeamContestProblemsData 处理团队比赛题目数据 - 私有方法
func (r *CfTeamContestsProblems) processTeamContestProblemsData(resp *models.CfTeamContestProblemsResponse) error {
	for _, cfProblem := range resp.Result.Problems {
		table := r.buildTeamProblemTable(&cfProblem)
		err := db.Insert_team_questions(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}

// buildTeamProblemTable 构建团队题目表 - 私有方法
func (r *CfTeamContestsProblems) buildTeamProblemTable(cfProblem *models.CfProblem) db.Cf_team_problems {
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
func (r *CfTeamContestsProblems) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *CfTeamContestsProblems) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}
