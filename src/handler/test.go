package handler

import (
	"log"
	"spider/src/services"
)

// 调试脚本（测试用）
func Test() {
	cfService := services.NewCfService()
	func() error {
		log.Println("开始获取cf提交记录")
		err := cfService.CfUserStatus.GetCfRecords()
		if err != nil {
			return err
		}
		cfService.CfUserStatus.SaveLog()
		cfService.CfUserStatus.SaveErr()
		cfService.CfUserStatus.Clear()
		log.Println("cf提交记录获取完成")
		return nil
	}()
}
