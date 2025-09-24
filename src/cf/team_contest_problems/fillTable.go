package team_contest_problems

import (
	"spider/config/db"
)

func fillTable(problem map[string]any, table *db.Cf_team_problems) {
	table.Team_contest_id = problem["id"].(int)
	table.Team_contest_name = problem["contest_name"].(string)
	table.Problem_name = problem["name"].(string)
	table.Official_contest_ID = officialProblems[table.Problem_name]
	rating, has := problem["rating"].(float64)
	if !has {
		table.Rating = -1
	} else {
		table.Rating = int(rating)
	}
}
