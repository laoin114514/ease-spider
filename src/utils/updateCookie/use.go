package updatecookie

import (
	"fmt"
	"spider/src/utils/updateCookie/component"
	"spider/src/utils/updateCookie/request"
)

func Use(isInServer bool) {
	request.Init()
	request.InitRedirect()
	request.ChooseMth("laoi请问")
	request.GetCaptcha()
	request.RedirCaptcha()
	captcha, err := component.Identify(isInServer)
	if err != nil {
		return
	}
	fmt.Println(captcha)
	request.Login(captcha)
}
