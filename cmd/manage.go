package main

import (
	"fmt"
	"os"
	"spider/config"
	"spider/config/db"
	"spider/internal/crawlers"
	"spider/internal/utils"

	ease "spider/pkg/crawler"
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

	ease.Logger.Println("系统初始化完成")
}
func main() {
	e := crawlers.Register()
	args := os.Args
	if len(args) < 2 {
		fmt.Println("用法: go run ./cmd/manage/manage.go run | devrun <crawlerKey> | list")
		return
	}

	cmd := args[1]
	switch cmd {
	case "run":
		e.Run()
	case "devrun":
		if len(args) < 3 {
			fmt.Println("用法: go run ./cmd/manage/manage.go devrun <crawlerKey>")
			return
		}
		key := args[2]
		e.DevRun(key)
	case "list":
		crawlers := e.GetCrawlers()
		if len(crawlers) == 0 {
			fmt.Println("没有已注册的插件")
			return
		}
		fmt.Println("已注册插件 key 列表:")
		for key := range crawlers {
			fmt.Println("-", key)
		}
	default:
		fmt.Println("未知命令, 可用: run | devrun | list")
	}
}

//
//                            _ooOoo_
//                           o8888888o
//                           88" . "88
//                           (| -_- |)
//                           O\  =  /O
//                        ____/`---'\____
//                      .'  \\|     |//  `.
//                     /  \\|||  :  |||//  \
//                    /  _||||| -:- |||||-  \
//                    |   | \\\  -  /// |   |
//                    | \_|  ''\---/''  |   |
//                    \  .-\__  `-`  ___/-. /
//                  ___`. .'  /--.--\  `. . __
//               ."" '<  `.___\_<|>_/___.'  >'"".
//              | | :  `- \`.;`\ _ /`;.`/ - ` : | |
//              \  \ `-.   \_ __\ /__ _/   .-` /  /
//         ======`-.____`-.___\_____/___.-`____.-'======
//                            `=---='
//        ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//                       佛祖保佑永无bug
