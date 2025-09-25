package updatecookie

import (
	"fmt"
	"spider/src/luogu/updateCookie/component"
	"spider/src/luogu/updateCookie/request"
)

func Use(isInServer bool) {
	request.Init()
	request.InitRedirect()
	request.ChooseMth("luo2908451607")
	request.GetCaptcha()
	request.RedirCaptcha()
	captcha, err := component.Identify(isInServer)
	if err != nil {
		return
	}
	fmt.Println(captcha)
	request.Login(captcha)
}
