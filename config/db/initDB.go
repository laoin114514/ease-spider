package db

import (
	"database/sql"
	"fmt"
	"log"
	"spider/config"

	_ "github.com/go-sql-driver/mysql"
)

var Pool *sql.DB

func Init() error {
	dataSource := fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True",
		config.AppConfig.DB_USER,
		config.AppConfig.DB_PASSWORD,
		config.AppConfig.DB_HOST,
		config.AppConfig.DB_NAME,
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
	log.Printf(
		"数据库连接成功,配置信息: \n"+
			"    用户名: %s\n"+
			"    密码: %s\n"+
			"    数据库: %s\n"+
			"    主机: %s\n",
		config.AppConfig.DB_USER,
		config.AppConfig.DB_PASSWORD,
		config.AppConfig.DB_NAME,
		config.AppConfig.DB_HOST,
	)
	return nil
}
