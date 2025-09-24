package luogu

import (
	"fmt"
	"spider/component"
	db2 "spider/config/db"
	"time"
)

func Debug() {
	uidThis := "372213"
	result := request(uidThis, 1)
	page := countPage(result)
	count := 0
	//hasMap := getOldData(uidThis)
	for i := 1; i <= page; i++ {
		data := request(uidThis, i)
		result, ok := data["result"].([]any)
		if !ok {
			fmt.Println("返回错误")
			continue
		}
		for _, v := range result {
			v := v.(map[string]any)
			var table db2.Luogu_all_submissions
			fillTable(v, &table)
			//if !hasMap[table.Sub_id] {
			//	fmt.Println(table.Sub_id)
			//}
			err := db2.Insert_luogu_sub(table)
			if err != nil {
				stamp := time.Unix(int64(v["submitTime"].(float64)), 1)
				fmt.Println(err, "\n", component.DateTime(table.Creation_time.Unix()), component.DateTime(stamp.Unix()))
			}
		}
	}
	fmt.Println(count)
}
