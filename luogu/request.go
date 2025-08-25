package luogu

import (
	"fmt"
	"spider/db"

	"github.com/go-resty/resty/v2"
)

func request(uid string) []db.Luogu_problem {
	c := resty.New()
	var result map[string]any
	c.R().
		SetResult(&result).
		Get(url(uid))
	code := result["code"].(float64)
	if code != 200 {
		fmt.Println("请求错误")
		return nil
	}
	passedProblems := result["currentData"].(map[string]any)["passedProblems"].([]any)
	var tables []db.Luogu_problem
	for _, p := range passedProblems {
		p := p.(map[string]any)
		var table db.Luogu_problem
		table.Uid = uid
		table.Pid = p["pid"].(string)
		table.Difficulty = int(p["difficulty"].(float64))
		table.Title = p["title"].(string)
		table.Type = p["type"].(string)
		tables = append(tables, table)
	}
	return tables
}
