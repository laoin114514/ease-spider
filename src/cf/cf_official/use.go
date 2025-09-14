package cfOfficial

import (
	"fmt"
	"spider/cf/genUrl"
	"spider/db"
)

var ProblemSet genUrl.ProblemSet

func Use() {
	url, err := ProblemSet.Problems()
	if err != nil {
		fmt.Println(err)
		return
	}
	result := reuquest(url)
	var table db.Cf_contest_official
	count := 0
	has := getOldData()
	for _, v := range result {
		v := v.(map[string]any)
		fillTable(v, &table)
		if has[table.Id] {
			continue
		}
		err := db.Insert_cf_official(table)
		if err != nil {
			continue
		}
		count++
	}
	fmt.Println("官方题库新增", count, "\n")
}
func getOldData() map[string]bool {
	rows, err := db.Pool.Query("select id from cf_contest_official")
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
