package luogu

import (
	"math"
	"spider/db"
)

func handle(result []any) error {
	for _, v := range result {
		v := v.(map[string]any)
		var table db.Luogu_all_submissions
		fillTable(v, &table)
		err := db.Insert_luogu_sub(table)
		if err != nil {
			return err
		}
	}
	return nil
}
func countPage(data map[string]any) int {
	count := data["count"].(float64)
	perPage := data["perPage"].(float64)
	page := int(math.Ceil(count / perPage))
	return page
}
