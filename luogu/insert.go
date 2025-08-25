package luogu

import (
	"database/sql"
	"spider/db"
)

func insertPass(table db.Luogu_problem, dataBase *sql.DB) error {
	_, err := dataBase.Exec(
		"insert into luogu_pass_sub (uid,pid,title,difficulty,type) values(?,?,?,?,?)",
		table.Uid,
		table.Pid,
		table.Title,
		table.Difficulty,
		table.Type,
	)
	if err != nil {
		return err
	}
	return nil
}
