package newLuogu

import "database/sql"

func insert(table *luogu_all_submissions, dataBase *sql.DB) error {
	_, err := dataBase.Exec("insert into luogu_all_submissions (subId, username, uid, isPass, subTime, problemName, difficulty, pid) values (?,?,?,?,?,?,?,?)",
		table.subId,
		table.username,
		table.uid,
		table.isPass,
		table.subTime,
		table.problemName,
		table.difficulty,
		table.pid,
	)
	if err != nil {
		return err
	}
	return nil
}
