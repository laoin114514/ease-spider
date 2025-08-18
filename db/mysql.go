package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func New() *sql.DB {
	godotenv.Load()
	account := os.Getenv("ACCOUNT")
	password := os.Getenv("PASSWORD")
	url := os.Getenv("URL")
	dateBase := os.Getenv("DATEBASE")
	dataSource := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True", account, password, url, dateBase)
	db, err := sql.Open(
		"mysql",
		dataSource,
	)
	if err != nil {
		fmt.Println(err)
	}
	if err := db.Ping(); err != nil {
		fmt.Println("数据库连接失败: ", err)
	}
	return db
}
