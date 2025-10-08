package cheat_judge

import (
	"fmt"
	"spider/config/db"
	"time"
)

// 时间模式检测
func (cd *CheatingDetector) detectTimePatterns(submissions []db.Luogu_all_submissions, stats UserSubmissionStats) []DetectionIndicator {
	var indicators []DetectionIndicator
	//难度区分，加权
	// 检测快速连续AC（30分钟内AC多道题），30分钟是否合适待定
	rapidACCount := 0
	for i := 1; i < len(stats.RecentACs); i++ {
		timeDiff := stats.RecentACs[i].Creation_time.Sub(stats.RecentACs[i-1].Creation_time)
		if timeDiff <= 30*time.Minute {
			rapidACCount++
		}
	}

	if rapidACCount >= 3 {
		indicators = append(indicators, DetectionIndicator{
			Type:        "RAPID_AC",
			Severity:    "HIGH",
			Score:       20.0,
			Description: "短时间内连续AC多道题目",
			Evidence:    fmt.Sprintf("30分钟内连续AC %d 道题", rapidACCount),
		})
	}

	// 检测异常提交时间（凌晨2-6点大量提交）这个另外说
	nightSubmissions := 0
	for _, sub := range submissions {
		hour := sub.Creation_time.Hour()
		if hour >= 2 && hour <= 6 {
			nightSubmissions++
		}
	}

	if nightSubmissions > len(submissions)/4 {
		indicators = append(indicators, DetectionIndicator{
			Type:        "NIGHT_SUBMISSIONS",
			Severity:    "MEDIUM",
			Score:       10.0,
			Description: "异常时间大量提交",
			Evidence: fmt.Sprintf("凌晨2-6点提交 %d 次，占总提交的 %.1f%%",
				nightSubmissions, float64(nightSubmissions)/float64(len(submissions))*100),
		})
	}

	return indicators
}

// 难度跳跃检测
func (cd *CheatingDetector) detectDifficultyJumps(submissions []db.Luogu_all_submissions, stats UserSubmissionStats) []DetectionIndicator {
	var indicators []DetectionIndicator

	// 获取AC的题目按时间排序
	var acSubmissions []db.Luogu_all_submissions
	for _, sub := range submissions {
		if sub.Is_pass {
			acSubmissions = append(acSubmissions, sub)
		}
	}

	// 检测难度跳跃
	for i := 1; i < len(acSubmissions); i++ {
		prevWeight := difficultyWeights[acSubmissions[i-1].Difficulty]
		currWeight := difficultyWeights[acSubmissions[i].Difficulty]

		// 如果难度跳跃超过3个等级，等级待定
		if currWeight-prevWeight > 3 {
			indicators = append(indicators, DetectionIndicator{
				Type:        "DIFFICULTY_JUMP",
				Severity:    "MEDIUM",
				Score:       15.0,
				Description: "不合理的难度跳跃",
				Evidence: fmt.Sprintf("从 %s 直接跳到 %s",
					acSubmissions[i-1].Difficulty, acSubmissions[i].Difficulty),
			})
		}
	}

	// 检测缺乏基础题练习,题数待定
	basicCount := stats.DifficultyCount["grey"] + stats.DifficultyCount["red"] + stats.DifficultyCount["brown"]
	advancedCount := stats.DifficultyCount["purple"] + stats.DifficultyCount["black"]

	if basicCount < 10 && advancedCount > 5 {
		indicators = append(indicators, DetectionIndicator{
			Type:        "INSUFFICIENT_BASICS",
			Severity:    "HIGH",
			Score:       25.0,
			Description: "缺乏基础题练习但能做高难度题",
			Evidence:    fmt.Sprintf("基础题仅 %d 道，高难度题 %d 道", basicCount, advancedCount),
		})
	}

	return indicators
}

// 提交行为检测
func (cd *CheatingDetector) detectSubmissionBehavior(submissions []db.Luogu_all_submissions, stats UserSubmissionStats) []DetectionIndicator {
	var indicators []DetectionIndicator

	// 检测提交频率异常（1小时内大量提交）
	submissionTimes := make([]time.Time, len(submissions))
	for i, sub := range submissions {
		submissionTimes[i] = sub.Creation_time
	}

	// 按时间窗口统计提交频率
	windowSize := time.Hour
	maxSubmissionsInWindow := 0

	for i := 0; i < len(submissionTimes); i++ {
		windowStart := submissionTimes[i]
		windowEnd := windowStart.Add(windowSize)
		count := 0

		for j := i; j < len(submissionTimes) && submissionTimes[j].Before(windowEnd); j++ {
			count++
		}

		if count > maxSubmissionsInWindow {
			maxSubmissionsInWindow = count
		}
	}

	if maxSubmissionsInWindow > 20 {
		indicators = append(indicators, DetectionIndicator{
			Type:        "HIGH_FREQUENCY",
			Severity:    "MEDIUM",
			Score:       12.0,
			Description: "提交频率异常",
			Evidence:    fmt.Sprintf("1小时内最多提交 %d 次", maxSubmissionsInWindow),
		})
	}

	return indicators
}

// // AC率异常检测
func (cd *CheatingDetector) detectACRateAnomalies(stats UserSubmissionStats) []DetectionIndicator {
	var indicators []DetectionIndicator

	// 检测AC率过高且缺乏试错
	if stats.ACRate > 0.8 && stats.TotalSubmissions > 50 {
		// 计算平均每道AC题目的尝试次数
		// 这里需要重新查询详细数据，简化处理
		avgAttemptsPerAC := float64(stats.TotalSubmissions) / float64(stats.TotalAC)

		if avgAttemptsPerAC < 1.5 {
			indicators = append(indicators, DetectionIndicator{
				Type:        "HIGH_AC_RATE",
				Severity:    "HIGH",
				Score:       20.0,
				Description: "AC率过高且缺乏试错过程",
				Evidence: fmt.Sprintf("AC率 %.1f%%，平均每道AC题目尝试 %.1f 次",
					stats.ACRate*100, avgAttemptsPerAC),
			})
		}
	}

	return indicators
}
