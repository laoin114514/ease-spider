package cfBetterSub

import (
	db2 "spider/config/db"
)

func neverPassHandle(results []any) (int, int) {
	passSet := calPassSet(results)
	insertArr := calInsertArr(results, passSet)
	nePaCount := 0
	for _, table := range insertArr {
		err := db2.Insert_never_pass(table)
		if err != nil {
			continue
		}
		nePaCount++
	}
	delCount := deletePass(handle, passSet)
	return nePaCount, delCount
}

func calInsertArr(results []any, passSet map[string]bool) []db2.Cf_never_pass {
	//返回需要插入的从未通过的题目数组
	insertArr := []db2.Cf_never_pass{}
	has := map[string]bool{}
	for _, result := range results {
		result := result.(map[string]any)
		var neverPassTable db2.Cf_never_pass
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
		var table db2.Cf_all_submissions
		fillSubTable(result, &table)
		if table.Verdict == "OK" {
			passSet[table.ProblemId] = true
		}
	}
	return passSet
}

func deletePass(handle string, passSet map[string]bool) int {
	rows, _ := db2.Pool.Query("select problemId from cf_never_pass where handle=?", handle)
	count := 0
	for rows.Next() {
		var problemId string
		rows.Scan(&problemId)
		if !passSet[problemId] {
			continue
		}
		db2.Pool.Exec("delete from cf_never_pass where problemId=?", problemId)
		count++
	}
	return count
}
