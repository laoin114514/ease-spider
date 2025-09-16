在 Go 中读取 `.env` 文件通常使用第三方库 **`github.com/joho/godotenv`**，以下是完整实现步骤：

### 1. 安装库
```bash
go get github.com/joho/godotenv
```

### 2. 创建 `.env` 文件
在项目根目录创建文件，内容示例：
```env
APP_PORT=8080
DB_URL=postgres://user:pass@localhost/dbname
DEBUG=true
```

### 3. 读取环境变量的代码示例
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatal("无法加载 .env 文件")
	}

	// 获取环境变量（方法1：直接读取）
	port := os.Getenv("APP_PORT")
	fmt.Println("APP_PORT:", port) // 输出: 8080

	// 获取环境变量（方法2：带默认值）
	dbUrl := getEnv("DB_URL", "default-db-url")
	fmt.Println("DB_URL:", dbUrl)

	// 获取布尔值
	debug := os.Getenv("DEBUG") == "true"
	fmt.Println("DEBUG模式:", debug)
}

// 辅助函数：获取环境变量，带默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
```

### 关键说明：
1. **加载文件**：
   ```go
   godotenv.Load() // 加载当前目录的 .env 文件
   ```
   - 指定路径：`godotenv.Load("config/.env")`
   - 加载多个文件：`godotenv.Load(".env", ".env.local")`（后面的文件优先级更高）

2. **读取变量**：
   ```go
   os.Getenv("KEY") // 不存在时返回空字符串
   ```

3. **高级选项**：
   ```go
   // 带展开功能（解析 $VAR 变量）
   godotenv.LoadWithOptions(".env", godotenv.Expand(true))
   
   // 覆盖已存在的系统环境变量
   godotenv.Overload() 
   ```

### 生产环境建议：
- **不要将 `.env` 提交到 Git**（添加到 `.gitignore`）
- 生产环境直接使用系统环境变量（Docker/K8s/云平台等）
- 开发环境建议使用 `.env.local` 覆盖个人配置

### 替代方案：
如果需要解析复杂结构，可考虑：
1. **`github.com/caarlos0/env`**（直接绑定到结构体）
2. **`github.com/spf13/viper`**（支持多种配置格式）

> 使用示例：https://github.com/joho/godotenv

按照这个流程即可在 Go 中安全高效地管理环境变量配置。