package dingding

import "time"

func fillTable(checkUpData map[string]any, table *dingCheckUp) {
	userObject := getUserId()
	table.userId = checkUpData["userId"].(string)
	table.time = time.UnixMilli(int64(checkUpData["userCheckTime"].(float64)))
	table.name = userObject[table.userId].(string)
	table.checkType = checkUpData["checkType"].(string)
}
