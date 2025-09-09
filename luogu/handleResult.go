package luogu

import (
	"fmt"
	"math"
	"spider/db"
)

func handle(result []any, count *int) error {
	for _, v := range result {
		v := v.(map[string]any)
		var table db.Luogu_all_submissions
		fillTable(v, &table)
		if has[table.SubId] {
			continue
		}
		err := db.Insert_luogu_sub(table)
		if err != nil {
			return err
		}
		*count++
	}
	return nil
}
func getOldData(uid string) map[string]bool {
	rows, err := db.Pool.Query("select subid from luogu_all_submissions where uid=?", uid)
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
	count := data["count"].(float64)
	perPage := data["perPage"].(float64)
	page := int(math.Ceil(count / perPage))
	return page
}
