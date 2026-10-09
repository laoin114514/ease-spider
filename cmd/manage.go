package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"spider/config"
	"spider/config/db"
	"spider/internal/crawlers"
	"spider/internal/utils"
	"syscall"

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
	// 无论走哪条分支，退出前都把插件日志文件句柄释放掉
	defer func() {
		if err := e.Close(); err != nil {
			ease.Logger.Errorf("关闭插件日志文件失败: %v", err)
		}
	}()

	args := os.Args
	if len(args) < 2 {
		fmt.Println("用法: ./spider run | devrun <crawlerKey> | list")
		return
	}

	cmd := args[1]
	switch cmd {
	case "run":
		runUntilSignal(e)
	case "devrun":
		if len(args) < 3 {
			fmt.Println("用法: ./spider devrun <crawlerKey>")
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

// runUntilSignal 启动引擎并阻塞到进程收到退出信号。
//
// 收到 SIGINT(Ctrl-C) / SIGTERM(docker stop) 后的退出流程：
//  1. 取消引擎 ctx：各插件跑完当前这轮任务后退出，正在进行的入库不被打断
//  2. 引擎最多等待 SetShutdownTimeout 指定的时长（默认 30s），超时直接返回
//  3. main 的 defer 调用 e.Close() 释放插件日志文件句柄
//
// 另外，第一次信号之后会恢复系统默认的信号处理：此时再发一次信号即可强制
// 终止进程，用于任务收尾太慢时立刻退出。
func runUntilSignal(e *ease.Engine) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		// 恢复默认行为：此后再收到信号会直接杀掉进程
		stop()
		ease.Logger.Println("收到退出信号，等待正在执行的任务收尾；再次发送信号可强制退出")
	}()

	e.RunWithContext(ctx)
	ease.Logger.Println("爬虫引擎已退出，进程即将结束")
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
