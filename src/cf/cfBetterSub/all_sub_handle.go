package cfBetterSub

import (
	"fmt"
	db2 "spider/config/db"
)

func allSubHandle(results []any) int {
	count := 0
	oldData := getOldData()
	for _, result := range results {
		result := result.(map[string]any)
		var table db2.Cf_all_submissions
		fillSubTable(result, &table)
		if oldData[table.Sub_id] {
			continue
		}
		err := db2.Insert_cf_all_sub(table)
		if err != nil {
			continue
		}
		count++
	}
	return count
}
func getOldData() map[int]bool {
	rows, err := db2.Pool.Query("select sub_id from cf_all_submissions where account=?", handle)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	has := map[int]bool{}
	for rows.Next() {
		var subId int
		rows.Scan(&subId)
		has[subId] = true
	}
	return has
}
