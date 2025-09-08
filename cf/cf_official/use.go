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
	for _, v := range result {
		v := v.(map[string]any)
		fillTable(v, &table)
		err := db.Insert_cf_official(table)
		if err != nil {
			break
		}
		count++
	}
	fmt.Println("官方题库新增", count, "\n")
}
