package cheat_judge

import (
	"fmt"
	"time"
)

// 生成检测报告
func (cd *CheatingDetector) GenerateReport(results []CheatingDetectionResult) {
	fmt.Println("=== 洛谷作弊检测报告 ===")
	fmt.Printf("检测时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("检测用户数: %d\n\n", len(results))

	// 按风险等级分组
	riskGroups := make(map[string][]CheatingDetectionResult)
	for _, result := range results {
		riskGroups[result.RiskLevel] = append(riskGroups[result.RiskLevel], result)
	}

	// 输出高风险用户
	if highRiskUsers := riskGroups[RISK_HIGH]; len(highRiskUsers) > 0 {
		fmt.Println("🔴 高风险用户:")
		for _, user := range highRiskUsers {
			fmt.Printf("  - %s (UID: %s) - 风险分数: %.1f\n",
				user.Username, user.UID, user.RiskScore)
			for _, indicator := range user.Indicators {
				if indicator.Severity == "HIGH" {
					fmt.Printf("    ⚠️  %s: %s\n", indicator.Type, indicator.Description)
				}
			}
		}
		fmt.Println()
	}

	// 输出极高风险用户
	if criticalUsers := riskGroups[RISK_CRITICAL]; len(criticalUsers) > 0 {
		fmt.Println("🚨 极高风险用户:")
		for _, user := range criticalUsers {
			fmt.Printf("  - %s (UID: %s) - 风险分数: %.1f\n",
				user.Username, user.UID, user.RiskScore)
			for _, indicator := range user.Indicators {
				fmt.Printf("    ⚠️  %s: %s\n", indicator.Type, indicator.Description)
			}
		}
		fmt.Println()
	}

	// 统计信息
	fmt.Println("📊 风险分布统计:")
	for level, users := range riskGroups {
		fmt.Printf("  %s: %d 人\n", level, len(users))
	}
}

// 保存检测结果到数据库
func (cd *CheatingDetector) SaveResults(results []CheatingDetectionResult) error {
	// 创建检测结果表（如果不存在）
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS cheating_detection_results (
		id INT PRIMARY KEY AUTO_INCREMENT,
		uid VARCHAR(20) NOT NULL,
		username VARCHAR(100) NOT NULL,
		risk_score DECIMAL(5,2) NOT NULL,
		risk_level ENUM('LOW', 'MEDIUM', 'HIGH', 'CRITICAL') NOT NULL,
		detection_time DATETIME NOT NULL,
		indicators JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_uid_time (uid, detection_time),
		INDEX idx_risk_level (risk_level)
	)`

	_, err := cd.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建检测结果表失败: %v", err)
	}

	// 插入检测结果
	insertSQL := `
	INSERT INTO cheating_detection_results 
	(uid, username, risk_score, risk_level, detection_time, indicators) 
	VALUES (?, ?, ?, ?, ?, ?)`

	for _, result := range results {
		// 将indicators转换为JSON
		indicatorsJSON := `[`
		for i, indicator := range result.Indicators {
			if i > 0 {
				indicatorsJSON += `,`
			}
			indicatorsJSON += fmt.Sprintf(`{"type":"%s","severity":"%s","score":%.1f,"description":"%s","evidence":"%s"}`,
				indicator.Type, indicator.Severity, indicator.Score,
				indicator.Description, indicator.Evidence)
		}
		indicatorsJSON += `]`

		_, err := cd.db.Exec(insertSQL,
			result.UID, result.Username, result.RiskScore,
			result.RiskLevel, result.DetectionTime, indicatorsJSON)
		if err != nil {
			fmt.Printf("保存用户 %s 检测结果失败: %v\n", result.Username, err)
		}
	}

	return nil
}

// 查询历史检测结果
func (cd *CheatingDetector) GetHistoryResults(limit int) ([]CheatingDetectionResult, error) {
	query := `
	SELECT uid, username, risk_score, risk_level, detection_time, indicators
	FROM cheating_detection_results 
	ORDER BY detection_time DESC 
	LIMIT ?
	`

	rows, err := cd.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CheatingDetectionResult
	for rows.Next() {
		var result CheatingDetectionResult
		var indicatorsJSON string

		err := rows.Scan(&result.UID, &result.Username, &result.RiskScore,
			&result.RiskLevel, &result.DetectionTime, &indicatorsJSON)
		if err != nil {
			continue
		}

		// 这里简化处理JSON解析，实际项目中应该使用json包
		result.Indicators = []DetectionIndicator{}
		results = append(results, result)
	}

	return results, nil
}
