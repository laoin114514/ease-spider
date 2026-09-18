# Spider 插件化爬虫项目说明

本项目基于 `pkg/ease-crawler` 构建，采用“**分组 + 插件**”的方式组织抓取任务。

你可以把它理解成：

- `engine` 负责调度（间隔执行、启动即跑、日志注入、panic 兜底）
- `internal/crawlers/*` 负责具体业务插件
- `internal/crawlers/registry.go` 负责统一注册所有插件

---

## 1. 插件在哪里注册

统一注册入口：

- `internal/crawlers/registry.go`

当前注册方式（示例）：

- `luogu` 分组注册洛谷相关插件
- `cf` 分组注册 Codeforces 相关插件
- `dingding` 分组注册钉钉相关插件

即在 `Run()` 里通过：

- `e := easecrawler.New()`
- `group := e.Group("xxx")`
- `group.Register(plugin)`
- `e.Run()`

完成启动。

---

## 2. 插件目录规范（推荐）

每个插件目录建议保持三层结构：

- `base.go`：插件定义与框架接口实现
- `model.go`：该插件专用结构体（请求/响应/中间结构）
- `service.go`：业务逻辑实现（抓取、解析、入库）

目录示例（洛谷分组目前只剩提交记录一个插件，因此没有 `model.go`）：

- `internal/crawlers/luogu/get_user_records/base.go`
- `internal/crawlers/luogu/get_user_records/service.go`

---

## 3. 插件体应该怎么写

最小插件需要实现 `Crawler` 接口：

```go
package demo

import (
    easecrawler "spider/pkg/ease-crawler"
    "time"
)

type DemoCrawler struct {
    log *easecrawler.EaseLogger
}

func NewDemoCrawler() *DemoCrawler {
    return &DemoCrawler{}
}

func (d *DemoCrawler) Name() string {
    return "demo_task"
}

func (d *DemoCrawler) Meta() easecrawler.Meta {
    return easecrawler.Meta{
        Interval:         1 * time.Minute,
        StartImmediately: true,
    }
}

func (d *DemoCrawler) Run(c *easecrawler.Context) error {
    d.log = easecrawler.GetCrawlerLogger(c)
    return d.DoWork()
}
```

说明：

- `Name()`：插件唯一标识（同分组下不可重复）
- `Meta()`：调度配置
  - `Interval`：执行周期
  - `StartImmediately`：启动时是否先执行一次
  - `Logger`：可选，自定义日志器
- `Run()`：框架调用入口，建议只做上下文准备，具体逻辑放 `service.go`

---

## 4. service.go 该怎么组织

`service.go` 建议只放业务步骤，示例模式：

1. 初始化上下文（请求客户端、cookie 等）
2. 拉取远端数据
3. 解析/转换
4. 入库
5. 记录日志并返回错误

常见结构：

- `Update()` / `GetAndStore()`：对外主流程
- `initXXX()`：初始化
- `fetchXXX()`：抓取
- `processXXX()`：处理
- `saveXXX()`：持久化

---

## 5. 日志规则

框架会为每个插件创建日志文件，默认路径：

- `logs/<group>/<plugin>.log`

例如：

- `logs/luogu/get_user_records.log`
- `logs/cf/get_user_records.log`

在插件内部通过：

- `easecrawler.GetCrawlerLogger(c)`

拿到注入日志器后，统一 `Printf/Errorf/Warnf` 输出。

---

## 6. 新增插件完整步骤

1. 新建目录：`internal/crawlers/<group>/<plugin>/`
2. 新增 `base.go`，实现 `Name/Meta/Run`
3. 新增 `service.go`，承载业务逻辑
4. （可选）新增 `model.go`，放插件专用结构
5. 在 `internal/crawlers/registry.go` 注册插件
6. 启动验证日志文件与任务执行

---

## 7. 运行方式

```bash
go run ./cmd run                          # 启动全部插件
go run ./cmd list                         # 查看已注册插件
go run ./cmd devrun luogu/get_user_records  # 单跑某个插件
```

洛谷提交记录经 Luogu2Api 服务（`pkg/luogu2api` SDK）获取，运行前需在
`config/config.dev.yml` / `config/config.prod.yml` 的 `luogu2api` 段填好 `baseUrl`
与 `adminToken`（后者与服务端 `ADMIN_TOKEN` 一致），两项为空时启动阶段就会失败。

---

## 8. 常见问题排查

### 1) 插件没有执行

检查：

- 是否在 `registry.go` 中注册
- `Meta.Interval` 是否设置合理
- 是否启动了主程序

### 2) 日志文件没生成

检查：

- 是否调用了 `GetCrawlerLogger(c)` 并有日志输出
- 进程对 `logs/` 目录是否有写权限

### 3) 运行时 panic 导致任务中断

框架调度层已做 panic recover，单插件 panic 不应导致主进程退出。
若仍退出，请检查 panic 是否发生在插件外部初始化链路。

### 4) 洛谷提交记录任务报错

先检查 `luogu2api.baseUrl` / `luogu2api.adminToken` 是否与服务端一致；任务运行中的
常见错误（业务码见 `pkg/luogu2api/errors.go`）：

- `401`：`adminToken` 与服务端 `ADMIN_TOKEN` 不一致
- `1001`（HTTP 503）：号池没有可用账号，稍后重试即可
- `1002`（HTTP 502）：号池内账号登录态全部失效，需要人工重新登录

插件日志在 `logs/luogu/get_user_records.log`。
