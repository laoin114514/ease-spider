package request

import (
	"spider/src/utils"

	"github.com/go-resty/resty/v2"
)

var JsonDB *utils.JSONDB

func Init() {
	c := RestyInit()
	resp, _ := c.R().Get("https://www.luogu.com.cn/auth/login")
	cookie := resp.Cookies()
	JsonDB = utils.NewJSONDB("tempDB.json")
	JsonDB.Set("cookie2", cookie[0].Name+"="+cookie[0].Value)
}

func InitRedirect() {
	c := RestyInit()
	resp, _ := c.R().
		SetHeader("Cookie", JsonDB.Get("cookie2").(string)).
		Get("https://www.luogu.com.cn/auth/login")
	cookie := resp.Cookies()
	JsonDB.Set("cookie1", cookie[0].Name+"="+cookie[0].Value)
}

func RestyInit() *resty.Client {
	c := resty.New()
	c.SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0")
	return c
}
