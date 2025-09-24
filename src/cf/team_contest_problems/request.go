package team_contest_problems

import (
	"spider/src/cf/genUrl"

	"github.com/go-resty/resty/v2"
)

func request(id int, name string, prepareBy string) []any {
	c := resty.New()
	url := contest.Standings(genUrl.Contest_standings{
		Handle:    prepareBy,
		ContestId: id,
		AsManager: true,
		From:      1,
		Count:     50000,
	})
	var result map[string]any
	resp, err := c.R().
		SetResult(&result).
		Get(url)
	status := resp.StatusCode()
	if status == 429 {
		errorCount++
	}
	if err != nil {
		return []any{}
	}
	if result["status"] != "OK" {
		return []any{}
	}
	arr := result["result"].(map[string]any)["problems"].([]any)
	return arr
}
