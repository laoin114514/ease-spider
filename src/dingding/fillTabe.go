package dingding

import (
	"spider/config/db"
	"time"
)

func fillTable(checkUpData map[string]any, table *db.DingCheckUp) {

	table.UserId = checkUpData["userId"].(string)
	table.Time = time.UnixMilli(int64(checkUpData["userCheckTime"].(float64)) + int64(8*3600*1000)).UTC()
	table.Name = obj[table.UserId].(string)
	_, ok := checkUpData["checkType"].(string)
	if !ok {
		table.CheckType = ""
		return
	}
	table.CheckType = checkUpData["checkType"].(string)
}
