package teamtainning

import (
	"github.com/go-resty/resty/v2"
)

func request(url string) []any {
	c := resty.New()

	var result map[string]any
	c.R().
		SetResult(&result).
		Get(url)
	arr := result["result"].([]any)
	return arr
}
