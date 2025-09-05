package newLuogu

import (
	"fmt"
	"spider/component"

	"github.com/go-resty/resty/v2"
)

var countTime component.CountTime
var tempDB component.TempDB

func request(uid string, page int) map[string]any {
	cookie := tempDB.Get("Cookie").(string)
	c := resty.New()
	var result map[string]any
	_, err := c.R().
		SetHeader("Cookie", cookie).
		SetResult(&result).
		Get(fmt.Sprintf("https://www.luogu.com.cn/record/list?user=%s&page=%d&_contentOnly=1", uid, page))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if result["code"].(float64) != 200 {
		fmt.Println(result["code"])
		return nil
	}
	return result["currentData"].(map[string]any)["records"].(map[string]any)
}
