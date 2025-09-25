package team_contest_problems

import (
	"fmt"
	"spider/component"
	db2 "spider/config/db"
	"spider/src/cf/genUrl"
	"sync"
	"time"
)

var errorCount int
var contest genUrl.Contest
var count int
var officialProblems map[string]string
var countTime component.CountTime
var wg sync.WaitGroup

type chParm struct {
	id        int
	name      string
	prepareBy string
}

func Use(concurrency int) {
	countTime.Start()
	ch := make(chan chParm, concurrency)
	officialProblems = getOfficialProblems()

	rows, err := db2.Pool.Query("select contest_id,contest_name,prepare_by from cf_team_contests")
	if err != nil {
		return
	}
	defer rows.Close()

	parmArr := []chParm{}
	for rows.Next() {
		var id int
		var name, prepareBy string
		rows.Scan(&id, &name, &prepareBy)
		parmArr = append(parmArr, chParm{id: id, name: name, prepareBy: prepareBy})
	}

	// 构建并发函数
	for i := 0; i < concurrency && i < len(parmArr); i++ {
		go worker(ch)
	}

	// 发送所有任务
	for _, parm := range parmArr {
		ch <- parm
		time.Sleep(100 * time.Millisecond)
	}

	close(ch)
	wg.Wait()
	fmt.Printf("插入%d条训练赛题目  请求超频 %d 次   ", count, errorCount)
	countTime.End()
}

func worker(ch <-chan chParm) {
	for parm := range ch {
		wg.Add(1)
		problems := request(parm.id, parm.name, parm.prepareBy)
		for _, v := range problems {
			v := v.(map[string]any)
			v["id"], v["contest_name"] = parm.id, parm.name
			var table db2.Cf_team_problems
			fillTable(v, &table)
			err := db2.Insert_team_questions(table)
			if err != nil {
				continue
			}
			count++
		}
		wg.Done()
	}
}
