package handler

import (
	"spider/src/services"
)

// 调试脚本（测试用）
func Test() {
	// func() error {
	// 	updateCookie := services.NewLuoguUpdateCookie()
	// 	err := updateCookie.Update()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	luogu := services.NewLuogu()
	// 	err = luogu.LuoguRecords.GetAndStore()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	luogu.LuoguRecords.SaveLog()
	// 	luogu.LuoguRecords.SaveErr()
	// 	luogu.LuoguRecords.Clear()
	// 	return nil
	// }()
	func() error {
		cf := services.NewCfService()
		err := cf.CfUserStatus.GetCfRecords()
		if err != nil {
			return err
		}
		cf.CfUserStatus.SaveLog()
		cf.CfUserStatus.SaveErr()
		return nil
	}()
}
