# 用于spider插件化爬虫项目的简易爬虫框架

## 框架结构
- `conCurrenter.go` 一个并发器, 用于控制并发, 启动多个爬虫任务
- `engine.go` 是框架核心调度器, 负责管理分组结构和已经注册的插件, 维护日志的输出和运行时的资源
- `crawler.go` 定义了爬虫插件的最小实现接口, 外部循环调度由`engine.go`负责
- `logger.go` 用于该框架的日志输出

---

## 框架使用方法

- **定义插件:** 创建实现 Crawler接口的结构体。
- **创建引擎:** 通过 easecrawler.New()初始化。
- **注册插件:** 在适当的分组下调用 Register方法。
- **启动引擎:** 调用 RunWithContext并传入一个可取消的上下文，以实现阻塞运行和优雅停止。
- **释放资源:** 程序退出前调用 engine.Close()。

## 代码示例

```go
package main

import (
	easecrawler "gitlab.unde.site/QingLuan/spider/pkg/easecrawler"
	"time"
)

// MyCrawler 示例：一个简单的采集插件
type MyCrawler struct {
	Name string
}

func (c *MyCrawler) Run(ctx *easecrawler.Context) error {
	// 1. 从上下文中获取日志器
	logger := easecrawler.GetCrawlerLogger(ctx)
	
	// 2. 执行你的采集任务
	logger.Printf("插件 [%s] 开始执行...", c.Name)
	// TODO: 你的业务逻辑，如 HTTP 请求、数据解析、存储等
	// 例如：resp, err := http.Get("https://example.com")
	
	logger.Printf("插件 [%s] 执行完成", c.Name)
	return nil // 返回 nil 表示成功，返回 error 表示失败
}

func (c *MyCrawler) Name() string {
	return c.Name
}

func (c *MyCrawler) Meta() easecrawler.Meta {
	return easecrawler.Meta{
		Interval:          10 * time.Second, // 每10秒执行一次
		StartImmediately:  true,             // 引擎启动后立即执行一次
		Logger:            nil,               // 使用框架分配的日志器
	}
}
```
