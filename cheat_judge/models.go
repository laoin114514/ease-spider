package cheat_judge

import (
	"database/sql"
	"spider/config/db"
	"time"
)

// 作弊检测结果结构
type CheatingDetectionResult struct {
	UID           string               `json:"uid"`
	Username      string               `json:"username"`
	RiskScore     float64              `json:"risk_score"`
	RiskLevel     string               `json:"risk_level"`
	Indicators    []DetectionIndicator `json:"indicators"`
	DetectionTime time.Time            `json:"detection_time"`
}

// 检测指标
type DetectionIndicator struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"` // LOW, MEDIUM, HIGH
	Score       float64 `json:"score"`
	Description string  `json:"description"`
	Evidence    string  `json:"evidence"`
}

// 用户提交统计
type UserSubmissionStats struct {
	UID              string
	TotalSubmissions int
	TotalAC          int
	ACRate           float64
	AvgTimeBetweenAC time.Duration
	DifficultyCount  map[string]int
	RecentACs        []db.Luogu_all_submissions
}

// 风险等级阈值
const (
	RISK_LOW      = "LOW"      // 0-30分
	RISK_MEDIUM   = "MEDIUM"   // 31-60分
	RISK_HIGH     = "HIGH"     // 61-80分
	RISK_CRITICAL = "CRITICAL" // 81-100分
)

// 检测器结构
type CheatingDetector struct {
	db *sql.DB
}

// 难度权重映射
var difficultyWeights = map[string]int{
	"grey":   1,
	"red":    2,
	"brown":  3,
	"yellow": 4,
	"green":  5,
	"blue":   6,
	"purple": 7,
	"black":  8,
}
