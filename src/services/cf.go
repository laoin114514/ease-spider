package services

import (
	"fmt"
	"spider/config/db"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"spider/src/utils/genUrl"
	"time"
)

type CfService struct {
	CfUserStatus cfUserStatus
}

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
	}
}

// //////////////////////////////////////////////// 获取cf用户提交记录////////////////////////////////////////////////////////

type cfUserStatus struct {
	moduleDetail[models.CfUserStatusResponse]
}

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
		cfUserData.OldDataSet, err = r.repo.GetCfRecordsInDbToset(cfUserData.Account)
		if err != nil {
			r.log.AddErr(fmt.Sprintf("%s 获取db中已有的提交记录失败 %v", cfUserData.RealName, err))
			return err
		}

		//按照cf规则拼接url
		url, err := genUrl.NewUser().Status(true, genUrl.User_status{
			Handle:         cfUserData.Account,
			From:           1,
			Count:          50000,
			IncludeSources: false,
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

// //////////////////////////////////////////////// 获取cf官方题目////////////////////////////////////////////////////////
type cfOfficialProblems struct {
	moduleDetail[models.CfOfficialProblemsResponse]
}
