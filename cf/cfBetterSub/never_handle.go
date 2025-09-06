package cfBetterSub

import (
	"spider/db"
)

func neverPassHandle(results []any) (int, int) {
	passSet := calPassSet(results)
	insertArr := calInsertArr(results, passSet)
	nePaCount := 0
	for _, table := range insertArr {
		err := db.Insert_never_pass(table)
		if err != nil {
			break
		}
		nePaCount++
	}
	delCount := deletePass(handle, passSet)
	return nePaCount, delCount
}

func calInsertArr(results []any, passSet map[string]bool) []db.Cf_never_pass {
	//返回需要插入的从未通过的题目数组
	insertArr := []db.Cf_never_pass{}
	has := map[string]bool{}
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
	return insertArr
}

func calPassSet(results []any) map[string]bool {
	passSet := map[string]bool{}
	for _, result := range results {
		result := result.(map[string]any)
		var table db.Cf_all_submissions
		fillSubTable(result, &table)
		if table.Verdict == "OK" {
			passSet[table.ProblemId] = true
		}
	}
	return passSet
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
