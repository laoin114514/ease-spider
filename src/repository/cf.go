package repository

import (
	"spider/config/db"
	"spider/src/models"
)

type CfRepository struct {
}

func NewCfRepository() *CfRepository {
	return &CfRepository{}
}

func (r *CfRepository) GetCfAccountData() ([]models.CfUserData, error) {
	rows, err := db.Pool.Query("SELECT p.cf_account,u.real_name FROM user as u,oj_account as p WHERE p.user_id=u.id and(role_id=1||role_id=3) and p.cf_account is not null and p.cf_account != ''")
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
