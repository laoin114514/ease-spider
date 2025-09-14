package teamtainning

import (
	"spider/config/db"
	"time"
)

func fillTable(result map[string]any, table *db.Cf_team_trainning) {
	table.Id = int(result["id"].(float64))
	table.Name = result["name"].(string)
	table.PrePareBy = result["preparedBy"].(string)
	start := int64(result["startTimeSeconds"].(float64))
	table.StartTime = time.Unix(start, 1)
}
