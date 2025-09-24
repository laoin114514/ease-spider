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
	resp, err := c.R().
		SetHeader("Cookie", cookie).
		SetResult(&result).
		Get(url)
	if resp.StatusCode() != 200 {
		errs = append(errs, fmt.Sprintf("请求失败  %d", resp.StatusCode()))
		return nil
	}
	if result["code"].(float64) != 200 {
		errs = append(errs, fmt.Sprintf("请求失败  %.0f", result["code"]))
		return nil
	}
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
