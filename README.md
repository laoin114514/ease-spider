# Spider 爬虫系统开发文档

## 项目简介

Spider 是一个基于 Go 语言开发的竞赛数据爬虫系统，主要用于爬取和监控 Codeforces 和洛谷平台的竞赛数据。系统采用定时任务机制，自动收集用户提交记录、比赛信息、题目数据等，并通过钉钉机器人进行消息推送。

## 核心功能

### 数据爬取
- **Codeforces**: 官方比赛、题目、用户提交记录、团队比赛数据
- **洛谷**: 用户提交记录、题解、源码、训练详情、团队数据
- **钉钉**: 消息推送和通知服务

### 定时任务
- 支持多种定时频率配置（秒/分钟/小时/天）
- 并发控制，避免API限制
- 自动重试和错误处理

### 🛠 工具集
- HTTP请求封装
- JSON数据库
- 并发控制器
- 日志管理
- 配置验证

## 项目结构

```
spider/
├── config/                 # 配置管理
│   ├── config.go          # 配置结构定义
│   └── db/                # 数据库相关
├── src/
│   ├── constants/         # 常量定义
│   ├── handler/           # 任务处理器
│   ├── models/           # 数据模型
│   ├── repository/       # 数据访问层
│   ├── services/         # 业务逻辑层
│   └── utils/            # 工具函数
├── help/                 # 文档
├── logs/                 # 日志文件
├── ocr/                  # OCR功能（可选）
└── main.go              # 程序入口
```

## 安装和运行

### 环境要求

- Go 1.19 或更高版本
- MySQL 5.7 或更高版本
- Python 3.8+ (用于OCR功能,不用时可忽略)

### 安装步骤

1. **克隆项目**
```bash
git clone https://gitlab.unde.site/QingLuan/spider.git
cd spider
git checkout feature/go
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置数据库**
```sql
-- 创建数据库
CREATE DATABASE gxuicpc CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'gxuicpc'@'%' IDENTIFIED BY 'gxuicpc';
GRANT ALL PRIVILEGES ON gxuicpc.* TO 'gxuicpc'@'%';
FLUSH PRIVILEGES;
```

4. **配置系统**
编辑 `config.yml` 文件，配置数据库连接、API密钥、定时任务频率等参数。

5. **运行程序**
```bash
# 开发环境运行
go run main.go

# 编译后运行
go build -o spider
./spider
```

6. **构建成linux可执行文件**
```bash
.\buildTolinux.bat
```

## 配置说明

### 主要配置项

```yaml
database:          # 数据库配置
  host: 210.36.22.245
  port: 3002
  user: gxuicpc
  password: gxuicpc
  dbName: gxuicpc

luogu:            # 洛谷配置
  luoguRecordsConcurrency: 4      # 并发数
  luoguTeamID: 116191             # 团队ID
  username: laoyin               # 用户名
  password: 683305SAO           # 密码

cf:               # Codeforces配置
  cfRecordsConcurrency: 4        # 并发数
  managerAccount: 233zhang      # 管理员账号

timerFrequency:   # 定时任务频率
  cf_records: 2m                # CF记录爬取频率
  luogu_records: 2m             # 洛谷记录爬取频率
  luogu_update_cookie: 1h       # Cookie更新频率
```

## 开发指南

### 数据库管理

1. **新建表格**: 在 `/config/db/tableModels.go` 中更新表结构
2. **数据查询**: 在各服务对应的 repository 中编写查询逻辑
3. **数据插入**: 使用预定义的表结构，避免字段错误

### 服务开发

1. **模型定义**: 在 `src/models/` 中定义数据结构
2. **业务逻辑**: 在 `src/services/` 中实现核心功能
3. **数据访问**: 在 `src/repository/` 中处理数据库操作
4. **任务处理**: 在 `src/handler/` 中处理定时任务

### 工具使用

- **HTTP请求**: 使用 `utils.NewRequest[T]()` 进行API调用
- **并发控制**: 使用 `utils.NewConCurrenter[T]()` 控制并发
- **日志管理**: 使用 `utils.NewLogContainer()` 收集日志
- **配置验证**: 使用 `utils.NewConfigValidator()` 验证配置

### 定时任务

在 `src/handler/timeTask.go` 中添加新的定时任务：

```go
timer.RunWithTimer(
    utils.NewDateFormat().AnalysisTimerFrequency("30m"),
    "任务描述",
    func() error {
        // 任务逻辑
        return nil
    },
)
```

## API接口

### Codeforces API
- 官方比赛列表
- 题目信息
- 用户提交记录
- 团队比赛数据

### 洛谷 API
- 用户提交记录
- 题解内容
- 源码获取
- 训练详情
- 团队信息

## 日志管理

系统自动生成日志文件到 `logs/` 目录：
- `.log` 文件：正常日志
- `.err.log` 文件：错误日志

## 故障排除

### 常见问题

1. **数据库连接失败**: 检查 `config.yml` 中的数据库配置
2. **API请求失败**: 检查网络连接和API密钥
3. **并发过高**: 调整配置文件中的并发数设置
4. **Cookie过期**: 系统会自动更新，检查用户名密码是否正确

### 调试模式

在 `config.yml` 中设置：
```yaml
debug:
  all: true
```

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request
