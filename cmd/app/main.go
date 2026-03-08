package main

import (
	"spider/config"
	"spider/config/db"
	"spider/internal/crawlers"
	"spider/internal/utils"
	ease "spider/pkg/ease-crawler"
)

func init() {
	// 初始化配置
	err := config.Init()
	if err != nil {
		ease.Logger.Fatalf("配置初始化失败: %v", err)
	}

	// 验证配置
	validator := utils.NewConfigValidator()
	if err := validator.ValidateConfig(); err != nil {
		ease.Logger.Fatalf("配置验证失败: %v", err)
	}

	// 初始化数据库
	err = db.Init()
	if err != nil {
		ease.Logger.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化JSON数据库和CF URL生成器
	utils.InitGlobalJSONDB("data.json")

	ease.Logger.Println("系统初始化完成，开始执行定时任务...")
}
func main() {
	crawlers.Run()
}
