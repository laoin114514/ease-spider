package cfBetterSub

import (
	"fmt"
	"spider/component"
	"spider/db"
)

func neverPassHandle(results []any, handle string, count int) {
	passSet := map[string]bool{}
	fiilPassSet(results, handle, passSet)
	neCount := insertNeverPass(results, handle, passSet)
	deCount := deletePass(handle, passSet)
	fmt.Printf("%s提交新增%d 未过题新增%d 删除已过题%d\n", handle, count, neCount, deCount)
}
func insertNeverPass(results []any, handle string, passSet map[string]bool) int {
	var countTime component.CountTime
	dateBase := db.New()
	defer dateBase.Close()
	count := 0
	insertArr := []db.Cf_never_pass{}
	has := map[string]bool{}
	countTime.Start()
	for _, result := range results {
		result := result.(map[string]any)
		var neverPassTable db.Cf_never_pass
		fiilNeverPass(result, handle, &neverPassTable)
		if passSet[neverPassTable.ProblemId] {
			continue
		}
		if !has[neverPassTable.ProblemId] {
			insertArr = append(insertArr, neverPassTable)
			has[neverPassTable.ProblemId] = true
		}
	}
	for _, table := range insertArr {
		_, err := dateBase.Exec(
			"insert into cf_never_pass (ProblemId,Handle,ProblemName,Rating ) values (?,?,?,?)",
			table.ProblemId,
			table.Handle,
			table.ProblemName,
			table.Rating,
		)
		if err != nil {
			break
		}
		count++
	}
	// countTime.End()
	return count
}
func deletePass(handle string, passSet map[string]bool) int {
	dateBase := db.New()
	defer dateBase.Close()
	rows, _ := dateBase.Query("select problemId from cf_never_pass where handle=?", handle)
	count := 0
	for rows.Next() {
		var problemId string
		rows.Scan(&problemId)
		if !passSet[problemId] {
			continue
		}
		dateBase.Exec("delete from cf_never_pass where problemId=?", problemId)
		count++
	}
	return count
}
