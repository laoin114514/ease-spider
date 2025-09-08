package luogu

import (
	"fmt"
	"spider/component"
	"spider/db"
)

var tempDB component.TempDB
var pass int

func Use() {
	fmt.Println("洛谷提交情况:")
	rows, err := db.Pool.Query("SELECT u.username,p.luogu FROM user as u,platform_id as p WHERE u.id=p.user_id&&(u.role_id=1||u.role_id=3);")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var username string
		var uid string
		rows.Scan(&username, &uid)
		result := request(uid, 1)
		if uid == "" {
			fmt.Println(username + "uid不存在")
			continue
		}
		if result == nil {
			fmt.Println(username + "uid不存在")
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
			err := handle(result, &count)
			if err != nil {
				break
			}
		}
		fmt.Printf("新增提交%d %s\n", count, username)
	}
	fmt.Printf("\n")
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
