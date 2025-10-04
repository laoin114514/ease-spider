package repository

import (
	"spider/config/db"
	"spider/src/models"
)

type LuoguRepository struct {
}

func NewLuoguRepository() *LuoguRepository {
	return &LuoguRepository{}
}

// 获取用户姓名和洛谷UID
func (r *LuoguRepository) GetUserNameMap() ([]map[string]string, error) {
	rows, err := db.Pool.Query("SELECT u.real_name,p.luogu_uid FROM user as u,oj_account as p WHERE u.id=p.user_id&&(u.role_id=1||u.role_id=3);")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userNameMaps := make([]map[string]string, 0)
	for rows.Next() {
		var realName, luoguUid string
		rows.Scan(&realName, &luoguUid)
		userNameMaps = append(userNameMaps, map[string]string{
			"real_name": realName,
			"luogu_uid": luoguUid,
		})
	}
	return userNameMaps, nil
}

// 获取用户提交记录数量
func (r *LuoguRepository) GetUserRecordsLength(uid string) (int, error) {
	var count int
	err := db.Pool.QueryRow("select count(*) from luogu_all_submissions where uid=?", uid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 批量插入用户提交记录
func (r *LuoguRepository) InsertUserRecords(record *models.LuoguRecordsResponse) error {
	result := record.CurrentData.Records.Result
	for _, record := range result {
		db.Pool.Exec("insert into luogu_all_submissions (uid,problem_id,problem_name,difficulty,is_pass,creation_time) values (?,?,?,?,?,?)", record.User.UID, record.Problem.PID, record.Problem.Title, record.Problem.Difficulty, record.Status == 12, record.SubmitTime)
	}
	return nil
}
