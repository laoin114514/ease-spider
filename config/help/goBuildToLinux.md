在 Windows 系统上将 Go 文件编译为 Linux 可执行文件，只需设置环境变量指定目标平台即可。以下是详细步骤：

### 方法一：命令行设置环境变量（推荐）
1. **打开命令提示符（CMD）或 PowerShell**
2. **设置目标平台环境变量**：
   ```bash
   # CMD 命令
   set GOOS=linux
   set GOARCH=amd64  # 适用于大多数64位Linux服务器

   # PowerShell 命令
   $env:GOOS = "linux"
   $env:GOARCH = "amd64"
   ```
3. **编译代码**：
   ```bash
   go build -o output_name your_file.go
   ```
   - `output_name`：生成的 Linux 可执行文件名（无后缀）
   - `your_file.go`：你的 Go 源文件

### 方法二：单行命令（适用于 CMD）
```bash
set GOOS=linux&& set GOARCH=amd64&& go build -o myapp main.go
```

### 方法三：使用 `.env` 文件（长期项目）
1. 在项目根目录创建 `.env` 文件：
   ```ini
   GOOS=linux
   GOARCH=amd64
   ```
2. 使用 `go build` 时自动加载环境变量（需第三方工具如 `direnv`）

---

### 关键参数说明
| 环境变量    | 值          | 作用                     |
|-------------|-------------|--------------------------|
| `GOOS`      | `linux`     | 目标操作系统为 Linux     |
| `GOARCH`    | `amd64`     | CPU 架构（64位）         |
| `CGO_ENABLED`| `0`         | 禁用 C 依赖（可选）      |

### 常见架构组合
| Linux 系统类型      | `GOARCH` 值 |
|---------------------|-------------|
| 64 位 Intel/AMD CPU | `amd64`     |
| ARM 64 位 (树莓派4) | `arm64`     |
| ARM 32 位           | `arm`       |

---

### 验证结果
1. 将生成的无后缀文件（如 `myapp`）上传到 Linux 服务器
2. 添加执行权限：
   ```bash
   chmod +x myapp
   ```
3. 运行测试：
   ```bash
   ./myapp
   ```

### 完整示例
```bash
# 编译 main.go 为 Linux 可执行文件
set GOOS=linux
set GOARCH=amd64
go build -o server main.go

# 生成文件: server (无后缀名)
```

> ⚠️ 注意：如果程序使用 CGo，需添加 `set CGO_ENABLED=0` 进行纯静态编译，避免 Linux 兼容性问题。