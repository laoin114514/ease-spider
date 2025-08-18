在 Go 语言中，将 JSON 字符串转换为 `map` 或结构体可以使用标准库 `encoding/json` 的 `Unmarshal` 函数。以下是两种方法的详细说明和示例：

---

### 1. 转换为 `map[string]interface{}`
适用于动态或未知结构的 JSON 数据：
```go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonStr := `{"name":"Alice","age":25,"is_student":true,"courses":["math","physics"]}`

	// 声明一个空 map
	var data map[string]interface{}

	// 解析 JSON
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		panic(err)
	}

	// 访问数据
	fmt.Println("Name:", data["name"])       // Alice
	fmt.Println("Age:", data["age"])         // 25 (float64 类型)
	fmt.Println("Is Student:", data["is_student"]) // true

	// 处理嵌套数组
	courses := data["courses"].([]interface{})
	fmt.Println("First Course:", courses[0]) // math
}
```

---

### 2. 转换为结构体 (推荐)
适用于固定结构的 JSON 数据（类型安全）：
```go
package main

import (
	"encoding/json"
	"fmt"
)

// 定义匹配 JSON 的结构体
type User struct {
	Name     string   `json:"name"`      // 字段标签指定 JSON 键名
	Age      int      `json:"age"`       // 类型自动转换
	IsStudent bool     `json:"is_student"`
	Courses  []string `json:"courses"`
}

func main() {
	jsonStr := `{"name":"Bob","age":30,"is_student":false,"courses":["art","history"]}`

	// 声明结构体实例
	var user User

	// 解析 JSON
	err := json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		panic(err)
	}

	// 直接访问字段
	fmt.Printf("User: %+v\n", user)
	// 输出: {Name:Bob Age:30 IsStudent:false Courses:[art history]}
	
	fmt.Println("First Course:", user.Courses[0]) // art
}
```

---

### 关键注意事项：
1. **字段导出**  
   结构体字段必须**首字母大写**（导出），否则 `json.Unmarshal` 无法访问。

2. **JSON 标签**  
   使用 `` `json:"field_name"` `` 标签指定 JSON 字段名（可处理蛇形命名/驼峰命名转换）：
   ```go
   type Example struct {
       FullName string `json:"full_name"` // 映射 JSON 中的 full_name
   }
   ```

3. **类型匹配**  
   JSON 类型与 Go 类型需兼容：
   - JSON 数字 → `int`, `float64`
   - JSON 数组 → `[]T` 或 `[]interface{}`
   - JSON 对象 → `map` 或嵌套结构体

4. **忽略字段**  
   添加 `json:"-"` 标签忽略字段：
   ```go
   type Example struct {
       IgnoreField string `json:"-"` // 跳过此字段
   }
   ```

5. **处理动态数据**  
   对于部分未知的 JSON，可混合使用结构体和 `map`：
   ```go
   type Mixed struct {
       Name string                 `json:"name"`
       Extras map[string]interface{} `json:"extras"` // 存储未知字段
   }
   ```

---

### 完整工作流程：
```go
jsonString := `{"id":123,"user_info":{"name":"Charlie"}}`

// 定义嵌套结构体
type UserInfo struct {
    Name string `json:"name"`
}

type Data struct {
    ID       int      `json:"id"`
    UserInfo UserInfo `json:"user_info"`
}

var d Data
if err := json.Unmarshal([]byte(jsonString), &d); err != nil {
    log.Fatal(err)
}
fmt.Println(d.UserInfo.Name) // Charlie
```

通过以上方法，你可以灵活地在 JSON 字符串与 Go 数据结构之间进行转换。结构体方式在类型安全和代码可读性上更优，而 `map` 方式则适合处理动态数据。