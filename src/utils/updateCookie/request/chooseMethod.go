package request

func ChooseMth(account string) {
	c := RestyInit()
	cookie := JsonDB.Get("cookie1").(string) + ";" + JsonDB.Get("cookie2").(string)
	c.R().SetHeader("Cookie", cookie).Get("https://www.luogu.com.cn/auth/login-methods?login=" + account)
}
