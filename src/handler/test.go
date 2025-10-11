package handler

import (
	"fmt"
	"spider/src/services"
)

func Test() {
	luoguUpdateCookieService := services.NewLuoguUpdateCookie()
	luoguUpdateCookieService.UpdateLuoguCookie()
	fmt.Println(luoguUpdateCookieService.LogService.GetLog())
	fmt.Println(luoguUpdateCookieService.LogService.GetErr())
	luoguUpdateCookieService.LogService.Clear()
	luoguSubmissionDetailService := services.NewLuoguSubmissionDetail()
	err := luoguSubmissionDetailService.GetRecordSourceCode()
	if err != nil {
		fmt.Println(err)
	}
	luoguSubmissionDetailService.LogService.Clear()
}
