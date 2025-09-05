package cfBetterSub

import (
	"fmt"
	"spider/component"
	"spider/db"
)

var countTime component.CountTime
var handle string

func Use() {
	fmt.Println("cf提交情况:")
	rows, _ := db.Pool.Query("select account from user where role_id=1||role_id=3")
	defer rows.Close()
	for rows.Next() {
		countTime.Start()
		rows.Scan(&handle)
		results, err := request(handle)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func() {
			count := 0
			for _, result := range results {
				result := result.(map[string]any)
				var table db.Cf_all_submissions
				fillTable(result, handle, &table)
				err := db.Insert_cf_all_sub(table)
				if err != nil {
					break
				}
				count++
			}
			go neverPassHandle(results, handle, count)
		}()
	}
}
