package request

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func Login(captcha string) {
	c := resty.New()
	cookie := tempDB.Get("cookie1").(string) + "; " + tempDB.Get("cookie2").(string)
	resp, err := c.R().
		SetBody(map[string]any{
			"username": "laoin",
			"password": "683305SAo",
			"captcha":  captcha,
		}).
		SetHeader("Cookie", cookie).
		Post("https://www.luogu.com.cn/do-auth/password")
	fmt.Println(resp.Status())
	if resp.StatusCode() != 200 {
		fmt.Println("验证码错误")
		return
	}
	if err != nil {
		fmt.Println(err)
	}
	arr := resp.Cookies()
	uid := arr[0].Name + "=" + arr[0].Value
	Cookie := cookie + ";" + uid
	tempDB.Set("Cookie", Cookie)
	fmt.Println("登录成功")
}
