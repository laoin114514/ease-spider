package repository

import (
	"errors"
	"spider/config/db"
	"spider/src/models"
)

type CfRepository struct {
}

func NewCfRepository() *CfRepository {
	return &CfRepository{}
}

func (r *CfRepository) GetCfAccountData() ([]models.CfUserData, error) {
	rows, err := db.Pool.Query("SELECT o.cf_account,u.real_name from user as u,oj_account as o where u.id=o.user_id and cf_account !=' ' and cf_account is not null and cf_account!='' and (u.role_id=1 or u.role_id=3)")
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
func (r *CfRepository) GetCfRecordsInDbToset(account string) (map[int]bool, error) {
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
func (r *CfRepository) GetCfApikey(handle string) (string, string, error) {
	var apikey, secret string
	db.Pool.QueryRow("SELECT cf_apikey, cf_secret FROM oj_account WHERE cf_account = ?", handle).Scan(&apikey, &secret)
	if apikey == "" || secret == "" {
		return "", "", errors.New("apikey不存在")
	}
	return apikey, secret, nil
}
