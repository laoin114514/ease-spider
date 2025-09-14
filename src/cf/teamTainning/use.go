package teamtainning

import (
	"fmt"
	"spider/component"
	db2 "spider/config/db"
	genUrl2 "spider/src/cf/genUrl"
)

var Contest genUrl2.Contest
var tempDB component.TempDB

func Use() {
	url := Contest.List(genUrl2.Contest_list{
		Handle:    "233zhang",
		Gym:       false,
		GroupCode: tempDB.Get("groupCode").(string),
	})
	result := request(url)
	count := 0
	has := getOldData()
	for _, v := range result {
		v := v.(map[string]any)
		var table db2.Cf_team_trainning
		fillTable(v, &table)
		if has[table.Id] {
			continue
		}
		err := db2.Insert_team_trainning(table)
		if err != nil {
			continue
		}
		count++
	}
	fmt.Printf("训练题单新增:%d\n", count)
}
func getOldData() map[int]bool {
	rows, err := db2.Pool.Query("select id from cf_team_training")
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
