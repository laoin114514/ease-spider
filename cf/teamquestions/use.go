package teamquestions

import (
	"fmt"
	"spider/cf/genUrl"
	"spider/db"
)

var contest genUrl.Contest
var id int
var name string
var prepareBy string

func Use() {
	rows, err := db.Pool.Query("select id,name,prepareBy from cf_team_training")
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
			var table db.Cf_team_question
			fillTable(v, &table)
			err := db.Insert_team_questions(table)
			if err != nil {
				continue
			}
			count++
		}
	}
	fmt.Printf("插入%d条训练赛题目\n", count)
}
