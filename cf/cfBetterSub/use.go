package cfBetterSub

import (
	"fmt"
	"spider/component"
	"spider/db"
)

var countTime component.CountTime
var handle string
var username string

func Use() {
	fmt.Println("cf提交情况:")
	rows, _ := db.Pool.Query("select account,username from user where role_id=1||role_id=3")
	defer rows.Close()
	for rows.Next() {
		countTime.Start()
		rows.Scan(&handle, &username)
		results, err := request(handle)
		if err != nil {
			fmt.Println(err)
			continue
		}
		subcount := allSubHandle(results)
		nePaCont, delCount := neverPassHandle(results)
		fmt.Printf("新增提交:%d 新增未过题:%d 删除未过题中已过题:%d  %s\n", subcount, nePaCont, delCount, username)
	}
}
