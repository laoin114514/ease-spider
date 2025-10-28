# Utils 工具目录使用说明

## 1. 并发器 (conCurrenter.go)

### 功能
控制并发执行任务，支持超时机制。

### 使用方法
```go
// 创建并发器（默认超时3600秒）
concurrenter := utils.NewConCurrenter[int](5) // 5个并发

// 创建带自定义超时的并发器
concurrenter := utils.NewConCurrenterWithTimeout[int](5, 30*time.Second)

// 执行任务
params := []int{1, 2, 3, 4, 5}
err := concurrenter.Run(params, func(param int) error {
    // 处理单个参数
    fmt.Println("处理参数:", param)
    return nil
})
```

## 2. 配置验证器 (config_validator.go)

### 功能
验证应用配置的完整性和有效性。

### 使用方法
```go
// 创建验证器
validator := utils.NewConfigValidator()

// 验证配置
err := validator.ValidateConfig()
if err != nil {
    log.Fatal("配置验证失败:", err)
}
```

## 3. 调试器 (debug.go)

### 功能
提供调试输出功能，可控制是否启用。

### 使用方法
```go
// 创建调试器
debug := utils.NewDebug(true) // true启用调试，false禁用

// 调试输出
debug.Debug("这是调试信息", 123, "test")
debug.Debugf("格式化输出: %s %d", "hello", 456)
```

## 4. CF URL生成器 (generateCFurl.go)

### 功能
生成Codeforces API的请求URL，支持API Key认证。

### 使用方法
```go
// 初始化
utils.InitGenerateCFurl()

// 生成用户状态查询URL
url, err := utils.GenerateCFurlInstance.User.Status(true, &models.UserStatusParams{
    Handle: "tourist",
    From:   1,
    Count:  10,
})

// 生成比赛列表URL
url, err := utils.GenerateCFurlInstance.Contest.List("handle", &models.ContestListParams{
    Gym: false,
})
```

## 5. 哈希编码器 (hashEncode.go)

### 功能
提供SHA512系列哈希编码功能。

### 使用方法
```go
// 创建编码器
encoder := utils.NewHashEncoder()

// SHA512编码
hash := encoder.Hash512("hello world")

// SHA512/224编码
hash224 := encoder.Hash512_224("hello world")

// SHA512/256编码
hash256 := encoder.Hash512_256("hello world")
```

## 6. JSON数据库 (josnDB.go)

### 功能
提供基于JSON文件的数据库操作和内存数据库。

### 使用方法
```go
// JSON文件数据库
utils.InitGlobalJSONDB("data.json")
utils.JsonDB.Set("key", "value")
value := utils.JsonDB.Get("key")
utils.JsonDB.Delete("key")

// 内存数据库
utils.InitGlobalRAMDB()
utils.RamDB.Set("key", "value")
value := utils.RamDB.Get("key")
utils.RamDB.Delete("key")
```

## 7. 日志容器 (logContainer.go)

### 功能
收集和管理日志信息，支持普通日志和错误日志。

### 使用方法
```go
// 创建日志容器
logContainer := utils.NewLogContainer()

// 添加日志
logContainer.AddLog("这是一条普通日志")
logContainer.AddErr("这是一条错误日志")

// 获取日志
logs := logContainer.GetLog()
errs := logContainer.GetErr()

// 清空日志
logContainer.ClearLog()
logContainer.ClearErr()
```

## 8. Markdown解析器 (markdownPaser.go)

### 功能
解析Markdown文档中的代码块。

### 使用方法
```go
// 创建解析器
parser := utils.NewMarkdownParser(markdownContent)

// 提取所有代码块
codeBlocks := parser.ExtractCodeBlocks()

// 按语言提取代码块
goBlocks := parser.ExtractCodeBlocksByLanguage("go")

// 获取代码内容数组
contents := parser.GetCodeContentArray()

// 按语言获取代码内容
goContents := parser.GetCodeContentArrayByLanguage("go")

// 统计代码块数量
totalCount := parser.CountCodeBlocks()
goCount := parser.CountCodeBlocksByLanguage("go")

// 获取代码块详细信息
info := parser.GetCodeBlockInfo()
```

## 9. HTTP请求器 (request.go)

### 功能
封装HTTP请求，支持GET/POST，自动JSON解析。

### 使用方法
```go
// 创建请求器
request := utils.NewRequest[ResponseType](true) // true启用严格JSON解析

// GET请求
response, err := request.Get("https://api.example.com/data", map[string]string{
    "param1": "value1",
    "param2": "value2",
})

// POST请求
response, err := request.Post("https://api.example.com/submit", requestBody)

// 设置Cookie
request.SetCookie("session=abc123")

// 设置自定义Header
request.SetHeader("Authorization", "Bearer token")

// 获取原始响应
rawResp := request.GetRawResp()
rawBody := request.GetRawRespBody()
```

## 10. 结构体转换器 (structTransfer.go)

### 功能
将结构体转换为Map和有序参数字符串。

### 使用方法
```go
// 创建转换器
transfer := utils.NewStructTransfer(myStruct)

// 转换为Map
dataMap, err := transfer.ToMap()

// 转换为有序参数字符串
paramStr, err := transfer.ToOrderStr()

// 添加参数
transfer.AddParam("newKey", "newValue")

// 批量添加参数
transfer.AddParamWithMap(map[string]any{
    "key1": "value1",
    "key2": "value2",
})
```

## 11. 计时器 (timer.go)

### 功能
提供函数执行时间测量和定时任务功能。

### 使用方法
```go
// 创建计时器
timer := utils.NewTimer()

// 测量函数执行时间（毫秒）
duration := timer.CountDurationInMs(func() {
    // 要测量的代码
    time.Sleep(100 * time.Millisecond)
})

// 测量函数执行时间（秒）
duration := timer.CountDurationInS(func() {
    // 要测量的代码
})

// 获取执行时间字符串
durationStr := timer.CountDurationStr(func() {
    // 要测量的代码
})

// 定时执行任务
timer.RunWithTimer(5*time.Second, "数据同步", func() error {
    // 定时任务逻辑
    return nil
})

// 多任务定时执行
timer.MultiRunWithTimer(10*time.Second, 
    func() error { return nil }, // 任务1
    func() error { return nil }, // 任务2
)

// 日期格式化
dateFormat := utils.NewDateFormat()

// 格式化时间戳
formatted := dateFormat.DateTimeWithSecond(1640995200)

// 当前时间
now := dateFormat.NowDateTime()

// 当前日期
today := dateFormat.NowDate()

// 指定天数前的时间
before := dateFormat.BeforDateTimeWithDay(7)

// 解析定时频率字符串
duration := dateFormat.AnalysisTimerFrequency("30s") // 30秒
duration = dateFormat.AnalysisTimerFrequency("5m")  // 5分钟
duration = dateFormat.AnalysisTimerFrequency("2h")  // 2小时
duration = dateFormat.AnalysisTimerFrequency("1d")  // 1天
```

## 使用建议

1. **并发器**: 适用于需要控制并发数量的场景，如批量API请求
2. **配置验证器**: 在应用启动时验证配置完整性
3. **调试器**: 开发阶段使用，生产环境建议关闭
4. **CF URL生成器**: 专门用于Codeforces API调用
5. **哈希编码器**: 用于密码加密、数据校验等场景
6. **JSON数据库**: 轻量级数据存储，适合配置和小数据量场景
7. **日志容器**: 收集运行时日志，便于后续分析和报告
8. **Markdown解析器**: 处理包含代码的Markdown文档
9. **HTTP请求器**: 统一的HTTP请求接口，自动处理JSON解析
10. **结构体转换器**: API参数构建、数据格式转换
11. **计时器**: 性能测试、定时任务、时间格式化
