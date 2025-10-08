package services

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/src/constants"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"strings"
	"time"
)

type CfService struct {
	CfUserStatus           cfUserStatus
	CfOfficialProblems     cfOfficialProblems
	CfTeamContests         cfTeamContests
	CfTeamContestsProblems cfTeamContestsProblems
	CfOfficialContests     cfOfficialContests
}

// T类型是请求响应的结构体
type moduleDetail[T any] struct {
	*LogService
	repo *repository.CfRepository
}

func NewCfService() *CfService {
	return &CfService{
		CfUserStatus: cfUserStatus{
			moduleDetail: moduleDetail[models.CfUserStatusResponse]{
				LogService: NewLogService("logs/cfUserStatus.log", "logs/cfUserStatus.err.log"),
				repo:       repository.NewCfRepository(),
			},
		},
		CfOfficialProblems: cfOfficialProblems{
			moduleDetail: moduleDetail[models.CfOfficialProblemsResponse]{
				LogService: NewLogService("logs/cfOfficialProblems.log", "logs/cfOfficialProblems.err.log"),
				repo:       repository.NewCfRepository(),
			},
			count: 0,
		},
		CfTeamContests: cfTeamContests{
			moduleDetail: moduleDetail[models.CfTeamContestsResponse]{
				LogService: NewLogService("logs/cfTeamContests.log", "logs/cfTeamContests.err.log"),
				repo:       repository.NewCfRepository(),
			},
			count:      0,
			useAccount: config.AppConfig.Cf.ManagerAccount,
		},
		CfTeamContestsProblems: cfTeamContestsProblems{
			moduleDetail: moduleDetail[models.CfTeamContestProblemsResponse]{
				LogService: NewLogService("logs/cfTeamContestsProblems.log", "logs/cfTeamContestsProblems.err.log"),
				repo:       repository.NewCfRepository(),
			},
			count: 0,
		},
		CfOfficialContests: cfOfficialContests{
			moduleDetail: moduleDetail[models.CfOfficialContestsResponse]{
				LogService: NewLogService("logs/cfOfficialContests.log", "logs/cfOfficialContests.err.log"),
				repo:       repository.NewCfRepository(),
			},
			count: 0,
		},
	}
}

// //////////////////////////////////////////////// 获取cf用户提交记录////////////////////////////////////////////////////////
type cfUserStatus struct {
	moduleDetail[models.CfUserStatusResponse]
}

// 获取cf用户提交记录
func (r *cfUserStatus) GetCfRecords() error {
	//构建工具类
	conCurrenter := utils.NewConCurrenter[models.CfUserData](config.AppConfig.Cf.CfRecordsConcurrency)

	//获取cf用户数据
	cfUserDatas, err := r.repo.GetCfAccountData()
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取cf用户数据失败 %v", err.Error()))
		return err
	}
	//并发获取cf提交记录
	conCurrenter.Run(cfUserDatas, func(cfUserData models.CfUserData) error {
		//获取db中已有的提交记录
		var err error
		cfUserData.OldDataSet, err = r.repo.GetCfRecordsIdToset(cfUserData.Account)
		if err != nil {
			r.AddErr(fmt.Sprintf("%s 获取db中已有的提交记录失败 %v", cfUserData.RealName, err))
			return err
		}
		//按照cf规则拼接url
		url, err := utils.GenerateCFurlInstance.User.Status(
			true,
			&models.UserStatusParams{
				Handle: cfUserData.Account,
				From:   1,
				Count:  constants.CfMaxRecords,
			},
		)
		if err != nil {
			r.AddErr(fmt.Sprintf("%s 获取url失败,尝试使用无apikey", cfUserData.RealName))
			url, err = utils.GenerateCFurlInstance.User.Status(
				false,
				&models.UserStatusParams{
					Handle: cfUserData.Account,
					From:   1,
					Count:  constants.CfMaxRecords,
				},
			)
			if err != nil {
				r.AddErr(fmt.Sprintf("%s 获取url失败 %v", cfUserData.RealName, err))
				return err
			}
			r.AddErr(fmt.Sprintf("%s 获取url失败,尝试使用无apikey成功", cfUserData.RealName))
		}

		//发起请求
		req := utils.NewRequest[models.CfUserStatusResponse]()
		resp, err := req.Get(url, map[string]string{})
		if err != nil {
			r.AddErr(fmt.Sprintf("%s 请求数据失败 %v", cfUserData.RealName, err))
			return err
		}

		// 处理cf提交记录
		err = r.handleCfRecords(&cfUserData, &resp)
		if err != nil {
			r.AddErr(fmt.Sprintf("%s 处理提交记录失败 %v", cfUserData.RealName, err))
			return err
		}

		r.AddLog(fmt.Sprintf("%s 处理提交记录成功 %d", cfUserData.RealName, cfUserData.InsertCount))
		return nil
	})
	return nil
}

// 处理cf提交记录
func (r *cfUserStatus) handleCfRecords(cfUserData *models.CfUserData, resp *models.CfUserStatusResponse) error {
	for _, cfRecord := range resp.Result {
		// 如果提交记录已存在，则跳过
		if cfUserData.OldDataSet[int(cfRecord.Id)] {
			continue
		}
		// 构建提交记录表
		table := r.buildTable(&cfRecord, cfUserData)
		// 插入提交记录
		err := db.Insert_cf_all_sub(table)
		if err != nil {
			r.log.AddErr(fmt.Sprintf("%s 插入提交记录失败 %v", cfUserData.RealName, err))
			return err
		}
		cfUserData.InsertCount++
	}
	return nil
}

// 构建提交记录表
func (r *cfUserStatus) buildTable(cfRecord *models.CfSubmission, cfUserData *models.CfUserData) db.Cf_all_submissions {
	if cfRecord.Problem.Rating == 0 {
		cfRecord.Problem.Rating = -1
	}
	return db.Cf_all_submissions{
		Sub_id:        int(cfRecord.Id),
		Account:       cfUserData.Account,
		Problem_id:    fmt.Sprintf("%d%s", cfRecord.Problem.ContestId, cfRecord.Problem.Index),
		Problem_name:  cfRecord.Problem.Name,
		Verdict:       cfRecord.Verdict,
		Rating:        int(cfRecord.Problem.Rating),
		Creation_time: time.Unix(cfRecord.CreationTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}

// //////////////////////////////////////////////// 获取cf官方题目////////////////////////////////////////////////////////
type cfOfficialProblems struct {
	moduleDetail[models.CfOfficialProblemsResponse]
	count int
}

func (r *cfOfficialProblems) GetCfOfficialProblems() error {
	r.count = 0
	url, err := utils.GenerateCFurlInstance.ProblemSet.Problems(&models.ProblemsetProblemsParams{})
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取url失败 %v", err))
		return err
	}
	req := utils.NewRequest[models.CfOfficialProblemsResponse]()
	resp, err := req.Get(url, map[string]string{})
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取数据失败 %v", err))
		return err
	}
	err = r.handleCfOfficialProblems(&resp)
	if err != nil {
		r.log.AddErr(fmt.Sprintf("处理数据失败 %v", err))
		return err
	}
	r.log.AddLog(fmt.Sprintf("插入%d条官方题目", r.count))
	return nil
}
func (r *cfOfficialProblems) handleCfOfficialProblems(resp *models.CfOfficialProblemsResponse) error {
	for _, cfProblem := range resp.Result.Problems {
		table := r.buildTable(&cfProblem)
		err := db.Insert_cf_official(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}
func (r *cfOfficialProblems) buildTable(cfProblem *models.CfProblem) db.Cf_official_problems {
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

// //////////////////////////////////////////////// 获取cf团队比赛////////////////////////////////////////////////////////
type cfTeamContests struct {
	moduleDetail[models.CfTeamContestsResponse]
	useAccount string
	count      int
}

func (r *cfTeamContests) GetCfTeamContests() error {
	r.count = 0
	GroupCode := config.AppConfig.Cf.GroupCode
	url, err := utils.GenerateCFurlInstance.Contest.List(r.useAccount, &models.ContestListParams{
		GroupCode: GroupCode,
	})
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取url失败 %v", err))
		return err
	}
	req := utils.NewRequest[models.CfTeamContestsResponse]()
	resp, err := req.Get(url, map[string]string{})
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取数据失败 %v", err))
		return err
	}
	err = r.handleCfTeamContests(&resp)
	if err != nil {
		r.log.AddErr(fmt.Sprintf("处理数据失败 %v", err))
		return err
	}
	r.log.AddLog(fmt.Sprintf("插入%d条团队比赛", r.count))
	return nil
}

func (r *cfTeamContests) handleCfTeamContests(resp *models.CfTeamContestsResponse) error {
	for _, cfTeamContest := range resp.Result {
		table := r.buildTable(&cfTeamContest)
		err := db.Insert_cf_team_contests(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}
func (r *cfTeamContests) buildTable(cfTeamContest *models.CfContest) db.Cf_team_contests {
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

// //////////////////////////////////////////////// 获取cf团队比赛题目////////////////////////////////////////////////////////
type cfTeamContestsProblems struct {
	moduleDetail[models.CfTeamContestProblemsResponse]
	count int
}

func (r *cfTeamContestsProblems) GetCfTeamContestsProblems() error {
	r.count = 0
	teamContests, err := r.repo.GetTeamContests()
	conCurrenter := utils.NewConCurrenter[db.Cf_team_contests](config.AppConfig.Cf.CfTeamContestsConcurrency)
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取团队比赛失败 %v", err))
		return err
	}
	conCurrenter.Run(teamContests, func(teamContest db.Cf_team_contests) error {
		var url string
		url, err = utils.GenerateCFurlInstance.Contest.Standings(teamContest.PrePare_by, &models.ContestStandingsParams{
			ContestID:      teamContest.Contest_id,
			AsManager:      true,
			From:           1,
			Count:          constants.CfMaxRecords,
			ShowUnofficial: true,
		})
		if err != nil {
			r.log.AddErr(fmt.Sprintf("获取团队比赛题目失败 %v", err))
			return err
		}
		req := utils.NewRequest[models.CfTeamContestProblemsResponse]()
		problems, err := req.Get(url, map[string]string{})
		if err != nil {
			r.log.AddErr(fmt.Sprintf("获取团队比赛题目失败 %v", err))
			return err
		}
		err = r.handleCfTeamContestsProblems(&problems)
		if err != nil {
			r.log.AddErr(fmt.Sprintf("处理团队比赛题目失败 %v", err))
			return err
		}
		return nil
	})
	r.log.AddLog(fmt.Sprintf("插入%d条团队比赛题目", r.count))
	return nil
}
func (r *cfTeamContestsProblems) handleCfTeamContestsProblems(resp *models.CfTeamContestProblemsResponse) error {
	for _, cfProblem := range resp.Result.Problems {
		table := r.buildTable(&cfProblem)
		err := db.Insert_team_questions(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}
func (r *cfTeamContestsProblems) buildTable(cfProblem *models.CfProblem) db.Cf_team_problems {
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

// ////////////////////////////////////////////// 获取cf官方比赛////////////////////////////////////////////////////////
type cfOfficialContests struct {
	moduleDetail[models.CfOfficialContestsResponse]
	count int
}

func (r *cfOfficialContests) GetCfOfficialContests() error {
	r.count = 0
	GroupCode := config.AppConfig.Cf.GroupCode
	url, err := utils.GenerateCFurlInstance.Contest.List("", &models.ContestListParams{
		GroupCode: GroupCode,
	})
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取url失败 %v", err))
		return err
	}
	req := utils.NewRequest[models.CfOfficialContestsResponse]()
	resp, err := req.Get(url, map[string]string{})
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取数据失败 %v", err))
		return err
	}
	err = r.handleCfOfficialContests(&resp)
	if err != nil {
		r.log.AddErr(fmt.Sprintf("处理数据失败 %v", err))
		return err
	}
	r.log.AddLog(fmt.Sprintf("插入%d条官方比赛", r.count))
	return nil
}

func (r *cfOfficialContests) handleCfOfficialContests(resp *models.CfOfficialContestsResponse) error {
	for _, cfTeamContest := range resp.Result {
		table := r.buildTable(&cfTeamContest)
		err := db.Insert_cf_official_contests(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}
func (r *cfOfficialContests) buildTable(cfTeamContest *models.CfContest) db.Cf_official_contests {
	return db.Cf_official_contests{
		Official_contest_id:   int(cfTeamContest.Id),
		Official_contest_name: cfTeamContest.Name,
		Phase:                 cfTeamContest.Phase,
		Start_time:            time.Unix(cfTeamContest.StartTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}
