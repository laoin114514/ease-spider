package cheat_judge

import (
	"fmt"
	"time"
)

// 检测所有role_id=3的用户
func (cd *CheatingDetector) DetectAllUsers() ([]CheatingDetectionResult, error) {
	// 获取所有role_id=3的用户
	query := `
		SELECT u.id, u.real_name, oa.luogu_uid 
		FROM user u 
		JOIN oj_account oa ON u.id = oa.user_id 
		WHERE u.role_id = 3 AND oa.luogu_uid IS NOT NULL AND oa.luogu_uid != ''
	`

	rows, err := cd.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	defer rows.Close()

	var results []CheatingDetectionResult

	for rows.Next() {
		var userID int
		var username, luoguUID string
		err := rows.Scan(&userID, &username, &luoguUID)
		if err != nil {
			continue
		}

		// 检测单个用户
		result, err := cd.DetectUser(luoguUID, username)
		if err != nil {
			fmt.Printf("检测用户 %s 失败: %v\n", username, err)
			continue
		}

		results = append(results, *result)
	}

	return results, nil
}

// 检测单个用户
func (cd *CheatingDetector) DetectUser(uid, username string) (*CheatingDetectionResult, error) {
	// 获取用户提交记录
	submissions, err := cd.getUserSubmissions(uid)
	if err != nil {
		return nil, err
	}

	if len(submissions) == 0 {
		return &CheatingDetectionResult{
			UID:           uid,
			Username:      username,
			RiskScore:     0,
			RiskLevel:     RISK_LOW,
			Indicators:    []DetectionIndicator{},
			DetectionTime: time.Now(),
		}, nil
	}

	// 计算用户统计信息
	stats := cd.calculateUserStats(submissions)

	// 运行各种检测器
	var allIndicators []DetectionIndicator

	//检测逻辑模块
	// 1. 时间模式检测
	timeIndicators := cd.detectTimePatterns(submissions, stats)
	allIndicators = append(allIndicators, timeIndicators...)

	// 2. 难度跳跃检测
	difficultyIndicators := cd.detectDifficultyJumps(submissions, stats)
	allIndicators = append(allIndicators, difficultyIndicators...)

	// 3. 提交行为检测
	behaviorIndicators := cd.detectSubmissionBehavior(submissions, stats)
	allIndicators = append(allIndicators, behaviorIndicators...)

	// 4. AC率异常检测
	acRateIndicators := cd.detectACRateAnomalies(stats)
	allIndicators = append(allIndicators, acRateIndicators...)

	// 计算综合风险分数
	riskScore := cd.calculateRiskScore(allIndicators)
	riskLevel := cd.determineRiskLevel(riskScore)

	return &CheatingDetectionResult{
		UID:           uid,
		Username:      username,
		RiskScore:     riskScore,
		RiskLevel:     riskLevel,
		Indicators:    allIndicators,
		DetectionTime: time.Now(),
	}, nil
}
