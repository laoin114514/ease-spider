package newLuogu

import "database/sql"

func handle(result []any, dataBase *sql.DB, count *int) error {
	for _, v := range result {
		v := v.(map[string]any)
		var table luogu_all_submissions
		fillTable(v, &table)
		err := insert(&table, dataBase)
		if err != nil {
			return err
		}
		*count++
	}
	return nil
}
