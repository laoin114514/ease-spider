package luogu

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func request(uid string, page int) map[string]any {
	cookie := tempDB.Get("Cookie").(string)
	c := resty.New()
	var result map[string]any
	url := fmt.Sprintf("https://www.luogu.com.cn/record/list?user=%s&page=%d&_contentOnly=1", uid, page)
	_, err := c.R().
		SetHeader("Cookie", cookie).
		SetResult(&result).
		Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if result == nil {
		return nil
	}
	if result["code"].(float64) != 200 {
		return nil
	}
	return result["currentData"].(map[string]any)["records"].(map[string]any)
}
