package newLuogu

import (
	"fmt"
	"time"
)

var difficultys = []string{"grey", "red", "brown", "yellow", "green", "blue", "purple", "black"}

func fillTable(obj map[string]any, table *luogu_all_submissions) {
	table.subId = fmt.Sprintf("%d", int(obj["id"].(float64)))
	problem := obj["problem"].(map[string]any)
	user := obj["user"].(map[string]any)
	difficulty := problem["difficulty"].(float64)
	table.difficulty = difficultys[int(difficulty)]
	if obj["status"].(float64) != 14 {
		table.isPass = false
	} else {
		table.isPass = true
	}
	table.problemName = problem["title"].(string)
	table.username = user["name"].(string)
	table.uid = fmt.Sprintf("%d", int(user["uid"].(float64)))
	stamp := time.Unix(int64(obj["submitTime"].(float64)), 1)
	table.subTime = stamp
	table.pid = problem["pid"].(string)
}
