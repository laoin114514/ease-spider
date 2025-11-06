package main

import (
	"log"
	"spider/config"
	"spider/config/db"
	"spider/src/handler"
	"spider/src/utils"
	"sync"
)

func main() {
	// 初始化配置文件
	err := config.MakeYml()
	if err != nil {
		log.Fatalf("配置文件初始化失败: %v", err)
	}
	// 初始化配置
	err = config.Init()
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
	utils.InitGlobalJSONDB("data.json")
	utils.InitGenerateCFurl()

	log.Println("系统初始化完成，开始执行定时任务...")

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	handler.TimeTask()
}
