# ease-crawler 使用文档

## 1. 框架定位

`ease-crawler` 是一个按“分组 + 插件”组织的轻量调度框架，支持：

- 分组链式注册（`Group(...).Group(...).Register(...)`）
- 插件定时执行（`Interval`）
- 启动即执行（`StartImmediately`）
- 插件日志按分组目录落盘（`logs/<group>/<crawler>.log`）
- 插件内注入专属 logger（通过 `ContextLoggerKey`）

---

## 2. 目录与日志规则

默认日志根目录是：

- `logs`

插件日志路径规则：

- `logs/<group-path>/<crawler-name>.log`

示例：

- 分组：`测试1/测试1`
- 插件：`task1`
- 日志文件：`logs/测试1/测试1/task1.log`

---

## 3. 插件接口

实现 `Crawler`：

```go
type Crawler interface {
    Run(c *Context) error
    Name() string
    Meta() Meta
}
```

`Meta`：

- `Interval`：执行间隔（<=0 自动回退为 1 秒）
- `StartImmediately`：启动后是否先跑一次
- `Logger`：可选，自定义 logger（不设置则由框架注入分组 logger）

---

## 4. 插件内日志写法

框架在每次执行前会注入 logger 到 `Context`：

- key: `ContextLoggerKey`

插件中建议这样写：

```go
func (t *Task) Run(c *easecrawler.Context) error {
    if logger, ok := easecrawler.GetAs[*log.Logger](c, easecrawler.ContextLoggerKey); ok {
        logger.Println("开始执行")
    }
    return nil
}
```

---

## 5. 快速开始

```go
package main

import (
    "log"
    "time"
    easecrawler "spider/pkg/ease-crawler"
)

type Task1 struct{}

func (t *Task1) Name() string { return "task1" }
func (t *Task1) Meta() easecrawler.Meta {
    return easecrawler.Meta{
        Interval:         2 * time.Second,
        StartImmediately: true,
    }
}
func (t *Task1) Run(c *easecrawler.Context) error {
    if logger, ok := easecrawler.GetAs[*log.Logger](c, easecrawler.ContextLoggerKey); ok {
        logger.Println("task1 running")
    }
    return nil
}

func main() {
    engine := easecrawler.New()
    defer engine.Close()

    g := engine.Group("测试1").Group("测试1")
    g.Register(&Task1{})

    // 阻塞运行（内部按 Interval 调度）
    engine.Run()
}
```

---

## 6. 可选配置

### 设置日志根目录

```go
engine.SetLogRootDir("my-logs")
```

### 按 Context 控制退出

```go
ctx, cancel := context.WithCancel(context.Background())
// ... 某处 cancel()
engine.RunWithContext(ctx)
```

---

## 7. 运行建议

1. 生产环境建议总是 `defer engine.Close()`（释放文件句柄）。
2. 插件内部统一用注入 logger，不建议直接 `fmt.Println`。
3. 分组名称尽量稳定，便于日志目录长期维护。
4. 同一分组内插件名要唯一（重复注册会被跳过并记录日志）。

---

## 8. 当前行为说明（已验证）

- `go test ./...` 可通过（当前仓库无测试文件，但编译链路正常）。
- 框架已可稳定运行并按分组目录输出日志。
