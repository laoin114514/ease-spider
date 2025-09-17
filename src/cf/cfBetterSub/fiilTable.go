package cfBetterSub

import (
	"fmt"
	"spider/config/db"
	"time"
)

func fillSubTable(result map[string]any, table *db.Cf_all_submissions) {
	table.Account = handle
	table.Sub_id = int(result["id"].(float64))
	table.Problem_id = fmt.Sprintf("%d", int(result["problem"].(map[string]any)["contestId"].(float64))) + result["problem"].(map[string]any)["index"].(string)
	table.Problem_name = result["problem"].(map[string]any)["name"].(string)
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
	table.Creation_time = time.Unix(int64(t), 0)
}
func fiilNeverPass(result map[string]any, table *db.Cf_never_pass) {
	table.Problem_id = fmt.Sprintf("%d", int(result["problem"].(map[string]any)["contestId"].(float64))) + result["problem"].(map[string]any)["index"].(string)
	table.Account = handle
	table.Problem_name = result["problem"].(map[string]any)["name"].(string)
	rating := result["problem"].(map[string]any)["rating"]
	if rating != nil {
		table.Rating = int(rating.(float64))
	} else {
		table.Rating = -1
	}
}
