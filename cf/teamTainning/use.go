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
	for _, v := range result {
		v := v.(map[string]any)
		var table db.Cf_team_trainning
		fillTable(v, &table)
		err := db.Insert_team_trainning(table)
		if err != nil {
			break
		}
		count++
	}
	fmt.Printf("训练题单新增:%d\n", count)
}
