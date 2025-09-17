package dingding

import (
	"spider/config/db"
	"time"
)

func fillTable(checkUpData map[string]any, table *db.Ding_checkUp) {

	table.Ding_id = checkUpData["userId"].(string)
	table.Time = time.UnixMilli(int64(checkUpData["userCheckTime"].(float64)) + int64(8*3600*1000)).UTC()
	table.Name = obj[table.Ding_id].(string)
	_, ok := checkUpData["checkType"].(string)
	if !ok {
		table.Check_type = ""
		return
	}
	table.Check_type = checkUpData["checkType"].(string)
}
