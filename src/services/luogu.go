package services

import (
	"spider/config"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
)

// Luogu 是洛谷服务的主结构，包含所有洛谷相关的子服务
type Luogu struct {
	LuoguSubmissionDetail *LuoguSubmissionDetail
	LuoguCookie           *LuoguUpdateCookie
	LuoguSolution         *LuoguSolution
	LuoguRecords          *LuoguRecords
	LuoguTeam             *LuoguTeam
}

// NewLuogu 创建并初始化洛谷服务的主结构
func NewLuogu() *Luogu {
	return &Luogu{
		LuoguSubmissionDetail: NewLuoguSubmissionDetail(),
		LuoguCookie:           NewLuoguUpdateCookie(),
		LuoguSolution:         NewLuoguSolution(),
		LuoguRecords:          NewLuoguRecords(),
		LuoguTeam:             NewLuoguTeam(),
	}
}

// ================================获取洛谷题解===============================================
func NewLuoguSolution() *LuoguSolution {
	return &LuoguSolution{
		LogService: NewLogService("logs/luoguSolution.log", "logs/luoguSolution.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		totalPage:  0,
		solutions:  []models.SolutionContent{},
	}
}

// ============================================爬取洛谷提交记录源代码===============================================
func NewLuoguSubmissionDetail() *LuoguSubmissionDetail {
	return &LuoguSubmissionDetail{
		LogService: NewLogService("logs/luoguSubmissionDetail.log", "logs/luoguSubmissionDetail.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		count:      0,
	}
}

// ============================================获取洛谷团队成员==========================================
func NewLuoguTeam() *LuoguTeam {
	return &LuoguTeam{
		LogService: NewLogService("logs/luoguTeam.log", "logs/luoguTeam.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		count:      0,
		results:    make(map[string]bool),
	}
}

// ================================获取洛谷用户提交记录===============================================
func NewLuoguRecords() *LuoguRecords {
	return &LuoguRecords{
		LogService: NewLogService("logs/luogu.log", "logs/luogu.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
	}
}

// ================================更新洛谷Cookie===============================================
func NewLuoguUpdateCookie() *LuoguUpdateCookie {
	return &LuoguUpdateCookie{
		LogService: NewLogService("logs/luoguUpdateCookie.log", "logs/luoguUpdateCookie.err.log"),
		cookiePool: make(map[string]string),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
	}
}
