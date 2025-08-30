package luogu

import (
	"fmt"
	"spider/db"
)

func Use() {
	fmt.Println("洛谷提交情况:")
	dataBase := db.New()
	defer dataBase.Close()
	rows, err := dataBase.Query("select id from luoguac")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var uid string
		rows.Scan(&uid)
		passProblems := request(uid)
		if len(passProblems) == 0 {
			continue
		}
		var table luoguAc
		table.id = uid
		count := handle(passProblems, table)
		todayCount := todayCount(dataBase, uid, count)
		err1 := insertPass(table, dataBase, todayCount)
		err2 := insertDayPass(table, dataBase, todayCount)
		if err1 != nil || err2 != nil {
			fmt.Println(err1, err2)
		}
		fmt.Printf("%s今日过题%d\n", uid, todayCount)
	}
}
