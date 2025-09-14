package cfBetterSub

import (
	"fmt"
	"spider/db"
)

func allSubHandle(results []any) int {
	count := 0
	oldData := getOldData()
	for _, result := range results {
		result := result.(map[string]any)
		var table db.Cf_all_submissions
		fillSubTable(result, &table)
		if oldData[table.SubId] {
			continue
		}
		err := db.Insert_cf_all_sub(table)
		if err != nil {
			continue
		}
		count++
	}
	return count
}
func getOldData() map[int]bool {
	rows, err := db.Pool.Query("select subId from cf_all_submissions where handle=?", handle)
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
