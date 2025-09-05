package dingding

import (
	"spider/db"
	"time"
)

func fillTable(checkUpData map[string]any, table *db.DingCheckUp) {
	userObject := tempDB.Get("dingUserId").(map[string]any)
	table.UserId = checkUpData["userId"].(string)
	table.Time = time.UnixMilli(int64(checkUpData["userCheckTime"].(float64)))
	table.Name = userObject[table.UserId].(string)
	table.CheckType = checkUpData["checkType"].(string)
}
