package cfBetterSub

import (
	"fmt"
	"spider/db"
	"time"
)

func fillTable(result map[string]any, handle string, table *db.Cf_all_submissions) {
	table.Handle = handle
	table.SubId = int(result["id"].(float64))
	table.ProblemId = fmt.Sprintf("%d", int(result["problem"].(map[string]any)["contestId"].(float64))) + result["problem"].(map[string]any)["index"].(string)
	table.ProblemName = result["problem"].(map[string]any)["name"].(string)
	rating := result["problem"].(map[string]any)["rating"]
	if rating != nil {
		table.Rating = int(rating.(float64))
	} else {
		table.Rating = -1
	}
	_, ok := result["verdict"].(string)
	if !ok {
		table.Verdict = ""
	} else {
		table.Verdict = result["verdict"].(string)
	}
	t := int(result["creationTimeSeconds"].(float64))
	table.CreationTime = time.Unix(int64(t), 0)
}
func fiilNeverPass(result map[string]any, handle string, table *db.Cf_never_pass) {
	table.ProblemId = fmt.Sprintf("%d", int(result["problem"].(map[string]any)["contestId"].(float64))) + result["problem"].(map[string]any)["index"].(string)
	table.Handle = handle
	table.ProblemName = result["problem"].(map[string]any)["name"].(string)
	rating := result["problem"].(map[string]any)["rating"]
	if rating != nil {
		table.Rating = int(rating.(float64))
	} else {
		table.Rating = -1
	}
}
func fiilPassSet(results []any, handle string, passSet map[string]bool) {
	for _, result := range results {
		result := result.(map[string]any)
		var table db.Cf_all_submissions
		fillTable(result, handle, &table)
		if table.Verdict == "OK" {
			passSet[table.ProblemId] = true
		}
	}
}
