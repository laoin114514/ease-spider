package request

import (
	"spider/component"

	"github.com/go-resty/resty/v2"
)

var tempDB component.TempDB

func Init() {
	c := RestyInit()
	resp, _ := c.R().Get("https://www.luogu.com.cn/auth/login")
	cookie := resp.Cookies()
	tempDB.Set("cookie2", cookie[0].Name+"="+cookie[0].Value)
}

func InitRedirect() {
	c := RestyInit()
	resp, _ := c.R().
		SetHeader("Cookie", tempDB.Get("cookie2").(string)).
		Get("https://www.luogu.com.cn/auth/login")
	cookie := resp.Cookies()
	tempDB.Set("cookie1", cookie[0].Name+"="+cookie[0].Value)
}

func RestyInit() *resty.Client {
	c := resty.New()
	c.SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0")
	return c
}
