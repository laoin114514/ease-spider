package teamquestions

import (
	"fmt"
	"spider/db"
)

func fillTable(problem map[string]any, table *db.Cf_team_question) {
	table.In_team_ID = id
	table.Contest_name = name
	table.Official_ID = fmt.Sprintf("%d", int(problem["contestId"].(float64))) + problem["index"].(string)
	table.Question_name = problem["name"].(string)
	rating, has := problem["rating"].(float64)
	if !has {
		table.Rating = -1
	} else {
		table.Rating = int(rating)
	}
}
