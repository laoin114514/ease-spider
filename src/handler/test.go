package handler

import (
	"spider/src/services"
)

// 调试脚本（测试用）
func Test() {
	func() error {
		updateCookie := services.NewLuoguUpdateCookie()
		err := updateCookie.Update()
		if err != nil {
			return err
		}
		luogu := services.NewLuogu()
		err = luogu.LuoguSubmissionDetail.GetAndStoreSourceCode()
		if err != nil {
			return err
		}
		luogu.LuoguSubmissionDetail.SaveLog()
		luogu.LuoguSubmissionDetail.SaveErr()
		luogu.LuoguSubmissionDetail.Clear()
		return nil
	}()
}
