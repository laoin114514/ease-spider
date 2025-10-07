package services

import (
	"fmt"
	"spider/config/db"
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
	repo *repository.CfRepository
	req  *utils.Request[T]
	log  *utils.LogContainer
}

func NewCfService() *CfService {
	return &CfService{
		CfUserStatus: cfUserStatus{
			moduleDetail: moduleDetail[models.CfUserStatusResponse]{
				repo: repository.NewCfRepository(),
				req:  utils.NewRequest[models.CfUserStatusResponse](),
				log:  utils.NewLogContainer(),
			},
		},
		CfOfficialProblems: cfOfficialProblems{
			moduleDetail: moduleDetail[models.CfOfficialProblemsResponse]{
				repo: repository.NewCfRepository(),
				req:  utils.NewRequest[models.CfOfficialProblemsResponse](),
				log:  utils.NewLogContainer(),
			},
			count: 0,
		},
		CfTeamContests: cfTeamContests{
			moduleDetail: moduleDetail[models.CfTeamContestsResponse]{
				repo: repository.NewCfRepository(),
				req:  utils.NewRequest[models.CfTeamContestsResponse](),
				log:  utils.NewLogContainer(),
			},
			count:      0,
			useAccount: "233zhang",
		},
		CfTeamContestsProblems: cfTeamContestsProblems{
			moduleDetail: moduleDetail[models.CfTeamContestProblemsResponse]{
				repo: repository.NewCfRepository(),
				req:  utils.NewRequest[models.CfTeamContestProblemsResponse](),
				log:  utils.NewLogContainer(),
			},
			count: 0,
		},
		CfOfficialContests: cfOfficialContests{
			moduleDetail: moduleDetail[models.CfOfficialContestsResponse]{
				repo: repository.NewCfRepository(),
				req:  utils.NewRequest[models.CfOfficialContestsResponse](),
				log:  utils.NewLogContainer(),
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
func (r *cfUserStatus) GetCfRecords(concurrency int) error {
	//构建工具类
	conCurrenter := utils.NewConCurrenter[models.CfUserData](concurrency)

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
			r.log.AddErr(fmt.Sprintf("%s 获取db中已有的提交记录失败 %v", cfUserData.RealName, err))
			return err
		}
		//按照cf规则拼接url
		url, err := utils.GenerateCFurlInstance.User.Status(&models.UserStatusParams{
			Handle: cfUserData.Account,
			From:   1,
			Count:  50000,
		})
		if err != nil {
			r.log.AddErr(fmt.Sprintf("%s 获取url失败 %v", cfUserData.RealName, err))
			return err
		}

		//发起请求
		resp, err := r.req.Get(url, map[string]string{})
		if err != nil {
			r.log.AddErr(fmt.Sprintf("%s 请求数据失败 %v", cfUserData.RealName, err))
			return err
		}

		// 处理cf提交记录
		err = r.handleCfRecords(&cfUserData, &resp)
		if err != nil {
			r.log.AddErr(fmt.Sprintf("%s 处理提交记录失败 %v", cfUserData.RealName, err))
			return err
		}

		r.log.AddLog(fmt.Sprintf("%s 处理提交记录成功 %d", cfUserData.RealName, cfUserData.InsertCount))
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
	return db.Cf_all_submissions{
		Sub_id:        int(cfRecord.Id),
		Account:       cfUserData.Account,
		Problem_id:    fmt.Sprintf("%d%s", cfRecord.Problem.ContestId, cfRecord.Problem.Index),
		Problem_name:  cfRecord.Problem.Name,
		Verdict:       cfRecord.Verdict,
		Rating:        int(cfRecord.Problem.Rating),
		Creation_time: time.Unix(cfRecord.CreationTimeSeconds, 0).Add(8 * time.Hour),
	}
}
func (r *cfUserStatus) GetLog() []string {
	return r.log.GetLog()
}
func (r *cfUserStatus) GetErr() []string {
	return r.log.GetErr()
}
func (r *cfUserStatus) Clear() error {
	r.log.ClearLog()
	r.log.ClearErr()
	return nil
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
	resp, err := r.req.Get(url, map[string]string{})
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
	return db.Cf_official_problems{
		Problem_id: fmt.Sprintf("%d%s", cfProblem.ContestId, cfProblem.Index),
		Title:      cfProblem.Name,
		Points:     int(cfProblem.Points),
		Rating:     int(cfProblem.Rating),
		Tags:       tagsArr,
	}
}
func (r *cfOfficialProblems) GetLog() []string {
	return r.log.GetLog()
}
func (r *cfOfficialProblems) GetErr() []string {
	return r.log.GetErr()
}
func (r *cfOfficialProblems) Clear() error {
	r.log.ClearLog()
	r.log.ClearErr()
	r.count = 0
	return nil
}

// //////////////////////////////////////////////// 获取cf团队比赛////////////////////////////////////////////////////////
type cfTeamContests struct {
	moduleDetail[models.CfTeamContestsResponse]
	useAccount string
	count      int
}

func (r *cfTeamContests) GetCfTeamContests() error {
	r.count = 0
	GroupCode := utils.JsonDB.Get("groupCode").(string)
	url, err := utils.GenerateCFurlInstance.Contest.List(r.useAccount, &models.ContestListParams{
		GroupCode: GroupCode,
	})
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取url失败 %v", err))
		return err
	}
	resp, err := r.req.Get(url, map[string]string{})
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
	return db.Cf_team_contests{
		Contest_id:   int(cfTeamContest.Id),
		Contest_name: cfTeamContest.Name,
		PrePare_by:   cfTeamContest.PreparedBy,
		Start_time:   time.Unix(cfTeamContest.StartTimeSeconds, 0).Add(8 * time.Hour),
	}
}
func (r *cfTeamContests) GetLog() []string {
	return r.log.GetLog()
}
func (r *cfTeamContests) GetErr() []string {
	return r.log.GetErr()
}
func (r *cfTeamContests) Clear() error {
	r.log.ClearLog()
	r.log.ClearErr()
	r.count = 0
	return nil
}

// //////////////////////////////////////////////// 获取cf团队比赛题目////////////////////////////////////////////////////////
type cfTeamContestsProblems struct {
	moduleDetail[models.CfTeamContestProblemsResponse]
	count int
}

func (r *cfTeamContestsProblems) GetCfTeamContestsProblems(concurrency int) error {
	r.count = 0
	teamContests, err := r.repo.GetTeamContests()
	conCurrenter := utils.NewConCurrenter[db.Cf_team_contests](concurrency)
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
			Count:          50000,
			ShowUnofficial: true,
		})
		if err != nil {
			r.log.AddErr(fmt.Sprintf("获取团队比赛题目失败 %v", err))
			return err
		}
		problems, err := r.req.Get(url, map[string]string{})
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
func (r *cfTeamContestsProblems) GetLog() []string {
	return r.log.GetLog()
}
func (r *cfTeamContestsProblems) GetErr() []string {
	return r.log.GetErr()
}
func (r *cfTeamContestsProblems) Clear() error {
	r.log.ClearLog()
	r.log.ClearErr()
	r.count = 0
	return nil
}

// ////////////////////////////////////////////// 计算cf未过题////////////////////////////////////////////////////////
// type cfUnsolvedProblems struct {
// 	moduleDetail[db.Cf_all_submissions]
// 	count int
// }

// func (r *cfUnsolvedProblems) CalCfUnsolvedProblems(concurrency int) error {
// 	cfUserDatas, err := r.repo.GetCfAccountData()
// 	if err != nil {
// 		r.log.AddErr(fmt.Sprintf("获取cf用户数据失败 %v", err))
// 		return err
// 	}
// 	conCurrenter := utils.NewConCurrenter[models.CfUserData](concurrency)
// 	conCurrenter.Run(cfUserDatas, func(cfUserData models.CfUserData) error {
// 		cfSubmissions, err := r.repo.GetCfSubmissions(cfUserData.Account)
// 		if err != nil {
// 			r.log.AddErr(fmt.Sprintf("获取cf用户提交记录失败 %v", err))
// 			return err
// 		}
// 		err = r.handleCfUnsolvedProblems(&cfSubmissions)
// 		if err != nil {
// 			r.log.AddErr(fmt.Sprintf("处理cf用户提交记录失败 %v", err))
// 			return err
// 		}
// 		r.log.AddLog(fmt.Sprintf("新增%s %d条cf未过题", cfUserData.RealName, r.count))
// 		return nil
// 	})
// 	return nil
// }
// func (r *cfUnsolvedProblems) handleCfUnsolvedProblems(cfSubmissions *[]db.Cf_all_submissions) error {
// 	for _, cfSubmission := range *cfSubmissions {
// 		if cfSubmission.Verdict != "OK" {
// 			r.count++
// 		}
// 	}
// 	return nil
// }

// ////////////////////////////////////////////// 获取cf官方比赛////////////////////////////////////////////////////////
type cfOfficialContests struct {
	moduleDetail[models.CfOfficialContestsResponse]
	count int
}

func (r *cfOfficialContests) GetCfOfficialContests() error {
	r.count = 0
	GroupCode := utils.JsonDB.Get("groupCode").(string)
	url, err := utils.GenerateCFurlInstance.Contest.List("", &models.ContestListParams{
		GroupCode: GroupCode,
	})
	if err != nil {
		r.log.AddErr(fmt.Sprintf("获取url失败 %v", err))
		return err
	}
	resp, err := r.req.Get(url, map[string]string{})
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
		Start_time:            time.Unix(cfTeamContest.StartTimeSeconds, 0).Add(8 * time.Hour),
	}
}
func (r *cfOfficialContests) GetLog() []string {
	return r.log.GetLog()
}
func (r *cfOfficialContests) GetErr() []string {
	return r.log.GetErr()
}
func (r *cfOfficialContests) Clear() error {
	r.log.ClearLog()
	r.log.ClearErr()
	r.count = 0
	return nil
}
