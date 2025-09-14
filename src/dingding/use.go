package dingding

import (
	"fmt"
	"spider/component"
	"spider/db"
)

var tempDB component.TempDB
var obj map[string]any

func Use() {
	token := getToken()
	count := 0
	for i := int64(0); i < 60; i++ {
		var from int64 = i * 7
		checkUpDatas, err := request(token, component.BeforDateTime(from+7), component.BeforDateTime(from))
		if err != nil {
			break
		}
		obj = tempDB.Get("dingUserId").(map[string]any)
		for _, v := range checkUpDatas {
			v := v.(map[string]any)
			var table db.DingCheckUp
			fillTable(v, &table)
			err := db.Insert_checkup(table)
			if err != nil {
				continue
			}
			fmt.Println(table)
			count++
		}
	}
	fmt.Printf("插入%d条打卡记录\n", count)
}
