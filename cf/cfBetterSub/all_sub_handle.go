package cfBetterSub

import "spider/db"

func allSubHandle(results []any) int {
	count := 0
	for _, result := range results {
		result := result.(map[string]any)
		var table db.Cf_all_submissions
		fillSubTable(result, &table)
		err := db.Insert_cf_all_sub(table)
		if err != nil {
			break
		}
		count++
	}
	return count
}
