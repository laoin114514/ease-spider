package luogu

import "database/sql"

func handle(passProblems []any, table luoguAc) int {
	offset := 0
	for _, v := range passProblems {
		difficulty := v.(map[string]any)["difficulty"].(float64)
		if difficulty == 0 {
			offset--
		} else if difficulty == 1 {
			table.red++
		} else if difficulty == 2 {
			table.brown++
		} else if difficulty == 3 {
			table.yellow++
		} else if difficulty == 4 {
			table.green++
		} else if difficulty == 5 {
			table.blue++
		} else if difficulty == 6 {
			table.purple++
		} else if difficulty == 7 {
			table.black++
		}
	}
	count := len(passProblems) + offset
	return count
}
func todayCount(dataBase *sql.DB, uid string, currentCount int) int {
	rows, queryErr := dataBase.Query("select red,brown,yellow,green,blue,purple,black from luoguac where id=?", uid)
	if queryErr != nil {
		return 0
	}
	formatCount := 0
	for rows.Next() {
		var formatTable luoguAc
		rows.Scan(
			&formatTable.red,
			&formatTable.brown,
			&formatTable.yellow,
			&formatTable.green,
			&formatTable.blue,
			&formatTable.purple,
			&formatTable.black,
		)
		formatCount += formatTable.red + formatTable.brown + formatTable.yellow + formatTable.green + formatTable.blue + formatTable.purple + formatTable.black
	}
	return currentCount - formatCount
}
