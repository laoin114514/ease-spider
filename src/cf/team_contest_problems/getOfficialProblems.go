package team_contest_problems

import (
	"spider/config/db"
)

func getOfficialProblems() (arr map[string]string) {
	arr = make(map[string]string)
	rows, err := db.Pool.Query("select Problem_id,Title from cf_official_problems")
	if err != nil {
		return nil
	}
	for rows.Next() {
		var table db.Cf_official_problems
		rows.Scan(&table.Problem_id, &table.Title)
		arr[table.Title] = table.Problem_id
	}
	return arr
}
