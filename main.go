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
	cfService.CfUserStatus.GetCfRecords(5)
	fmt.Println(cfService.CfUserStatus.GetLog())
	fmt.Println(cfService.CfUserStatus.GetErr())
	fmt.Println("--------------------------------")
	fmt.Println(timer.CountDurationStr(func() {
		cfService.CfOfficialProblems.GetCfOfficialProblems()
	}))
	fmt.Println(cfService.CfOfficialProblems.GetLog())
	fmt.Println(cfService.CfOfficialProblems.GetErr())
}
