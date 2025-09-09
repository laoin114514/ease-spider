package teamtainning

import (
	"fmt"
	"spider/cf/genUrl"
	"spider/component"
	"spider/db"
)

var Contest genUrl.Contest
var tempDB component.TempDB

func Use() {
	url := Contest.List(genUrl.Contest_list{
		Handle:    "233zhang",
		Gym:       false,
		GroupCode: tempDB.Get("groupCode").(string),
	})
	result := request(url)
	count := 0
	has := getOldData()
	for _, v := range result {
		v := v.(map[string]any)
		var table db.Cf_team_trainning
		fillTable(v, &table)
		if has[table.Id] {
			continue
		}
		err := db.Insert_team_trainning(table)
		if err != nil {
			continue
		}
		count++
	}
	fmt.Printf("训练题单新增:%d\n", count)
}
func getOldData() map[int]bool {
	rows, err := db.Pool.Query("select id from cf_team_training")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	has := map[int]bool{}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		has[id] = true
	}
	return has
}
