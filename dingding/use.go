package dingding

import (
	"fmt"
	"spider/component"
)

func Use() {
	token := getToken()
	var from int64 = 0
	checkUpDatas := request(token, component.BeforDateTime(from+7), component.BeforDateTime(from))
	count := 0
	for _, v := range checkUpDatas {
		v := v.(map[string]any)
		var table dingCheckUp
		fillTable(v, &table)
		insert(table, &count)
	}
	fmt.Printf("插入%d条打卡记录\n", count)
}
