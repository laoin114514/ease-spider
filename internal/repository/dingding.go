package repository

import "spider/config/db"

type DingdingRepository struct {
}

func NewDingdingRepository() *DingdingRepository {
	return &DingdingRepository{}
}

func (d *DingdingRepository) GetDingUserMap() (map[string]any, error) {
	rows, err := db.Pool.Query("SELECT real_name,ding_id FROM user WHERE ding_id IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dingUserMap := make(map[string]any)
	for rows.Next() {
		var realName, dingId string
		rows.Scan(&realName, &dingId)
		dingUserMap[dingId] = realName
	}
	return dingUserMap, nil
}
