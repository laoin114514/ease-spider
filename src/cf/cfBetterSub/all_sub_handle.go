package cfBetterSub

import (
	"fmt"
	db2 "spider/config/db"
)

func allSubHandle(results []any, parm chanParm) int {
	count := 0
	oldData := getOldData(parm)
	for _, result := range results {
		result := result.(map[string]any)
		var table db2.Cf_all_submissions
		fillSubTable(result, &table, parm)
		if oldData[table.Sub_id] {
			continue
		}
		err := db2.Insert_cf_all_sub(table)
		if err != nil {
			fmt.Println(err)
			continue
		}
		count++
	}
	return count
}
func getOldData(parm chanParm) map[int]bool {
	rows, err := db2.Pool.Query("select sub_id from cf_all_submissions where account=?", parm.handle)
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
