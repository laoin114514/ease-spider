在 Go 语言中，时间操作主要通过 `time` 包实现。以下是详细的常用时间获取函数和方法：

---

### **1. 获取当前时间**
#### `time.Now()`
返回当前本地时间的 `time.Time` 对象（精确到纳秒）。
```go
now := time.Now()
fmt.Println(now) // 输出：2023-10-01 12:34:56.789012345 +0800 CST
```

---

### **2. 获取时间戳**
#### (1) 秒级时间戳
```go
unixSeconds := time.Now().Unix() // 返回 int64
```
#### (2) 毫秒级时间戳
```go
unixMilli := time.Now().UnixMilli() // Go 1.17+
```
#### (3) 微秒级时间戳
```go
unixMicro := time.Now().UnixMicro() // Go 1.17+
```
#### (4) 纳秒级时间戳
```go
unixNano := time.Now().UnixNano() // 返回 int64
```

---

### **3. 获取时间组成部分**
从 `time.Time` 对象中提取具体部分：
```go
t := time.Now()
year := t.Year()       // 年（int）
month := t.Month()     // 月（time.Month 枚举，可转 int）
day := t.Day()         // 日（int）
hour := t.Hour()       // 时（int）
minute := t.Minute()   // 分（int）
second := t.Second()   // 秒（int）
nsec := t.Nanosecond() // 纳秒（int）
weekday := t.Weekday() // 星期（time.Weekday 枚举）
```

---

### **4. 获取时区信息**
```go
loc := time.Now().Location() // 返回 *time.Location
fmt.Println(loc.String())   // 输出时区，如 "Asia/Shanghai"
```

---

### **5. 获取特定时区的时间**
```go
// 加载时区
loc, _ := time.LoadLocation("America/New_York")
nyTime := time.Now().In(loc) // 转换为纽约时间
```

---

### **6. 构造特定时间**
#### `time.Date()`
创建指定时间点的 `time.Time` 对象：
```go
t := time.Date(2023, time.October, 1, 12, 0, 0, 0, time.UTC)
```

---

### **7. 解析字符串时间**
#### `time.Parse(layout, value)`
将字符串解析为 UTC 时间：
```go
t, _ := time.Parse("2006-01-02 15:04:05", "2023-10-01 12:00:00")
```
#### `time.ParseInLocation(layout, value, loc)`
解析为指定时区的时间：
```go
t, _ := time.ParseInLocation("2006-01-02", "2023-10-01", time.Local)
```

> **注意**：布局字符串必须使用参考时间 `2006-01-02 15:04:05`（Go 的诞生时间）。

---

### **8. 获取时间差**
#### `time.Since(t)`
计算从时间 `t` 到现在的时间间隔：
```go
start := time.Now()
// ... 执行操作
elapsed := time.Since(start) // 返回 time.Duration
fmt.Println(elapsed)         // 输出：1.234s
```
#### `time.Until(t)`
计算现在到未来时间 `t` 的剩余时间。

---

### **9. 定时器相关**
#### `time.After(duration)`
返回一个通道，在指定时间后发送当前时间：
```go
select {
case <-time.After(2 * time.Second):
    fmt.Println("2秒后触发")
}
```
#### `time.Tick(duration)`
返回一个通道，每隔指定时间发送一次时间点（用于周期性任务）。

---

### **10. 获取单调时间（避免时钟回拨）**
```go
start := time.Now()
// 使用 time.Since(start) 计算耗时（基于单调时钟，不受系统时间调整影响）
```

---

### **完整示例**
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. 当前时间
	now := time.Now()
	fmt.Println("当前时间:", now)

	// 2. 时间戳
	fmt.Println("秒级时间戳:", now.Unix())
	fmt.Println("毫秒级时间戳:", now.UnixMilli())
	fmt.Println("纳秒级时间戳:", now.UnixNano())

	// 3. 时间组成部分
	fmt.Printf("日期: %d-%02d-%02d\n", now.Year(), now.Month(), now.Day())
	fmt.Printf("时间: %02d:%02d:%02d\n", now.Hour(), now.Minute(), now.Second())

	// 4. 时区信息
	fmt.Println("时区:", now.Location())

	// 5. 特定时区时间
	nyLoc, _ := time.LoadLocation("America/New_York")
	nyTime := now.In(nyLoc)
	fmt.Println("纽约时间:", nyTime)

	// 6. 计算时间差
	start := time.Now()
	time.Sleep(1 * time.Second)
	fmt.Println("耗时:", time.Since(start))
}
```

---

### **关键说明**
1. **时区处理**：  
   Go 默认使用本地时区，可用 `time.Local`、`time.UTC` 或自定义时区。
2. **时间格式化**：  
   使用 `time.Format(layout string)` 方法，布局字符串必须基于 `2006-01-02 15:04:05`。
3. **单调时间**：  
   `time.Now()` 包含单调时钟读数，适合测量时间间隔（如性能分析）。

通过组合这些函数，可以满足绝大多数时间操作需求。详细文档参考：[Go time package](https://pkg.go.dev/time)。