package luogu

import (
	"database/sql"
	"time"
)

func insertPass(table luoguAc, dataBase *sql.DB, todayCount int) error {

	_, err := dataBase.Exec(
		"update luoguac set red=?,brown=?,yellow=?,green=?,blue=?,purple=?,black=?,todayCount=? where id=?",
		table.red,
		table.brown,
		table.yellow,
		table.green,
		table.blue,
		table.purple,
		table.black,
		todayCount,
		table.id,
	)
	if err != nil {
		return err
	}
	return nil
}
func insertDayPass(table luoguAc, dataBase *sql.DB, todayCount int) error {
	now := time.Now()
	_, err := dataBase.Exec("insert into luogudayac (date,luogu_uid,count) values (?,?,?)", time.Unix(now.Unix(), 0), table.id, todayCount)
	if err != nil {
		return err
	}
	return nil
}
