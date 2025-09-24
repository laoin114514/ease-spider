package teamtainning

import (
	"spider/config/db"
	"time"
)

func fillTable(result map[string]any, table *db.Cf_team_contests) {
	table.Contest_id = int(result["id"].(float64))
	table.Contest_name = result["name"].(string)
	table.PrePare_by = result["preparedBy"].(string)
	start := int64(result["startTimeSeconds"].(float64))
	table.Start_time = time.Unix(start+8*3600, 1)
}
