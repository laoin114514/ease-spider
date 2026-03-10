package db

import (
	"database/sql"
	"fmt"
	"spider/config"

	_ "github.com/go-sql-driver/mysql"
)

var Pool *sql.DB

func Init() error {
	dataSource := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.DbName,
	)
	var err error
	Pool, err = sql.Open(
		"mysql",
		dataSource,
	)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	if err := Pool.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	return nil
}
