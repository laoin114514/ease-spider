在 Go 语言中使用 SHA-512 哈希算法非常简单，标准库 `crypto/sha512` 提供了完整的实现。以下是详细使用指南：

---

### 1. 基本用法（处理小数据）
```go
package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func main() {
	// 原始数据
	data := []byte("Hello, SHA-512!")

	// 计算 SHA-512 哈希值
	hash := sha512.Sum512(data) // 返回 [64]byte 数组

	// 转换为十六进制字符串
	hashStr := hex.EncodeToString(hash[:])
	fmt.Println("SHA-512:", hashStr)
}
```

---

### 2. 流式处理（适合大文件/数据流）
```go
package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {
	// 打开文件
	file, err := os.Open("largefile.dat")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建 SHA-512 哈希器
	hasher := sha512.New()

	// 将文件内容写入哈希器
	if _, err := io.Copy(hasher, file); err != nil {
		panic(err)
	}

	// 计算最终哈希值
	hash := hasher.Sum(nil)

	// 输出结果
	fmt.Println("File SHA-512:", hex.EncodeToString(hash))
}
```

---

### 3. 使用 SHA-512 变种
Go 还支持 SHA512/224 和 SHA512/256 两种变体：
```go
func main() {
	data := []byte("Hello")

	// SHA-512/224
	hash224 := sha512.Sum512_224(data)
	fmt.Println("SHA-512/224:", hex.EncodeToString(hash224[:]))

	// SHA-512/256
	hash256 := sha512.Sum512_256(data)
	fmt.Println("SHA-512/256:", hex.EncodeToString(hash256[:]))
}
```

---

### 4. 关键函数说明
| 函数/方法                     | 作用                                                                 |
|-------------------------------|----------------------------------------------------------------------|
| `sha512.Sum512(data []byte)`  | 直接计算数据的 SHA-512，返回 `[64]byte` 数组                         |
| `sha512.New()`                | 创建 `hash.Hash` 对象，用于流式处理                                  |
| `h.Write(data []byte)`        | 向哈希器追加数据（可多次调用）                                       |
| `h.Sum([]byte)`               | 计算当前哈希值，可选追加额外数据                                     |
| `sha512.Sum512_224()`         | 计算 SHA-512/224（输出 28 字节）                                     |
| `sha512.Sum512_256()`         | 计算 SHA-512/256（输出 32 字节）                                     |

---

### 5. 注意事项
1. **输出长度**：
   - SHA-512: **64字节** (128字符十六进制)
   - SHA-512/224: **28字节** (56字符十六进制)
   - SHA-512/256: **32字节** (64字符十六进制)

2. **性能建议**：
   - 小数据直接用 `Sum512()`
   - 大数据（如文件）使用 `io.Copy` 流式处理

3. **安全性**：
   - SHA-512 是抗碰撞的加密哈希函数
   - 适用于密码存储（需加盐）、数据完整性校验等场景

---

### 6. 完整示例（带盐值密码哈希）
```go
func hashPassword(password string, salt []byte) string {
	// 拼接密码和盐值
	pwdBytes := []byte(password)
	salted := append(pwdBytes, salt...)

	// 计算哈希
	hash := sha512.Sum512(salted)
	return hex.EncodeToString(hash[:])
}

func main() {
	password := "MySecret123"
	salt := []byte{0x12, 0x34, 0x56} // 实际应使用随机盐

	hashed := hashPassword(password, salt)
	fmt.Println("Hashed password:", hashed)
}
```

> **安全提示**：实际密码存储应使用 `bcrypt` 或 `scrypt` 等专门算法，SHA-512 需配合随机盐和多次迭代使用。

通过以上方法，您可以灵活地在 Go 项目中实现 SHA-512 哈希功能。