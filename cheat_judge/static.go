package cheat_judge

// 计算综合风险分数
func (cd *CheatingDetector) calculateRiskScore(indicators []DetectionIndicator) float64 {
	var totalScore float64

	for _, indicator := range indicators {
		// 根据严重程度调整权重
		weight := 1.0
		switch indicator.Severity {
		case "LOW":
			weight = 0.5
		case "MEDIUM":
			weight = 1.0
		case "HIGH":
			weight = 1.5
		}

		totalScore += indicator.Score * weight
	}

	// 限制最高分数为100
	if totalScore > 100 {
		totalScore = 100
	}

	return totalScore
}

// 确定风险等级
func (cd *CheatingDetector) determineRiskLevel(score float64) string {
	switch {
	case score <= 30:
		return RISK_LOW
	case score <= 60:
		return RISK_MEDIUM
	case score <= 80:
		return RISK_HIGH
	default:
		return RISK_CRITICAL
	}
}
