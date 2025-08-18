package cfBetterSub

import (
	"fmt"
	"spider/db"
)

func Use() {
	dataBase := db.New()
	defer dataBase.Close()
	rows, _ := dataBase.Query("select account from user where role_id=1")
	defer rows.Close()
	for rows.Next() {
		var handle string
		rows.Scan(&handle)
		results, err := request(handle)
		if err != nil {
			fmt.Println(err)
			continue
		}
		count := 0
		for _, result := range results {
			result := result.(map[string]any)
			var table db.Cf_all_submissions
			handleArr := fillTable(result, &table)
			err := insert(table, dataBase)
			if err != nil {
				if len(handleArr) > 1 {
					continue
				}
				break
			}
			count++
		}
		fmt.Printf("%s新增：%d   ", handle, count)
		neverPassHandle(results, handle)
	}
}
