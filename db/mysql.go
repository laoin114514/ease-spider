package db

import (
	"database/sql"
	"fmt"
	"os"

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
