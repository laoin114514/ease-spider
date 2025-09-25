package luogu

import (
	"fmt"
	"spider/component"
	"spider/config/db"
	"sync"
)

var tempDB component.TempDB
var pass int
var outputs []string
var errs []string

type parm struct {
	username   string
	uid        string
	has        map[string]bool
	totalCount float64
}

var wg sync.WaitGroup

func initData() {
	pass = 0
	outputs = []string{}
	errs = []string{}
}
func Use(concurrency int) {
	var countTime component.CountTime
	countTime.Start()
	ch := make(chan parm, concurrency)
	initData()
	rows, err := db.Pool.Query("SELECT u.real_name,p.luogu_uid FROM user as u,oj_account as p WHERE u.id=p.user_id&&(u.role_id=1||u.role_id=3);")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	var parms []parm
	for rows.Next() {
		var username string
		var uid string
		rows.Scan(&username, &uid)
		has := getOldData(uid)
		parms = append(parms, parm{username: username, uid: uid, has: has})
	}
	//构建并发函数
	for i := 0; i < concurrency && i < len(parms); i++ {
		go worker(ch)
	}
	//发布所有任务
	for _, v := range parms {
		ch <- v
	}
	close(ch)
	wg.Wait()
	content := ""
	dateTime := component.NowDateTime()
	content += "============================ " + dateTime + " ============================\n"
	for _, v := range outputs {
		content += v + "\n"
	}
	fmt.Println()
	for _, v := range errs {
		content += v + "\n"
	}
	timeSum := 0
	for _, v := range reqTimes {
		timeSum += v
	}
	content += fmt.Sprintf("请求平均时长 %dms\n", timeSum/len(reqTimes))
	content += countTime.EndWithStr()
	err1 := component.CoverFile("log/luogu.log", content)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println("获取洛谷过题成功")
}
func worker(ch <-chan parm) {
	for c := range ch {
		wg.Add(1)
		if c.uid == "" {
			errs = append(errs, c.username+"uid不存在")
			wg.Done()
			continue
		}
		result := request(c.uid, 1)
		if result == nil {
			errs = append(errs, c.username+"uid不存在")
			wg.Done()
			continue
		}
		page := countPage(result, &c)
		count := loopRequest(page, false, c)
		if len(c.has)+count != int(c.totalCount) {
			errs = append(errs, fmt.Sprintf("%s少插入%d条 重新获取中...", c.username, int(c.totalCount)-len(c.has)))
			loopRequest(page, true, c)
			wg.Done()
			continue
		}
		outputs = append(outputs, fmt.Sprintf("新增提交%d %s", count, c.username))
		wg.Done()
	}
}
func loopRequest(page int, allCatch bool, parm parm) int {
	count := 0
	for i := 1; i <= page; i++ {
		data := request(parm.uid, i)
		result, ok := data["result"].([]any)
		if !ok {
			continue
		}
		err := handle(result, &count, allCatch, parm)
		if err != nil && !allCatch {
			break
		}
	}
	if allCatch {
		fmt.Println(fmt.Sprintf("洛谷重新插入%d条 %s", count, parm.username))
	}
	return count
}
