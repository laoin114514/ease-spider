package main

import (
	"fmt"
	"log"
	"spider/config"
	"spider/config/db"
	"spider/src/handler"
	"spider/src/services"
	"spider/src/utils"
	"sync"
)

func main() {
	// 初始化配置
	err := config.Init()
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 验证配置
	validator := utils.NewConfigValidator()
	if err := validator.ValidateConfig(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 初始化数据库
	err = db.Init()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化JSON数据库和CF URL生成器
	utils.InitGlobalJSONDB("config.json")
	utils.InitGenerateCFurl()

	log.Println("系统初始化完成，开始执行定时任务...")
	// func() error {
	// 	log.Println("开始获取洛谷用户提交记录")
	// 	luoguRecordsService := services.NewLuoguRecordsService()
	// 	luoguUpdateCookieService := services.NewLuoguUpdateCookie()
	// 	err := luoguUpdateCookieService.UpdateLuoguCookie()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	luoguUpdateCookieService.SaveLog()
	// 	luoguUpdateCookieService.SaveErr()
	// 	luoguUpdateCookieService.Clear()
	// 	err = luoguRecordsService.GetLuoguUsersRecords()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = luoguRecordsService.ChangePrivateProblem()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	luoguRecordsService.SaveLog()
	// 	luoguRecordsService.SaveErr()
	// 	luoguRecordsService.Clear()
	// 	log.Println("洛谷用户提交记录获取完成")
	// 	return nil
	// }()
	func() error {
		log.Println("开始获取CF用户提交记录")
		cfRecordsService := services.NewCfService()
		err := cfRecordsService.CfUserStatus.GetCfRecords()
		if err != nil {
			return err
		}
		log := cfRecordsService.CfUserStatus.GetLog()
		fmt.Println(len(log))
		errLog := cfRecordsService.CfUserStatus.GetErr()
		fmt.Println(len(errLog))
		return nil
	}()
	return
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	handler.TimeTask()
}
