package cheat_judge

import (
	"spider/config/db"
	"time"
)

// 获取用户提交记录
func (cd *CheatingDetector) getUserSubmissions(uid string) ([]db.Luogu_all_submissions, error) {
	//T开头为私人比赛题目，role_id=3的用户为预备役
	query := `
		SELECT sub_id, uid, problem_id, problem_name, difficulty, is_pass, creation_time
		FROM luogu_all_submissions 
		WHERE uid = ? and left(problem_id, 1) != 'T'
		ORDER BY creation_time ASC
	`

	rows, err := cd.db.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []db.Luogu_all_submissions
	for rows.Next() {
		var sub db.Luogu_all_submissions
		err := rows.Scan(&sub.Sub_id, &sub.Uid, &sub.Problem_id, &sub.Problem_name,
			&sub.Difficulty, &sub.Is_pass, &sub.Creation_time)
		if err != nil {
			continue
		}
		submissions = append(submissions, sub)
	}

	return submissions, nil
}

// 计算用户统计信息
func (cd *CheatingDetector) calculateUserStats(submissions []db.Luogu_all_submissions) UserSubmissionStats {
	stats := UserSubmissionStats{
		UID:              submissions[0].Uid,
		TotalSubmissions: len(submissions),
		DifficultyCount:  make(map[string]int),
	}

	var acSubmissions []db.Luogu_all_submissions
	var acTimes []time.Time

	for _, sub := range submissions {
		if sub.Is_pass {
			stats.TotalAC++
			acSubmissions = append(acSubmissions, sub)
			acTimes = append(acTimes, sub.Creation_time)
		}
		stats.DifficultyCount[sub.Difficulty]++
	}

	if stats.TotalSubmissions > 0 {
		stats.ACRate = float64(stats.TotalAC) / float64(stats.TotalSubmissions)
	}

	// 计算平均AC间隔时间
	if len(acTimes) > 1 {
		var totalDuration time.Duration
		for i := 1; i < len(acTimes); i++ {
			totalDuration += acTimes[i].Sub(acTimes[i-1])
		}
		stats.AvgTimeBetweenAC = totalDuration / time.Duration(len(acTimes)-1)
	}

	// 获取最近的AC记录（最近30天）
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	for _, sub := range acSubmissions {
		if sub.Creation_time.After(thirtyDaysAgo) {
			stats.RecentACs = append(stats.RecentACs, sub)
		}
	}

	return stats
}
