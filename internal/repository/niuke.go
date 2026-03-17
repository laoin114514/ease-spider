package repository

import (
	"spider/config/db"
	"spider/internal/models"
)

type NiukeRepository struct {
}

func NewNiukeRepository() *NiukeRepository {
	return &NiukeRepository{}
}

// GetUserNameMap 获取用户姓名和牛客UID
// 从user表获取有牛客UID的用户
func (r *NiukeRepository) GetUserNameMap() ([]models.NiukeUserDeliver, error) {
	rows, err := db.Pool.Query(`
		SELECT real_name, niuke_uid 
		FROM user 
		WHERE (role_id=1 or role_id=3) 
		AND niuke_uid IS NOT NULL 
		AND niuke_uid != '';
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userNameMaps := make([]models.NiukeUserDeliver, 0)
	for rows.Next() {
		var realName, niukeUid string
		rows.Scan(&realName, &niukeUid)

		// 获取该用户已爬取的提交ID集合（用于增量更新）
		oldDataSet, err := r.GetUserOldDataSet(niukeUid)
		if err != nil {
			return nil, err
		}

		userNameMaps = append(userNameMaps, models.NiukeUserDeliver{
			RealName:   realName,
			Uid:        niukeUid,
			Count:      0,
			OldDataSet: oldDataSet,
		})
	}
	return userNameMaps, nil
}

// GetUserRecordsLength 获取用户提交记录数量
func (r *NiukeRepository) GetUserRecordsLength(uid string) (int, error) {
	var count int
	err := db.Pool.QueryRow("select count(*) from niuke_submissions where uid=?", uid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetUserOldDataSet 获取用户已爬取的提交ID集合
func (r *NiukeRepository) GetUserOldDataSet(uid string) (map[string]bool, error) {
	rows, err := db.Pool.Query("select sub_id from niuke_submissions where uid=?", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	oldDataSet := make(map[string]bool)
	for rows.Next() {
		var subId string
		rows.Scan(&subId)
		oldDataSet[subId] = true
	}
	return oldDataSet, nil
}

// InsertUserRecords 批量插入用户提交记录
func (r *NiukeRepository) InsertUserRecords(uid string, records []models.NiukeSubmission) error {
	for _, record := range records {
		// 生成唯一提交ID：uid + problem_id + submit_time
		subId := uid + "_" + record.ProblemId + "_" + record.SubmitTime

		_, err := db.Pool.Exec(
			"INSERT INTO niuke_submissions (sub_id, uid, user_name, problem_id, problem_name, submit_time, status, language) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE updated_at=CURRENT_TIMESTAMP",
			subId,
			uid,
			record.UserName,
			record.ProblemId,
			record.ProblemName,
			record.SubmitTime,
			record.Status,
			record.Language,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
