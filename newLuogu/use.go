package newLuogu

import (
	"fmt"
	"spider/db"
)

func Use() {
	fmt.Println("洛谷提交情况:")
	dataBase := db.New()
	defer dataBase.Close()
	rows, err := dataBase.Query("select id from luoguac")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var uid string
		rows.Scan(&uid)
		result := request(uid, 1)
		if result == nil {
			continue
		}
		page := countPage(result)
		count := 0
		for i := 1; i <= page; i++ {
			data := request(uid, i)
			result, ok := data["result"].([]any)
			if !ok {
				continue
			}
			err := handle(result, dataBase, &count)
			if err != nil {
				break
			}
		}
		fmt.Println(uid, count)
	}
}

// {
//     time: 41,
//     memory: 832,
//     problem: {
//       pid: 'P1113',
//       title: '杂务',
//       difficulty: 3,
//       fullScore: 100,
//       type: 'P'
//     },
//     contest: null,
//     sourceCodeLength: 453,
//     submitTime: 1756215408,
//     language: 28,
//     user: {
//       uid: 1811873,
//       name: 'lllllllke',
//       avatar: 'https://cdn.luogu.com.cn/upload/usericon/1811873.png',
//       slogan: '',
//       badge: null,
//       isAdmin: false,
//       isBanned: false,
//       color: 'Blue',
//       ccfLevel: 0,
//       xcpcLevel: 0,
//       background: ''
//     },
//     id: 233308785,
//     status: 14,
//     enableO2: true,
//     score: 0
//   }
