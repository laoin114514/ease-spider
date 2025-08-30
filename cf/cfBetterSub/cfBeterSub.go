package cfBetterSub

import (
	"fmt"
	"spider/component"
	"spider/db"
)

func Use() {
	fmt.Println("cf提交情况:")
	var countTime component.CountTime
	dataBase := db.New()
	defer dataBase.Close()
	rows, _ := dataBase.Query("select account from user where role_id=1||role_id=3")
	defer rows.Close()
	for rows.Next() {
		countTime.Start()
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
			handleArr := fillTable(result, handle, &table)
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
