package luogu

import (
	"fmt"
	"spider/component"
	"spider/config/db"
	"time"
)

var tempDB component.TempDB
var pass int
var has map[string]bool
var totalCount float64
var username string
var uid string
var outputs []string
var errs []string

func Use() {
	rows, err := db.Pool.Query("SELECT u.real_name,p.luogu_uid FROM user as u,oj_account as p WHERE u.id=p.user_id&&(u.role_id=1||u.role_id=3);")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		rows.Scan(&username, &uid)
		result := request(uid, 1)
		if uid == "" {
			errs = append(errs, username+"uid不存在")
			continue
		}
		if result == nil {
			errs = append(errs, username+"uid不存在")
			continue
		}
		page := countPage(result)
		has = getOldData(uid)
		count := loopRequest(page, false)
		if len(has) != int(totalCount)+count {
			errs = append(errs, fmt.Sprintf("%s少插入%d条 重新获取中...", username, int(totalCount)-len(has)))
			go loopRequest(page, true)
			time.Sleep(2 * time.Second)
			continue
		}
		outputs = append(outputs, fmt.Sprintf("新增提交%d %s", count, username))
	}
	fmt.Println("洛谷提交情况:")
	for _, v := range outputs {
		fmt.Println(v)
	}
	for _, v := range errs {
		fmt.Println(v)
	}
}
func loopRequest(page int, allCatch bool) int {
	uid1, username1 := uid, username
	hasMap := has
	count := 0
	for i := 1; i <= page; i++ {
		data := request(uid1, i)
		result, ok := data["result"].([]any)
		if !ok {
			continue
		}
		err := handle(result, &count, allCatch, hasMap)
		if err != nil && !allCatch {
			break
		}
	}
	if allCatch {
		fmt.Printf("洛谷重新插入%d条 %s\n", count, username1)
	}
	return count
}
