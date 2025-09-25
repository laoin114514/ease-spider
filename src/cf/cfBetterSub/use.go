package cfBetterSub

import (
	"fmt"
	"spider/component"
	"spider/config/db"
	"sync"
)

// 注意这些数据都要初始化
var outputs []string
var errs []string
var wg sync.WaitGroup
var errCount int
var content string
var countTime component.CountTime

type chanParm struct {
	username string
	handle   string
}

func initData() {
	outputs = []string{}
	errs = []string{}
	errCount = 0
	countTime.Start()
}

func Use(concurrency int) {
	initData()
	ch := make(chan chanParm, concurrency-1)
	rows, err := db.Pool.Query("SELECT p.cf_account,u.real_name FROM user as u,oj_account as p WHERE p.user_id=u.id&&(role_id=1||role_id=3)")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	var parmArr []chanParm
	for rows.Next() {
		var handle string = ""
		var username string = ""
		rows.Scan(&handle, &username)
		parmArr = append(parmArr, chanParm{handle: handle, username: username})
	}

	//构建并发器
	for i := 0; i < concurrency && i < len(parmArr); i++ {
		go worker(ch)
	}

	//发布任务
	for _, parm := range parmArr {
		ch <- parm
	}
	wg.Wait()

	content = "====================================" + component.NowDateTime() + "===========================================\n"
	for _, v := range outputs {
		content += v + "\n"
	}
	for _, v := range errs {
		content += v + "\n"
	}
	content += fmt.Sprintf("超频次数%d\n", errCount)
	content += countTime.EndWithStr()
	component.CoverFile("log/cf.log", content)
	fmt.Println("cf获取成功")
}
func worker(ch <-chan chanParm) {
	for c := range ch {
		wg.Add(1)
		results, err := request(c)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%v", err))
			wg.Done()
			continue
		}
		subcount := allSubHandle(results, c)
		nePaCont, delCount := neverPassHandle(results, c)
		outputs = append(outputs, fmt.Sprintf("新增提交:%d 新增未过题:%d 删除未过题中已过题:%d  %s", subcount, nePaCont, delCount, c.username))
		wg.Done()
	}
}
