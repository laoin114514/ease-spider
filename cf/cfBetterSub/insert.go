package cfBetterSub

import (
	"database/sql"
	"errors"
	"spider/db"
)

func insert(table db.Cf_all_submissions, dataBase *sql.DB) error {
	_, err := dataBase.Exec(
		"insert into cf_all_submissions (subId,problemId,handle,problemName,rating,verdict) values(?,?,?,?,?,?)",
		table.SubId,
		table.ProblemId,
		table.Handle,
		table.ProblemName,
		table.Rating,
		table.Verdict,
	)
	if err != nil {
		return errors.New("数据更新完毕")
	}
	return nil
}
