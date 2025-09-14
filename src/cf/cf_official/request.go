package cfOfficial

import "github.com/go-resty/resty/v2"

func reuquest(url string) []any {
	c := resty.New()
	var data map[string]any
	c.R().
		SetResult(&data).
		Get(url)
	result := data["result"].(map[string]any)["problems"].([]any)
	return result
}
