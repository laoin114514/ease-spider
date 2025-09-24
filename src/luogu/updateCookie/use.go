package updatecookie

import (
	"fmt"
	"spider/src/luogu/updateCookie/component"
	request2 "spider/src/luogu/updateCookie/request"
)

func Use() {
	request2.Init()
	request2.InitRedirect()
	request2.ChooseMth("luo2908451607")
	request2.GetCaptcha()
	request2.RedirCaptcha()
	captcha, err := component.Identify(true)
	if err != nil {
		return
	}
	fmt.Println(captcha)
	request2.Login(captcha)
}
