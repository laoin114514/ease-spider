package main

import (
	"fmt"
	"spider/config/db"
	"spider/src/services"
	"spider/src/utils"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db.Init()
	utils.InitGlobalJSONDB("tempDB.json")
	timer := utils.NewTimer()
	cfService := services.NewCfService()
	str := timer.CountDurationStr(func() {
		cfService.CfUserStatus.GetCfRecords(2)
	})
	fmt.Println("日志", cfService.CfUserStatus.GetLog())
	fmt.Println("错误", cfService.CfUserStatus.GetErr())
	fmt.Println("时间", str)
}
