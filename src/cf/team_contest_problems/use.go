package team_contest_problems

import (
	"fmt"
	db2 "spider/config/db"
	"spider/src/cf/genUrl"
)

var contest genUrl.Contest
var id int
var name string
var prepareBy string

func Use() {
	rows, err := db2.Pool.Query("select contest_id,contest_name,prepare_by from cf_team_contests")
	if err != nil {
		return
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		rows.Scan(&id, &name, &prepareBy)
		problems := request(id, name, prepareBy)
		for _, v := range problems {
			v := v.(map[string]any)
			var table db2.Cf_team_problems
			fillTable(v, &table)
			err := db2.Insert_team_questions(table)
			if err != nil {
				continue
			}
			count++
		}
	}
	fmt.Printf("插入%d条训练赛题目\n", count)
}
