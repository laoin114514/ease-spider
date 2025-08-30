package dingding

import (
	"fmt"
	"spider/db"
)

func insert(table dingCheckUp, count *int) {
	dataBase := db.New()
	_, err := dataBase.Exec("insert into checkup (name,dingID,time,checkType) values (?,?,?,?)", table.name, table.userId, table.time, table.checkType)
	if err != nil {
		fmt.Println(err)
		return
	}
	*count++
}
