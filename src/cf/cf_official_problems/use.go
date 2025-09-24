package cf_official_problems

import (
	"fmt"
	db2 "spider/config/db"
	"spider/src/cf/genUrl"
)

var ProblemSet genUrl.ProblemSet

func Use() {
	url, err := ProblemSet.Problems()
	if err != nil {
		fmt.Println(err)
		return
	}
	result := reuquest(url)
	var table db2.Cf_official_problems
	count := 0
	has := getOldData()
	for _, v := range result {
		v := v.(map[string]any)
		fillTable(v, &table)
		if has[table.Problem_id] {
			continue
		}
		err := db2.Insert_cf_official(table)
		if err != nil {
			continue
		}
		count++
	}
	fmt.Println("官方题库新增", count)
}
func getOldData() map[string]bool {
	rows, err := db2.Pool.Query("select problem_id from cf_official_problems")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	has := map[string]bool{}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		has[id] = true
	}
	return has
}
