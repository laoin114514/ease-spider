在Go中使用MySQL，通常通过`database/sql`包配合MySQL驱动实现。以下是详细步骤和示例代码：

---

### 1. 安装MySQL驱动
```bash
go get -u github.com/go-sql-driver/mysql
```

---

### 2. 基础使用示例
```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // 匿名导入（只执行init()注册驱动）
)

func main() {
	// 连接数据库
	db, err := sql.Open("mysql", "用户名:密码@tcp(IP:端口)/数据库名?charset=utf8mb4&parseTime=True")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // 确保关闭连接

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("数据库连接失败: ", err)
	}
	fmt.Println("MySQL连接成功")

	// 示例操作
	createTable(db)
	insertData(db)
	queryData(db)
}

// 创建表
func createTable(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE
	)`
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal("建表失败: ", err)
	}
	fmt.Println("表创建成功")
}

// 插入数据
func insertData(db *sql.DB) {
	result, err := db.Exec(
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"张三", "zhangsan@example.com",
	)
	if err != nil {
		log.Fatal("插入失败: ", err)
	}

	id, _ := result.LastInsertId()
	fmt.Printf("插入成功，ID: %d\n", id)
}

// 查询数据
func queryData(db *sql.DB) {
	type User struct {
		ID    int
		Name  string
		Email string
	}

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Fatal("查询失败: ", err)
	}
	defer rows.Close() // 必须关闭

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			log.Fatal("扫描行失败: ", err)
		}
		users = append(users, u)
	}

	fmt.Printf("查询到 %d 条数据:\n", len(users))
	for _, u := range users {
		fmt.Printf("ID: %d, 姓名: %s, 邮箱: %s\n", u.ID, u.Name, u.Email)
	}
}
```

---

### 3. 关键点说明
#### (1) 连接字符串格式
```go
"user:password@tcp(host:port)/dbname?param=value"
```
- **常用参数**：
  - `charset=utf8mb4`：支持完整UTF-8（包括emoji）
  - `parseTime=True`：将数据库时间类型解析为Go的`time.Time`
  - `loc=Local`：时区设置（中国可用`Asia/Shanghai`）

#### (2) 预编译语句（防SQL注入）
使用 `?` 作为占位符：
```go
db.Exec("UPDATE users SET name=? WHERE id=?", "李四", 1)
```

#### (3) 事务处理
```go
tx, err := db.Begin()
if err != nil {
    log.Fatal(err)
}

_, err = tx.Exec("INSERT INTO users (...) VALUES (...)")
if err != nil {
    tx.Rollback() // 回滚
    log.Fatal(err)
}

err = tx.Commit() // 提交
if err != nil {
    log.Fatal(err)
}
```

#### (4) 错误处理
- 特殊错误检查：
```go
if err == sql.ErrNoRows {
    // 查询结果为空
}
```

---

### 4. 连接池配置
```go
db.SetMaxOpenConns(25)      // 最大连接数
db.SetMaxIdleConns(10)      // 最大空闲连接
db.SetConnMaxLifetime(time.Minute * 5) // 连接最大存活时间
```

---

### 5. 推荐工具库
- **[sqlx](https://github.com/jmoiron/sqlx)**：简化数据扫描和常用操作
- **[gorm](https://gorm.io)**：全功能ORM框架

---

### 常见问题解决
1. **连接超时**：
   - 检查防火墙/端口（默认3306）
   - 确认MySQL允许远程连接（`bind-address=0.0.0.0`）
2. **时区问题**：
   - 连接字符串添加 `loc=Asia%2FShanghai`（需URL编码）
3. **驱动未注册**：
   - 确保匿名导入 `_ "github.com/go-sql-driver/mysql"`

> 提示：生产环境避免硬编码密码，使用环境变量或配置管理工具（如Vault）。