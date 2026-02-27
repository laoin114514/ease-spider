package services

import (
	"spider/internal/repository"
	"spider/internal/utils"
)

// CfService 是CF服务的主结构，包含所有CF相关的子服务
type CfService struct {
	CfUserStatus           *CfUserStatus
	CfOfficialProblems     *CfOfficialProblems
	CfTeamContests         *CfTeamContests
	CfTeamContestsProblems *CfTeamContestsProblems
	CfOfficialContests     *CfOfficialContests
}

// NewCfService 创建并初始化CF服务的主结构
func NewCfService() *CfService {
	return &CfService{
		CfUserStatus:           NewCfUserStatus(),
		CfOfficialProblems:     NewCfOfficialProblems(),
		CfTeamContests:         NewCfTeamContests(),
		CfTeamContestsProblems: NewCfTeamContestsProblems(),
		CfOfficialContests:     NewCfOfficialContests(),
	}
}

// T类型是请求响应的结构体
type moduleDetail[T any] struct {
	*LogService
	repo  *repository.CfRepository
	debug *utils.Debug
}
