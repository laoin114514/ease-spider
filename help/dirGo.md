在 Go 语言中读取本地文件数据主要有以下几种方式，我将详细介绍每种方法的使用场景和示例代码：

---

### 1. 一次性读取整个文件（小文件）
使用 `os.ReadFile()` 函数（Go 1.16+ 推荐）
```go
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// 一次性读取文件所有内容（返回 []byte）
	content, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal("文件读取失败: ", err)
	}

	// 转换为字符串输出
	fmt.Println("文件内容:")
	fmt.Println(string(content))
}
```

---

### 2. 按行读取（大文件推荐）
使用 `bufio.Scanner` 逐行处理
```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// 打开文件
	file, err := os.Open("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // 确保关闭文件

	// 创建扫描器逐行读取
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() { // 每次扫描一行
		line := scanner.Text() // 获取文本内容
		fmt.Printf("行号 %d: %s\n", lineNum, line)
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("扫描错误: ", err)
	}
}
```

---

### 3. 分块读取（自定义缓冲区）
使用 `os.Open` + `bufio.NewReader`
```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("largefile.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 4096) // 4KB 缓冲区

	for {
		// 读取数据到缓冲区
		n, err := reader.Read(buffer)
		if err != nil {
			break // 读到文件末尾
		}

		// 处理当前读取块
		fmt.Printf("读取 %d 字节: %q\n", n, string(buffer[:n]))
	}
}
```

---

### 4. 使用 ioutil（旧版，Go 1.16 前）
```go
// 注意：Go 1.16+ 已弃用，建议改用 os.ReadFile
content, err := ioutil.ReadFile("data.txt")
```

---

### 关键注意事项：
1. **文件路径**：
   - 相对路径：相对于可执行文件位置
   - 绝对路径：`/home/user/file.txt`（Linux）或 `C:\file.txt`（Windows）

2. **错误处理**：
   - 必须检查所有文件操作的错误
   - 使用 `log.Fatal` 或自定义错误处理

3. **资源释放**：
   - 使用 `defer file.Close()` 确保文件描述符关闭
   - 大文件操作避免一次性读取

4. **文件权限**：
   - 程序需要文件读取权限
   - 使用 `os.OpenFile()` 可控制打开模式

---

### 完整流程示例（带错误检查）
```go
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// 安全打开文件
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 创建扫描器
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // 处理每行数据
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取错误: %v\n", err)
	}
}
```

---

### 选择建议：
| 场景              | 推荐方法               |
|-------------------|-----------------------|
| 小文件 (<10MB)    | `os.ReadFile()`       |
| 日志文件/逐行处理 | `bufio.Scanner`       |
| 二进制大文件      | 分块读取 (`Reader.Read`) |
| 需要精细控制      | `os.Open` + `bufio.Reader` |

根据实际需求选择合适的方法，大文件务必使用流式读取避免内存溢出！