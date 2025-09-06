package cfOfficial

import (
	"fmt"
	"spider/db"
)

func fillTable(result map[string]any, table *db.Cf_contest_official) {
	table.Id = fmt.Sprintf("%d%s", int(result["contestId"].(float64)), result["index"])
	table.Title = result["name"].(string)
	if result["points"] == nil {
		table.Points = -1
	} else {
		table.Points = int(result["points"].(float64))
	}
	if result["rating"] == nil {
		table.Rating = -1
	} else {
		table.Rating = int(result["rating"].(float64))
	}
	table.Tags = arrToString(result["tags"].([]any))
}
func arrToString(arr []any) string {
	str := "["
	for _, v := range arr {
		v := v.(string)
		str += "\"" + v + "\"" + ","
	}
	str = str[:len(str)-1] + "]"
	return str
}
