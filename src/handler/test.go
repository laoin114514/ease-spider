package handler

import (
	"log"
	"spider/src/services"
)

func Test() {
	func() error {
		luoguRecordsService := services.NewLuoguRecordsService()
		log.Println("开始获取洛谷用户提交记录")
		err := luoguRecordsService.GetLuoguUsersRecords()
		if err != nil {
			return err
		}
		err = luoguRecordsService.ChangePrivateProblem()
		if err != nil {
			return err
		}
		luoguRecordsService.SaveLog()
		luoguRecordsService.SaveErr()
		luoguRecordsService.Clear()
		log.Println("洛谷用户提交记录获取完成")
		return nil
	}()
}
