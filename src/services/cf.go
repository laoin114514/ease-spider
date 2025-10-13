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
	repo  *repository.CfRepository
	debug *utils.Debug
}

func NewCfService() *CfService {
	return &CfService{
		CfUserStatus: cfUserStatus{
			moduleDetail: moduleDetail[models.CfUserStatusResponse]{
				LogService: NewLogService("logs/cfUserStatus.log", "logs/cfUserStatus.err.log"),
				repo:       repository.NewCfRepository(),
				debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
			},
		},
		CfOfficialProblems: cfOfficialProblems{
			moduleDetail: moduleDetail[models.CfOfficialProblemsResponse]{
				LogService: NewLogService("logs/cfOfficialProblems.log", "logs/cfOfficialProblems.err.log"),
				repo:       repository.NewCfRepository(),
				debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
			},
			count: 0,
		},
		CfTeamContests: cfTeamContests{
			moduleDetail: moduleDetail[models.CfTeamContestsResponse]{
				LogService: NewLogService("logs/cfTeamContests.log", "logs/cfTeamContests.err.log"),
				repo:       repository.NewCfRepository(),
				debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
			},
			count:      0,
			useAccount: config.AppConfig.Cf.ManagerAccount,
		},
		CfTeamContestsProblems: cfTeamContestsProblems{
			moduleDetail: moduleDetail[models.CfTeamContestProblemsResponse]{
				LogService: NewLogService("logs/cfTeamContestsProblems.log", "logs/cfTeamContestsProblems.err.log"),
				repo:       repository.NewCfRepository(),
				debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
			},
			count: 0,
		},
		CfOfficialContests: cfOfficialContests{
			moduleDetail: moduleDetail[models.CfOfficialContestsResponse]{
				LogService: NewLogService("logs/cfOfficialContests.log", "logs/cfOfficialContests.err.log"),
				repo:       repository.NewCfRepository(),
				debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
			},
			count: 0,
		},
	}
}

// ================================获取CF用户提交记录===============================================
type cfUserStatus struct {
	moduleDetail[models.CfUserStatusResponse]
}

// ================================公有接口方法===============================================

// GetCfRecords 获取CF用户提交记录 - 对外提供的主要接口
func (r *cfUserStatus) GetCfRecords() error {
	r.debug.Debug("开始获取CF用户提交记录")

	// 构建并发工具类
	conCurrenter := utils.NewConCurrenter[models.CfUserData](config.AppConfig.Cf.CfRecordsConcurrency)

	// 获取CF用户数据
	cfUserDatas, err := r.repo.GetCfAccountData()
	if err != nil {
		r.logError("获取CF用户数据失败", err)
		return err
	}
	r.debug.Debug(fmt.Sprintf("成功获取到 %d 个CF用户", len(cfUserDatas)))

	// 并发获取CF提交记录
	conCurrenter.Run(cfUserDatas, func(cfUserData models.CfUserData) error {
		return r.processUserRecords(cfUserData)
	})

	r.debug.Debug("CF用户提交记录获取完成")
	return nil
}

// ================================私有实现方法===============================================

// processUserRecords 处理单个用户的提交记录 - 私有方法
func (r *cfUserStatus) processUserRecords(cfUserData models.CfUserData) error {

	// 获取数据库中已有的提交记录
	err := r.loadExistingRecords(&cfUserData)
	if err != nil {
		r.logError(fmt.Sprintf("%s获取数据库中已有提交记录失败", cfUserData.RealName), err)
		return err
	}

	// 构建请求URL
	url, err := r.buildRequestURL(cfUserData.Account, true)
	if err != nil {
		r.logError(fmt.Sprintf("%s构建请求URL失败", cfUserData.RealName), err)
		return err
	}

	// 发起请求
	resp, err := r.fetchUserStatus(url)
	if err != nil {
		r.logError(fmt.Sprintf("%s请求数据失败，尝试无apiKey请求", cfUserData.RealName), err)
		url, err = r.buildRequestURL(cfUserData.Account, false)
		if err != nil {
			r.logError(fmt.Sprintf("%s构建请求URL失败", cfUserData.RealName), err)
			return err
		}
		resp, err = r.fetchUserStatus(url)
		if err != nil {
			r.logError(fmt.Sprintf("%s请求数据失败,URL:%s", cfUserData.RealName, url), err)
			return err
		}
	}

	// 处理CF提交记录
	err = r.processSubmissionRecords(&cfUserData, &resp)
	if err != nil {
		r.logError(fmt.Sprintf("%s处理提交记录失败", cfUserData.RealName), err)
		return err
	}

	// 如果插入记录为0，则不进行日志记录
	if cfUserData.InsertCount == 0 {
		return nil
	}
	r.logSuccess(fmt.Sprintf("%s处理提交记录成功", cfUserData.RealName), cfUserData.InsertCount)
	return nil
}

// loadExistingRecords 加载数据库中已有的提交记录 - 私有方法
func (r *cfUserStatus) loadExistingRecords(cfUserData *models.CfUserData) error {
	oldDataSet, err := r.repo.GetCfRecordsIdToset(cfUserData.Account)
	if err != nil {
		return err
	}
	cfUserData.OldDataSet = oldDataSet
	return nil
}

// buildRequestURL 构建请求URL - 私有方法
func (r *cfUserStatus) buildRequestURL(account string, isApiKey bool) (string, error) {
	// 尝试使用HTTPS
	url, err := utils.GenerateCFurlInstance.User.Status(
		isApiKey,
		&models.UserStatusParams{
			Handle: account,
			From:   1,
			Count:  constants.CfMaxRecords,
		},
	)
	if err != nil {
		// 回退到HTTP
		url, err = utils.GenerateCFurlInstance.User.Status(
			false,
			&models.UserStatusParams{
				Handle: account,
				From:   1,
				Count:  constants.CfMaxRecords,
			},
		)
	}
	return url, err
}

// fetchUserStatus 获取用户状态数据 - 私有方法
func (r *cfUserStatus) fetchUserStatus(url string) (models.CfUserStatusResponse, error) {
	req := utils.NewRequest[models.CfUserStatusResponse]()
	return req.Get(url, map[string]string{})
}

// processSubmissionRecords 处理提交记录 - 私有方法
func (r *cfUserStatus) processSubmissionRecords(cfUserData *models.CfUserData, resp *models.CfUserStatusResponse) error {

	skippedCount := 0
	for _, cfRecord := range resp.Result {
		if r.shouldSkipRecord(cfRecord, cfUserData) {
			skippedCount++
			continue
		}

		table := r.buildSubmissionTable(&cfRecord, cfUserData)
		err := db.Insert_cf_all_sub(table)
		if err != nil {
			r.logError(fmt.Sprintf("%s插入提交记录失败", cfUserData.RealName), err)
			return err
		}
		cfUserData.InsertCount++
	}

	r.debug.Debug(fmt.Sprintf("用户 %s 提交记录处理完成，跳过 %d 条，新增 %d 条", cfUserData.RealName, skippedCount, cfUserData.InsertCount))
	return nil
}

// shouldSkipRecord 判断是否应该跳过该提交记录 - 私有方法
func (r *cfUserStatus) shouldSkipRecord(cfRecord models.CfSubmission, cfUserData *models.CfUserData) bool {
	// 如果提交记录已存在，则跳过
	if cfUserData.OldDataSet[int(cfRecord.Id)] {
		return true
	}
	// 如果提交记录还在测试中，则跳过
	if cfRecord.Verdict == "TESTING" {
		return true
	}
	return false
}

// buildSubmissionTable 构建提交记录表 - 私有方法
func (r *cfUserStatus) buildSubmissionTable(cfRecord *models.CfSubmission, cfUserData *models.CfUserData) db.Cf_all_submissions {
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

// ================================辅助方法===============================================

// logError 记录错误日志 - 私有方法
func (r *cfUserStatus) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *cfUserStatus) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}

// ================================获取CF官方题目===============================================
type cfOfficialProblems struct {
	moduleDetail[models.CfOfficialProblemsResponse]
	count int
}

// ================================公有接口方法===============================================

// GetCfOfficialProblems 获取CF官方题目 - 对外提供的主要接口
func (r *cfOfficialProblems) GetCfOfficialProblems() error {
	r.count = 0

	// 构建请求URL
	url, err := r.buildProblemsURL()
	if err != nil {
		r.logError("构建请求URL失败", err)
		return err
	}

	// 获取数据
	resp, err := r.fetchProblemsData(url)
	if err != nil {
		r.logError("获取数据失败", err)
		return err
	}

	// 处理官方题目数据
	err = r.processOfficialProblems(&resp)
	if err != nil {
		r.logError("处理数据失败", err)
		return err
	}

	r.logSuccess("插入官方题目", r.count)
	return nil
}

// ================================私有实现方法===============================================

// buildProblemsURL 构建题目请求URL - 私有方法
func (r *cfOfficialProblems) buildProblemsURL() (string, error) {
	return utils.GenerateCFurlInstance.ProblemSet.Problems(&models.ProblemsetProblemsParams{})
}

// fetchProblemsData 获取题目数据 - 私有方法
func (r *cfOfficialProblems) fetchProblemsData(url string) (models.CfOfficialProblemsResponse, error) {
	req := utils.NewRequest[models.CfOfficialProblemsResponse]()
	return req.Get(url, map[string]string{})
}

// processOfficialProblems 处理官方题目数据 - 私有方法
func (r *cfOfficialProblems) processOfficialProblems(resp *models.CfOfficialProblemsResponse) error {
	for _, cfProblem := range resp.Result.Problems {
		table := r.buildProblemTable(&cfProblem)
		err := db.Insert_cf_official(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}

// buildProblemTable 构建题目表 - 私有方法
func (r *cfOfficialProblems) buildProblemTable(cfProblem *models.CfProblem) db.Cf_official_problems {
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

// ================================辅助方法===============================================

// logError 记录错误日志 - 私有方法
func (r *cfOfficialProblems) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *cfOfficialProblems) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}

// ================================获取CF团队比赛===============================================
type cfTeamContests struct {
	moduleDetail[models.CfTeamContestsResponse]
	useAccount string
	count      int
}

// ================================公有接口方法===============================================

// GetCfTeamContests 获取CF团队比赛 - 对外提供的主要接口
func (r *cfTeamContests) GetCfTeamContests() error {
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

// ================================私有实现方法===============================================

// buildTeamContestsURL 构建团队比赛请求URL - 私有方法
func (r *cfTeamContests) buildTeamContestsURL() (string, error) {
	GroupCode := config.AppConfig.Cf.GroupCode
	return utils.GenerateCFurlInstance.Contest.List(r.useAccount, &models.ContestListParams{
		GroupCode: GroupCode,
	})
}

// fetchTeamContestsData 获取团队比赛数据 - 私有方法
func (r *cfTeamContests) fetchTeamContestsData(url string) (models.CfTeamContestsResponse, error) {
	req := utils.NewRequest[models.CfTeamContestsResponse]()
	return req.Get(url, map[string]string{})
}

// processTeamContests 处理团队比赛数据 - 私有方法
func (r *cfTeamContests) processTeamContests(resp *models.CfTeamContestsResponse) error {
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
func (r *cfTeamContests) buildTeamContestTable(cfTeamContest *models.CfContest) db.Cf_team_contests {
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

// ================================辅助方法===============================================

// logError 记录错误日志 - 私有方法
func (r *cfTeamContests) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *cfTeamContests) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}

// ================================获取CF团队比赛题目===============================================
type cfTeamContestsProblems struct {
	moduleDetail[models.CfTeamContestProblemsResponse]
	count int
}

// ================================公有接口方法===============================================

// GetCfTeamContestsProblems 获取CF团队比赛题目 - 对外提供的主要接口
func (r *cfTeamContestsProblems) GetCfTeamContestsProblems() error {
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

// ================================私有实现方法===============================================

// processTeamContestProblems 处理单个团队比赛的题目 - 私有方法
func (r *cfTeamContestsProblems) processTeamContestProblems(teamContest db.Cf_team_contests) error {
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

	return nil
}

// buildTeamContestProblemsURL 构建团队比赛题目请求URL - 私有方法
func (r *cfTeamContestsProblems) buildTeamContestProblemsURL(teamContest db.Cf_team_contests) (string, error) {
	return utils.GenerateCFurlInstance.Contest.Standings(teamContest.PrePare_by, &models.ContestStandingsParams{
		ContestID:      teamContest.Contest_id,
		AsManager:      true,
		From:           1,
		Count:          constants.CfMaxRecords,
		ShowUnofficial: true,
	})
}

// fetchTeamContestProblemsData 获取团队比赛题目数据 - 私有方法
func (r *cfTeamContestsProblems) fetchTeamContestProblemsData(url string) (models.CfTeamContestProblemsResponse, error) {
	req := utils.NewRequest[models.CfTeamContestProblemsResponse]()
	return req.Get(url, map[string]string{})
}

// processTeamContestProblemsData 处理团队比赛题目数据 - 私有方法
func (r *cfTeamContestsProblems) processTeamContestProblemsData(resp *models.CfTeamContestProblemsResponse) error {
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
func (r *cfTeamContestsProblems) buildTeamProblemTable(cfProblem *models.CfProblem) db.Cf_team_problems {
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

// ================================辅助方法===============================================

// logError 记录错误日志 - 私有方法
func (r *cfTeamContestsProblems) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *cfTeamContestsProblems) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}

// ================================获取CF官方比赛===============================================
type cfOfficialContests struct {
	moduleDetail[models.CfOfficialContestsResponse]
	count int
}

// ================================公有接口方法===============================================

// GetCfOfficialContests 获取CF官方比赛 - 对外提供的主要接口
func (r *cfOfficialContests) GetCfOfficialContests() error {
	r.count = 0

	// 构建请求URL
	url, err := r.buildOfficialContestsURL()
	if err != nil {
		r.logError("构建请求URL失败", err)
		return err
	}

	// 获取数据
	resp, err := r.fetchOfficialContestsData(url)
	if err != nil {
		r.logError("获取数据失败", err)
		return err
	}

	// 处理官方比赛数据
	err = r.processOfficialContests(&resp)
	if err != nil {
		r.logError("处理数据失败", err)
		return err
	}

	r.logSuccess("插入官方比赛", r.count)
	return nil
}

// ================================私有实现方法===============================================

// buildOfficialContestsURL 构建官方比赛请求URL - 私有方法
func (r *cfOfficialContests) buildOfficialContestsURL() (string, error) {
	GroupCode := config.AppConfig.Cf.GroupCode
	return utils.GenerateCFurlInstance.Contest.List("", &models.ContestListParams{
		GroupCode: GroupCode,
	})
}

// fetchOfficialContestsData 获取官方比赛数据 - 私有方法
func (r *cfOfficialContests) fetchOfficialContestsData(url string) (models.CfOfficialContestsResponse, error) {
	req := utils.NewRequest[models.CfOfficialContestsResponse]()
	return req.Get(url, map[string]string{})
}

// processOfficialContests 处理官方比赛数据 - 私有方法
func (r *cfOfficialContests) processOfficialContests(resp *models.CfOfficialContestsResponse) error {
	for _, cfTeamContest := range resp.Result {
		table := r.buildOfficialContestTable(&cfTeamContest)
		err := db.Insert_cf_official_contests(table)
		if err != nil {
			continue
		}
		r.count++
	}
	return nil
}

// buildOfficialContestTable 构建官方比赛表 - 私有方法
func (r *cfOfficialContests) buildOfficialContestTable(cfTeamContest *models.CfContest) db.Cf_official_contests {
	return db.Cf_official_contests{
		Official_contest_id:   int(cfTeamContest.Id),
		Official_contest_name: cfTeamContest.Name,
		Phase:                 cfTeamContest.Phase,
		Start_time:            time.Unix(cfTeamContest.StartTimeSeconds, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}

// ================================辅助方法===============================================

// logError 记录错误日志 - 私有方法
func (r *cfOfficialContests) logError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += fmt.Sprintf(" %v", err)
	}
	r.log.AddErr(errorMsg)
	r.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志 - 私有方法
func (r *cfOfficialContests) logSuccess(message string, count int) {
	successMsg := fmt.Sprintf("%s %d条", message, count)
	r.log.AddLog(successMsg)
	r.debug.Debug(successMsg)
}
