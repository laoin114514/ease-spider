package luogu

import (
	"fmt"
	"spider/db"
)

func Use() {
	dataBase := db.New()
	defer dataBase.Close()
	rows, err := dataBase.Query("select luogu from platform_id where luogu is not null")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var uid string
		rows.Scan(&uid)
		tables := request(uid)
		count := 0
		for _, table := range tables {
			err := insertPass(table, dataBase)
			if err != nil {
				fmt.Println(err)
				break
			}
			count++
		}
		fmt.Printf("新增过题%d", count)
	}
}
