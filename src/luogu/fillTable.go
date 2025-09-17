package luogu

import (
	"errors"
	"fmt"
	"spider/config/db"
	"time"
)

var difficultys = []string{"grey", "red", "brown", "yellow", "green", "blue", "purple", "black"}

func fillTable(obj map[string]any, table *db.Luogu_all_submissions) error {
	table.Sub_id = fmt.Sprintf("%d", int(obj["id"].(float64)))
	problem := obj["problem"].(map[string]any)
	user, ok := obj["user"].(map[string]any)
	if !ok {
		return errors.New("不存在user")
	}
	difficulty := problem["difficulty"].(float64)
	table.Difficulty = difficultys[int(difficulty)]
	if obj["status"].(float64) != 12 {
		table.Is_pass = false
	} else {
		pass++
		table.Is_pass = true
	}
	table.Problem_name = problem["title"].(string)
	table.Uid = fmt.Sprintf("%d", int(user["uid"].(float64)))
	stamp := time.Unix(int64(obj["submitTime"].(float64)), 1)
	table.Creation_time = stamp
	table.Problem_id = problem["pid"].(string)
	return nil
}
