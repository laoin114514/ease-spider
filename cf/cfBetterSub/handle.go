package cfBetterSub

import (
	"fmt"
	"spider/db"
)

func neverPassHandle(results []any, count int) {
	passSet := map[string]bool{}
	fiilPassSet(results, passSet)
	neCount := insertNeverPass(results, passSet)
	deCount := deletePass(handle, passSet)
	fmt.Printf("%s提交新增%d 未过题新增%d 删除已过题%d\n", handle, count, neCount, deCount)
}
func insertNeverPass(results []any, passSet map[string]bool) int {
	count := 0
	insertArr := []db.Cf_never_pass{}
	has := map[string]bool{}
	countTime.Start()
	for _, result := range results {
		result := result.(map[string]any)
		var neverPassTable db.Cf_never_pass
		fiilNeverPass(result, &neverPassTable)
		if passSet[neverPassTable.ProblemId] {
			continue
		}
		if !has[neverPassTable.ProblemId] {
			insertArr = append(insertArr, neverPassTable)
			has[neverPassTable.ProblemId] = true
		}
	}
	for _, table := range insertArr {
		err := db.Insert_never_pass(table)
		if err != nil {
			break
		}
		count++
	}
	// countTime.End()
	return count
}
func deletePass(handle string, passSet map[string]bool) int {
	rows, _ := db.Pool.Query("select problemId from cf_never_pass where handle=?", handle)
	count := 0
	for rows.Next() {
		var problemId string
		rows.Scan(&problemId)
		if !passSet[problemId] {
			continue
		}
		db.Pool.Exec("delete from cf_never_pass where problemId=?", problemId)
		count++
	}
	return count
}
