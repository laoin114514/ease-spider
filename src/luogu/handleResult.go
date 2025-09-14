package luogu

import (
	"errors"
	"fmt"
	"math"
	db2 "spider/config/db"
)

func handle(result []any, count *int, allCatch bool, hasMap map[string]bool) error {
	for _, v := range result {
		v := v.(map[string]any)
		var table db2.Luogu_all_submissions
		fillTable(v, &table)
		if hasMap[table.SubId] && !allCatch {
			return errors.New("重复")
		} else if hasMap[table.SubId] && allCatch {
			continue
		}
		db2.Insert_luogu_sub(table)
		*count++
	}
	return nil
}
func getOldData(uid string) map[string]bool {
	rows, err := db2.Pool.Query("select subid from luogu_all_submissions where uid=?", uid)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	has := map[string]bool{}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		has[id] = true
	}
	return has
}
func countPage(data map[string]any) int {
	totalCount = data["count"].(float64)
	perPage := data["perPage"].(float64)
	page := int(math.Ceil(totalCount / perPage))
	return page
}
