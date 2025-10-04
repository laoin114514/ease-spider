package request

func ChooseMth(account string) {
	c := RestyInit()
	cookie := tempDB.Get("cookie1").(string) + ";" + tempDB.Get("cookie2").(string)
	c.R().SetHeader("Cookie", cookie).Get("https://www.luogu.com.cn/auth/login-methods?login=" + account)
}
