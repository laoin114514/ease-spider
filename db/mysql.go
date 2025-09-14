package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var Pool *sql.DB

func Init() {
	godotenv.Load()
	account := os.Getenv("ACCOUNT")
	password := os.Getenv("PASSWORD")
	url := os.Getenv("URL")
	mysqldb := os.Getenv("DATEBASE")
	dataSource := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True", account, password, url, mysqldb)
	var err error
	Pool, err = sql.Open(
		"mysql",
		dataSource,
	)
	if err != nil {
		fmt.Println(err)
	}
	if err := Pool.Ping(); err != nil {
		fmt.Println("数据库连接失败: ", err)
	}
	fmt.Println("数据库连接成功")

}
func Backup() error {
	backupPath := "gxuicpc.sql"

	// 执行mysqldump命令
	cmd := exec.Command(
		"mysqldump",
		"-u", "gxuicpc",
		"-p", "gxuicpc",
		"-h", "127.0.0.1",
		"-P", "3002",
		"--databases", "gxuicpc",
		"-r", backupPath,
	)

	// 执行并等待完成
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	// 验证备份文件
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return err
	}

	return nil
}
