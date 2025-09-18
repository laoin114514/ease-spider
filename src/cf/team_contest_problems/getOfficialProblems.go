package team_contest_problems

import "spider/config/db"

func getOfficialProblems() (arr []db.Cf_official_problems) {
	rows, err := db.Pool.Query("select * from cf_official_problems")
	if err != nil {
		return nil
	}
	for rows.Next() {
		var table db.Cf_official_problems
		rows.Scan(&table.Problem_id, &table.Title)
		arr = append(arr, table)
	}
	return arr
}
