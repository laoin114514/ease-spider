package luogu

import (
	"fmt"
	"spider/config/db"
	"time"
)

var difficultys = []string{"grey", "red", "brown", "yellow", "green", "blue", "purple", "black"}

func fillTable(obj map[string]any, table *db.Luogu_all_submissions) {
	table.SubId = fmt.Sprintf("%d", int(obj["id"].(float64)))
	problem := obj["problem"].(map[string]any)
	user := obj["user"].(map[string]any)
	difficulty := problem["difficulty"].(float64)
	table.Difficulty = difficultys[int(difficulty)]
	if obj["status"].(float64) != 12 {
		table.IsPass = false
	} else {
		pass++
		table.IsPass = true
	}
	table.ProblemName = problem["title"].(string)
	table.Username = user["name"].(string)
	table.Uid = fmt.Sprintf("%d", int(user["uid"].(float64)))
	stamp := time.Unix(int64(obj["submitTime"].(float64)), 1)
	table.SubTime = stamp
	table.Pid = problem["pid"].(string)
}
