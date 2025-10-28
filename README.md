# 青鸾系统爬虫项目
## 这是一个项目介绍  
开发文档移步根目录下help文件夹
## 项目简介

青鸾系统是一个基于Go语言开发的多平台数据爬虫系统，主要用于爬取和监控以下平台的数据：

- **Codeforces** - 获取用户提交记录、官方比赛、团队比赛、题目信息
- **洛谷 (Luogu)** - 获取用户提交记录、题目信息
- **钉钉 (DingTalk)** - 获取队员打卡数据

系统采用定时任务机制，支持并发爬取，具备反作弊检测功能，适用于企业级数据采集和监控场景。

## 项目架构

```
spider/
├── config/                 # 配置文件
│   ├── config.yml         # 主配置文件
│   ├── config.json        # JSON配置
│   └── help/              # 帮助文档
├── src/
│   ├── constants/         # 常量定义
│   ├── handler/           # 处理器
│   ├── models/            # 数据模型
│   ├── repository/        # 数据访问层
│   ├── services/          # 业务逻辑层
│   └── utils/             # 工具类
├── logs/                  # 日志文件
├── main.go               # 主程序入口
└── go.mod                # Go模块文件
```

## 配置说明

### 主配置文件 (config.yml)

```yaml
database:
  host: 210.36.22.245
  port: 3002
  user: gxuicpc
  password: gxuicpc
  dbName: gxuicpc

dingding:
  appKey: dingsgji79glruqayy49
  appSecret: 0jzUXpbjMFs4jK3NNE1Hn3PTBJ8TZfVKurIR4_XI_2jDVmGD3Dkle7xQ4xtHH8e6
luogu:
  isInServer: false
  luoguRecordsConcurrency: 4
  luoguSolutionConcurrency: 5
  luoguSubmissionDetailConcurrency: 5
  luoguTeamID: 116191
  username: laoyin
  password: 683305SAO
  userAgent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0
cf:
  cfRecordsConcurrency: 4
  cfTeamContestProblemsConcurrency: 2
  managerAccount: 233zhang
  groupCode: mczZWy5WU7
#==================定时任务配置==================
#每个任务的执行时间，单位填在数字后面，如30s，30m，30h，30d，默认分钟,更改后需要重启程序
timerFrequency:
  cf_official_contests: 1h
  cf_official_problems: 1h
  cf_team_contests: 1h
  cf_team_contests_problems: 1h
  dingding: 1h
  luogu_update_cookie: 1h
  luogu_submission_detail: 1h
  luogu_solution: 1h
  cf_records: 2m
  luogu_records: 2m
debug:
  all: true
```

### JSON配置文件 (config.json)

```json
{
    "Cookie": "__client_id=ee3223357f3634d768e648a26c59c972e289d422; C3VK=696cf2;_uid=1851093",
    "dingUserId": {
        "01176421420026309256": "杜永坤",
        "03030233253124066409": "张健明",
        "0305145237061099840053": "韦书临",
        "03051506462725991660": "曾佳宁",
        "225968356337662726": "陈君屹",
        "481366495340096859": "黄翠婷",
        "49203919162724192": "陶康",
        "514450420223269518": "孙怿翔",
        "54210856181176852": "郑毅"
    }
}
```

## API数据结构

### 1. Codeforces API

#### 1.1 用户提交记录

**请求参数:**
```json
{
    "handle": "tourist",
    "from": 1,
    "count": 100
}
```

**响应数据结构:**
```json
{
    "status": "OK",
    "result": [
        {
            "id": 123456789,
            "contestId": 1000,
            "creationTimeSeconds": 1640995200,
            "relativeTimeSeconds": 1800,
            "problem": {
                "contestId": 1000,
                "index": "A",
                "name": "Watermelon",
                "type": "PROGRAMMING",
                "tags": ["brute force", "math"],
                "points": 500.0,
                "rating": 800
            },
            "author": {
                "contestId": 1000,
                "participantId": 12345,
                "members": [
                    {
                        "handle": "tourist"
                    }
                ],
                "participantType": "CONTESTANT",
                "teamId": 0,
                "teamName": "",
                "ghost": false,
                "startTimeSeconds": 1640995200
            },
            "programmingLanguage": "GNU C++17",
            "verdict": "OK",
            "testset": "TESTS",
            "passedTestCount": 10,
            "timeConsumedMillis": 1000,
            "memoryConsumedBytes": 1024000
        }
    ]
}
```

#### 1.2 官方比赛列表

**请求参数:**
```json
{
    "gym": false
}
```

**响应数据结构:**
```json
{
    "status": "OK",
    "result": [
        {
            "id": 1000,
            "name": "Codeforces Round #1000",
            "type": "CF",
            "startTime": 1640995200,
            "preparedBy": "tourist",
            "phase": "FINISHED",
            "frozen": false,
            "durationSeconds": 7200,
            "startTimeSeconds": 1640995200,
            "relativeTimeSeconds": 306246828,
            "country": "Russia",
            "city": "Moscow",
            "kind": "Official",
            "season": "2021-2022",
            "difficulty": 5
        }
    ]
}
```

#### 1.3 题目信息

**请求参数:**
```json
{
    "tags": "implementation;math",
    "problemsetName": "acmsguru"
}
```

**响应数据结构:**
```json
{
    "status": "OK",
    "result": {
        "problems": [
            {
                "contestId": 1000,
                "index": "A",
                "name": "Watermelon",
                "type": "PROGRAMMING",
                "tags": ["brute force", "math"],
                "points": 500.0,
                "rating": 800
            }
        ]
    }
}
```

#### 1.4 团队比赛

**请求参数:**
```json
{
    "groupCode": "mczZWy5WU7"
}
```

**响应数据结构:**
```json
{
    "status": "OK",
    "result": [
        {
            "id": 100003,
            "name": "2008-2009 Всероссийская командная олимпиада школьников по программированию",
            "type": "ICPC",
            "phase": "FINISHED",
            "frozen": false,
            "durationSeconds": 18000,
            "startTimeSeconds": 1453514400,
            "relativeTimeSeconds": 306246828,
            "preparedBy": "Edvard",
            "difficulty": 3,
            "kind": "Official School Contest",
            "country": "Russia",
            "city": "Saint Petersburg",
            "season": "2008-2009"
        }
    ]
}
```

### 2. 洛谷 API

#### 2.1 用户提交记录

**请求参数:**
```json
{
    "user": "1851093",
    "page": "1",
    "_contentOnly": "1"
}
```

**响应数据结构:**
```json
{
    "code": 200,
    "currentTemplate": "record_list",
    "currentData": {
        "records": {
            "result": [
                {
                    "time": 1000,
                    "memory": 1024000,
                    "problem": {
                        "pid": "P1000",
                        "title": "A+B Problem",
                        "difficulty": 1,
                        "fullScore": 100,
                        "type": "P"
                    },
                    "contest": {
                        "id": 1000,
                        "name": "洛谷月赛",
                        "startTime": 1640995200,
                        "endTime": 1640998800,
                        "type": 1,
                        "rated": true,
                        "description": "洛谷月赛"
                    },
                    "sourceCodeLength": 50,
                    "submitTime": 1640995200,
                    "language": 1,
                    "user": {
                        "uid": 1851093,
                        "name": "tourist",
                        "avatar": "https://cdn.luogu.com.cn/upload/usericon/1851093.png",
                        "slogan": "Hello World",
                        "badge": null,
                        "isAdmin": false,
                        "isBanned": false,
                        "color": "Red",
                        "ccfLevel": 9,
                        "xcpcLevel": 0,
                        "background": "",
                        "verified": true,
                        "followingCount": 0,
                        "followerCount": 0,
                        "ranking": 1,
                        "eloValue": 3000,
                        "blogAddress": null,
                        "unreadMessageCount": 0,
                        "unreadNoticeCount": 0
                    },
                    "id": 123456789,
                    "status": 12,
                    "enableO2": false,
                    "score": 100
                }
            ],
            "perPage": 20,
            "count": 1000
        }
    },
    "currentTitle": "提交记录",
    "currentTheme": {
        "id": 1,
        "header": {
            "imagePath": null,
            "color": [[255, 255, 255]],
            "blur": 0,
            "brightness": 0,
            "degree": 0,
            "repeat": 0,
            "position": [0, 0],
            "size": [100, 100],
            "type": 0,
            "__CLASS_NAME": "LuoguTheme"
        },
        "sideNav": {
            "logoBackgroundColor": [255, 255, 255],
            "color": [255, 255, 255],
            "invertColor": false,
            "__CLASS_NAME": "LuoguTheme"
        },
        "footer": {
            "imagePath": null,
            "color": [[255, 255, 255]],
            "blur": 0,
            "brightness": 0,
            "degree": 0,
            "repeat": 0,
            "position": [0, 0],
            "size": [100, 100],
            "type": 0,
            "__CLASS_NAME": "LuoguTheme"
        }
    },
    "currentTime": 1640995200,
    "currentUser": {
        "uid": 1851093,
        "name": "tourist",
        "avatar": "https://cdn.luogu.com.cn/upload/usericon/1851093.png",
        "slogan": "Hello World",
        "badge": null,
        "isAdmin": false,
        "isBanned": false,
        "color": "Red",
        "ccfLevel": 9,
        "xcpcLevel": 0,
        "background": "",
        "verified": true,
        "followingCount": 0,
        "followerCount": 0,
        "ranking": 1,
        "eloValue": 3000,
        "blogAddress": null,
        "unreadMessageCount": 0,
        "unreadNoticeCount": 0
    }
}
```

### 3. 钉钉 API

#### 3.1 获取访问令牌

**请求参数:**
```json
{
    "appkey": "ding_accessToken",
    "appsecret": "ding_appSecret"
}
```

**响应数据结构:**
```json
{
    "accessToken": "access_token_string",
    "expiresIn": 7200
}
```

#### 3.2 打卡记录

**请求参数:**
```json
{
    "workDateFrom": 1640995200000,
    "workDateTo": 1641081600000,
    "userIdList": ["01176421420026309256"],
    "offset": 0,
    "limit": 100
}
```

**响应数据结构:**
```json
{
    "errcode": 0,
    "errmsg": "ok",
    "recordresult": [
        {
            "baseCheckTime": 1640995200000,
            "baseMacAddr": "00:11:22:33:44:55",
            "bizId": "biz123456",
            "checkType": "OnDuty",
            "corpId": "corp123456",
            "deviceSN": "device123456",
            "gmtCreate": 1640995200000,
            "gmtModified": 1640995200000,
            "groupId": 123456,
            "id": 123456789,
            "isLegal": "Y",
            "locationMethod": "ATM",
            "locationResult": "Normal",
            "planCheckTime": 1640995200000,
            "planId": 123456,
            "sourceType": "ATM",
            "timeResult": "Normal",
            "userAddress": "北京市朝阳区",
            "userCheckTime": 1640995200000,
            "userId": "01176421420026309256",
            "workDate": 1640995200000
        }
    ]
}
```

## 安装和运行

### 环境要求

- Go 1.19 或更高版本
- MySQL 5.7 或更高版本
- Python 3.8+ (用于OCR功能)

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd spider
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
   - 复制 `config.yml.example` 为 `config.yml`
   - 复制 `config.json.example` 为 `config.json`
   - 根据实际情况修改配置参数

5. **运行程序**
```bash
# 开发环境运行
go run main.go

# 编译后运行
go build -o spider
./spider
```

### 交叉编译

```bash
# 编译为Linux版本
GOOS=linux GOARCH=amd64 go build -o spider-linux main.go

# 编译为Windows版本
GOOS=windows GOARCH=amd64 go build -o spider.exe main.go
```

## 使用指南

### 1. 配置用户列表

在 `config.json` 中配置需要监控的用户：

```json
{
    "Cookie": "你的洛谷Cookie",
    "dingUserId": {
        "用户ID1": "用户姓名1",
        "用户ID2": "用户姓名2"
    }
}
```

### 2. 启动定时任务

系统启动后会自动开始定时任务，包括：

- **Codeforces数据爬取** (每2-30分钟)
- **洛谷数据爬取** (每2分钟)
- **钉钉打卡数据** (每30分钟)

### 3. 查看日志

```bash
# 查看实时日志
tail -f logs/spider.log

# 查看错误日志
tail -f logs/spider.err.log

# 查看特定平台日志
tail -f logs/cfUserStatus.log
tail -f logs/luogu.log
```

### 4. 监控系统状态

系统提供以下监控指标：

- 爬取成功率
- 数据更新频率
- 错误率统计
- 并发任务状态

## 数据模型说明

### Codeforces 数据模型

#### CfSubmission (提交记录)
```go
type CfSubmission struct {
    Id                  int64     `json:"id"`                  // 提交ID
    ContestId           int64     `json:"contestId"`           // 比赛ID
    CreationTimeSeconds int64     `json:"creationTimeSeconds"` // 创建时间
    Problem             CfProblem `json:"problem"`             // 题目信息
    Author              CfAuthor  `json:"author"`              // 作者信息
    ProgrammingLanguage string    `json:"programmingLanguage"` // 编程语言
    Verdict             string    `json:"verdict"`             // 判题结果
    TimeConsumedMillis  int64     `json:"timeConsumedMillis"`  // 时间消耗
    MemoryConsumedBytes int64     `json:"memoryConsumedBytes"` // 内存消耗
}
```

#### CfContest (比赛信息)
```go
type CfContest struct {
    Id                  int64  `json:"id"`                  // 比赛ID
    Name                string `json:"name"`                // 比赛名称
    Type                string `json:"type"`                // 比赛类型
    Phase               string `json:"phase"`               // 比赛阶段
    Frozen              bool   `json:"frozen"`              // 是否冻结
    DurationSeconds     int64  `json:"durationSeconds"`     // 持续时间
    StartTimeSeconds    int64  `json:"startTimeSeconds"`    // 开始时间
    RelativeTimeSeconds int64  `json:"relativeTimeSeconds"` // 相对时间
    PreparedBy          string `json:"preparedBy"`          // 出题人
    Difficulty          int    `json:"difficulty"`          // 难度
    Kind                string `json:"kind"`                // 种类
    Country             string `json:"country"`             // 国家
    City                string `json:"city"`                // 城市
    Season              string `json:"season"`              // 赛季
}
```

### 洛谷数据模型

#### LuoguRecord (提交记录)
```go
type LuoguRecord struct {
    Time             int           `json:"time"`             // 运行时间
    Memory           int           `json:"memory"`           // 内存使用
    Problem          LuoguProblem  `json:"problem"`          // 题目信息
    Contest          *LuoguContest `json:"contest"`          // 比赛信息
    SourceCodeLength int           `json:"sourceCodeLength"` // 代码长度
    SubmitTime       int64         `json:"submitTime"`       // 提交时间
    Language         int           `json:"language"`         // 编程语言
    User             LuoguUser     `json:"user"`             // 用户信息
    ID               int64         `json:"id"`               // 记录ID
    Status           int           `json:"status"`           // 状态
    EnableO2         bool          `json:"enableO2"`         // 是否开启O2优化
    Score            int           `json:"score"`            // 得分
}
```

## 反作弊机制

系统内置多种反作弊检测算法：

### 1. 频率限制
- 请求间隔控制
- 并发数量限制
- 每日请求配额

### 2. 用户代理轮换
- 随机User-Agent
- 模拟真实浏览器行为
- 请求头随机化

### 3. 代理池支持
- 支持HTTP/HTTPS代理
- 自动代理轮换
- 代理健康检查

### 4. 验证码处理
- 集成Python OCR
- 自动识别验证码
- 人工干预机制

## 性能优化

### 1. 并发控制
```go
// 配置并发数量
cf:
  cfRecordsConcurrency: 4        // CF记录爬取并发数
  cfTeamContestsConcurrency: 4   // CF团队比赛并发数

luogu:
  luoguRecordsConcurrency: 10    // 洛谷记录爬取并发数
```

### 2. 增量更新
- 基于时间戳的增量爬取
- 避免重复数据获取
- 智能数据去重

### 3. 缓存机制
- 内存缓存热点数据
- 数据库查询优化
- 减少网络请求

## 故障排除

### 常见问题

1. **数据库连接失败**
   ```
   错误: dial tcp 210.36.22.245:3002: connect: connection refused
   解决: 检查数据库服务是否启动，网络是否连通
   ```

2. **洛谷登录失败**
   ```
   错误: 验证码识别失败
   解决: 检查Cookie是否过期，更新用户凭据
   ```

3. **钉钉API调用失败**
   ```
   错误: access_token过期
   解决: 检查appKey和appSecret配置
   ```

4. **ocr配置问题**
   ```
   因为ocr依赖较大，因此不上传，后续需要图像识别可以自主使用其它图像识别
   ```