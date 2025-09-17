package teamquestions

import (
	"fmt"
	"spider/config/db"
)

func fillTable(problem map[string]any, table *db.Cf_team_problems) {
	table.Team_contest_id = id
	table.Team_contest_name = name
	table.Official_contest_ID = fmt.Sprintf("%d", int(problem["contestId"].(float64))) + problem["index"].(string)
	table.Problem_name = problem["name"].(string)
	rating, has := problem["rating"].(float64)
	if !has {
		table.Rating = -1
	} else {
		table.Rating = int(rating)
	}
}
