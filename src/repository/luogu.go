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
func (r *LuoguRepository) GetUserNameMap() ([]models.LuoguUserDeliver, error) {
	rows, err := db.Pool.Query("SELECT u.real_name,p.luogu_uid FROM user as u,oj_account as p WHERE u.id=p.user_id&&(u.role_id=1||u.role_id=3) and p.luogu_uid is not null and p.luogu_uid != '';")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userNameMaps := make([]models.LuoguUserDeliver, 0)
	for rows.Next() {
		var realName, luoguUid string
		rows.Scan(&realName, &luoguUid)
		oldDataSet, err := r.GetUserOldDataSet(luoguUid)
		if err != nil {
			return nil, err
		}
		userNameMaps = append(userNameMaps, models.LuoguUserDeliver{
			RealName:   realName,
			Uid:        luoguUid,
			Count:      0,
			OldDataSet: oldDataSet,
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

func (r *LuoguRepository) GetUserOldDataSet(uid string) (map[string]bool, error) {
	rows, err := db.Pool.Query("select sub_id from luogu_all_submissions where uid=?", uid)
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

// 批量插入用户提交记录
func (r *LuoguRepository) InsertUserRecords(record *models.LuoguRecordsResponse) error {
	result := record.CurrentData.Records.Result
	for _, record := range result {
		db.Pool.Exec("insert into luogu_all_submissions (uid,problem_id,problem_name,difficulty,is_pass,creation_time) values (?,?,?,?,?,?)", record.User.UID, record.Problem.PID, record.Problem.Title, record.Problem.Difficulty, record.Status == 12, record.SubmitTime)
	}
	return nil
}

// 获取没有源代码的提交记录ID,仅获取role_id=3的用户（预备役）
func (r *LuoguRepository) GetSubidNoSourceCode() ([]string, error) {
	rows, err := db.Pool.Query("select sub_id from luogu_all_submissions as l,user as u,oj_account as o where NOT EXISTS(select 1 from luogu_source_code as l2 where l2.sub_id=l.sub_id ) AND l.uid=o.luogu_uid AND u.id=o.user_id AND u.role_id=3")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	subids := make([]string, 0)
	for rows.Next() {
		var subid string
		rows.Scan(&subid)
		subids = append(subids, subid)
	}
	return subids, nil
}
func (r *LuoguRepository) InsertSourceCode(subid string, source_code string) error {
	_, err := db.Pool.Exec("insert into luogu_source_code (sub_id, source_code) values (?,?)", subid, source_code)
	if err != nil {
		return err
	}
	return nil
}
func (r *LuoguRepository) GetNameBySubid(subid string) (string, error) {
	var name string
	err := db.Pool.QueryRow("SELECT u.real_name FROM user as u,oj_account as o,luogu_all_submissions as s WHERE u.id=o.user_id AND o.luogu_uid=s.uid AND s.sub_id=?", subid).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
func (r *LuoguRepository) GetProblemIdHasSourceCode() ([]string, error) {
	rows, err := db.Pool.Query("SELECT s.problem_id FROM luogu_source_code as c,luogu_all_submissions as s WHERE s.sub_id=c.sub_id AND c.source_code!='无' AND NOT EXISTS(SELECT 1 FROM luogu_solutions as l WHERE l.problem_id=s.problem_id) GROUP BY problem_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	problemIds := make([]string, 0)
	for rows.Next() {
		var problemId string
		rows.Scan(&problemId)
		problemIds = append(problemIds, problemId)
	}
	return problemIds, nil
}
func (r *LuoguRepository) ChangePrivateProblem() (int64, error) {
	result, err := db.Pool.Exec("UPDATE luogu_all_submissions SET difficulty='unknow' WHERE problem_id LIKE 'T%'")
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
