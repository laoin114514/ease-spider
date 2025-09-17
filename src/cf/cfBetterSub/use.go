package cfBetterSub

import (
	"fmt"
	"spider/component"
	"spider/config/db"
)

var countTime component.CountTime
var handle string
var username string
var outputs []string
var errs []string

func Use() {
	rows, err := db.Pool.Query("SELECT p.cf_account,u.real_name FROM user as u,oj_account as p WHERE p.user_id=u.id&&(role_id=1||role_id=3)")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		countTime.Start()
		rows.Scan(&handle, &username)
		results, err := request(handle)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%v", err))
			continue
		}
		subcount := allSubHandle(results)
		nePaCont, delCount := neverPassHandle(results)
		outputs = append(outputs, fmt.Sprintf("新增提交:%d 新增未过题:%d 删除未过题中已过题:%d  %s", subcount, nePaCont, delCount, username))
	}
	fmt.Println("cf提交情况:")
	for _, v := range outputs {
		fmt.Println(v)
	}
	for _, v := range errs {
		fmt.Println(v)
	}
	fmt.Printf("\n")
}
