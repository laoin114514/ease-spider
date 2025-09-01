package newLuogu

import (
	"fmt"
	"time"
)

func fillTable(obj map[string]any, table *luogu_all_submissions) {
	table.subId = fmt.Sprintf("%d", int(obj["id"].(float64)))
	problem := obj["problem"].(map[string]any)
	user := obj["user"].(map[string]any)
	difficulty := problem["difficulty"].(float64)
	if difficulty == 0 {
		table.difficulty = "grey"
	} else if difficulty == 1 {
		table.difficulty = "red"
	} else if difficulty == 2 {
		table.difficulty = "brown"
	} else if difficulty == 3 {
		table.difficulty = "yellow"
	} else if difficulty == 4 {
		table.difficulty = "green"
	} else if difficulty == 5 {
		table.difficulty = "blue"
	} else if difficulty == 6 {
		table.difficulty = "purple"
	} else if difficulty == 7 {
		table.difficulty = "black"
	}
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
