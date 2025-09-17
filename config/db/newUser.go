package db

import (
	"fmt"
)

func NewUser() {
	var user User
	var uid string
	fmt.Scan(&user.Account, &user.Real_name, &uid)
	_, err := Pool.Exec("insert into user (account,Real_name) values (?,?)", user.Account, user.Real_name)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := Pool.Query("select id from user where Real_name=?", user.Real_name)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		_, err := Pool.Exec("insert into oj_account (user_id,luogu_uid) values (?,?)", id, uid)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
