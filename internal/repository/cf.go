package repository

import (
	"spider/config/db"
	"spider/internal/models"
	"time"
)

type CfRepository struct {
}

func NewCfRepository() *CfRepository {
	return &CfRepository{}
}

func (r *CfRepository) GetCfAccountData() ([]models.CfUserData, error) {
	rows, err := db.Pool.Query("SELECT u.cf_account,u.real_name from user as u where u.cf_account !=' ' and u.cf_account is not null and u.cf_account!='' and (u.role_id=1 or u.role_id=3)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cfUserDatas := make([]models.CfUserData, 0)
	for rows.Next() {
		var account, realName string
		err := rows.Scan(&account, &realName)
		if err != nil {
			return nil, err
		}
		cfUserDatas = append(cfUserDatas, models.CfUserData{Account: account, RealName: realName})
	}
	return cfUserDatas, nil
}
func (r *CfRepository) GetCfRecordsIdToset(account string) (map[int]bool, error) {
	rows, err := db.Pool.Query("SELECT sub_id FROM cf_all_submissions WHERE account=?", account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cfRecords := make(map[int]bool)
	for rows.Next() {
		var subId int
		err := rows.Scan(&subId)
		if err != nil {
			return nil, err
		}
		cfRecords[subId] = true
	}
	return cfRecords, nil
}

// CfApiKey 是单个 CF 账号的 API 凭据（对应 user 表的 cf_apikey / cf_secret）。
type CfApiKey struct {
	ApiKey string
	Secret string
}

func (r *CfRepository) GetCfApikeyPool() (map[string]CfApiKey, error) {
	rows, err := db.Pool.Query("SELECT cf_account, cf_apikey, cf_secret FROM user WHERE cf_account != '' AND cf_account IS NOT NULL AND cf_account != ' ' AND (role_id = 1 OR role_id = 3) AND cf_apikey IS NOT NULL AND cf_secret IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	apiKeyPool := make(map[string]CfApiKey)
	for rows.Next() {
		var account, apikey, secret string
		err := rows.Scan(&account, &apikey, &secret)
		if err != nil {
			return nil, err
		}
		apiKeyPool[account] = CfApiKey{ApiKey: apikey, Secret: secret}
	}
	return apiKeyPool, nil
}
func (r *CfRepository) GetTeamContests() ([]db.Cf_team_contests, error) {
	rows, err := db.Pool.Query("SELECT Contest_id, Contest_name, Start_time, PrePare_by FROM cf_team_contests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cfTeamContests := make([]db.Cf_team_contests, 0)
	for rows.Next() {
		var id int
		var name, prepareBy string
		var startTime time.Time
		rows.Scan(&id, &name, &startTime, &prepareBy)
		cfTeamContests = append(cfTeamContests, db.Cf_team_contests{
			Contest_id:   id,
			Contest_name: name,
			PrePare_by:   prepareBy,
			Start_time:   startTime,
		})
	}
	return cfTeamContests, nil
}
func (r *CfRepository) GetTeamOfficailId() ([]int, error) {
	rows, err := db.Pool.Query("SELECT official_contest_id from cf_team_problem")
	if err != nil {
		return []int{}, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}
