package db

// Upsert_luogu_sub 幂等写入一条洛谷提交记录，返回本次是否真的新增了一行。
//
// 表的主键是 sub_id，所以"插入撞主键"只说明这条记录已经爬过，不是错误：
//   - 已存在时返回 (false, nil)，调用方据此把它当作"已有"而不是失败，也就不会再有
//     Error 1062 Duplicate entry 刷屏；
//   - ownerKnown 为 true（洛谷记录里带回了提交者 uid）时，顺带把 uid 纠正为洛谷给出的
//     真实提交者，这样 user.luogu_uid 变更后遗留在旧 uid 下的行会被逐步拉回来；
//   - ownerKnown 为 false 时 uid 是兜底值，不能拿来改写库里已有的行。
func Upsert_luogu_sub(table Luogu_all_submissions, ownerKnown bool) (bool, error) {
	query := "insert into luogu_all_submissions (Sub_id, Uid, Is_pass, Creation_time, Problem_name, Difficulty, Problem_id) values (?,?,?,?,?,?,?)"
	args := []interface{}{
		table.Sub_id,
		table.Uid,
		table.Is_pass,
		table.Creation_time,
		table.Problem_name,
		table.Difficulty,
		table.Problem_id,
	}
	if ownerKnown {
		// 刻意不用 values(Uid)：该写法在 MySQL 8.0.20 起已废弃，重复传一次参数兼容所有版本
		query += " on duplicate key update Uid = ?"
		args = append(args, table.Uid)
	} else {
		query += " on duplicate key update Sub_id = Sub_id"
	}

	result, err := Pool.Exec(query, args...)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	// affected：1 = 新插入；0 = 已存在且无需更新；2 = 已存在但纠正了 uid。后两种都算"已爬过"
	return affected == 1, nil
}

// 批量插入，性能更好
func BatchInsert_luogu_sub(tables []Luogu_all_submissions) error {
	if len(tables) == 0 {
		return nil
	}

	query := "insert into luogu_all_submissions (Sub_id, Uid, Is_pass, Creation_time, Problem_name, Difficulty, Problem_id) values "
	args := make([]interface{}, 0, len(tables)*7)

	for i, table := range tables {
		if i > 0 {
			query += ","
		}
		query += "(?,?,?,?,?,?,?)"
		args = append(args, table.Sub_id, table.Uid, table.Is_pass, table.Creation_time, table.Problem_name, table.Difficulty, table.Problem_id)
	}

	_, err := Pool.Exec(query, args...)
	return err
}

// 批量插入CF提交记录
func BatchInsert_cf_all_sub(tables []Cf_all_submissions) error {
	if len(tables) == 0 {
		return nil
	}

	query := "insert into cf_all_submissions (Sub_id, Problem_id, Account, Problem_name, Rating, Verdict, Creation_time) values "
	args := make([]interface{}, 0, len(tables)*7)

	for i, table := range tables {
		if i > 0 {
			query += ","
		}
		query += "(?,?,?,?,?,?,?)"
		args = append(args, table.Sub_id, table.Problem_id, table.Account, table.Problem_name, table.Rating, table.Verdict, table.Creation_time)
	}

	_, err := Pool.Exec(query, args...)
	return err
}

// 批量插入钉钉打卡记录
func BatchInsert_checkup(tables []Ding_checkUp) error {
	if len(tables) == 0 {
		return nil
	}

	query := "insert into ding_checkup (name, ding_id, time, check_type) values "
	args := make([]interface{}, 0, len(tables)*4)

	for i, table := range tables {
		if i > 0 {
			query += ","
		}
		query += "(?,?,?,?)"
		args = append(args, table.Name, table.Ding_id, table.Time, table.Check_type)
	}

	_, err := Pool.Exec(query, args...)
	return err
}

func Insert_checkup(table Ding_checkUp) error {
	_, err := Pool.Exec("insert into ding_checkup (name,ding_id,time,check_type) values (?,?,?,?)",
		table.Name,
		table.Ding_id,
		table.Time,
		table.Check_type,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_cf_all_sub(table Cf_all_submissions) error {
	_, err := Pool.Exec(
		"insert into cf_all_submissions (Sub_id,Problem_id,Account,Problem_name,Rating,Verdict,Creation_time) values(?,?,?,?,?,?,?)",
		table.Sub_id,
		table.Problem_id,
		table.Account,
		table.Problem_name,
		table.Rating,
		table.Verdict,
		table.Creation_time,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_never_pass(table Cf_never_pass) error {
	_, err := Pool.Exec(
		"insert into cf_never_pass (Problem_id,Account,Problem_name,Rating ) values (?,?,?,?)",
		table.Problem_id,
		table.Account,
		table.Problem_name,
		table.Rating,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_cf_official(table Cf_official_problems) error {
	_, err := Pool.Exec(
		"insert into cf_official_problems (Problem_id, title, points, rating, tags) values (?,?,?,?,?)",
		table.Problem_id,
		table.Title,
		table.Points,
		table.Rating,
		table.Tags,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_cf_team_contests(table Cf_team_contests) error {
	_, err := Pool.Exec(
		"insert into cf_team_contests (Contest_id, Contest_name, Start_time, PrePare_by) values (?,?,?,?)",
		table.Contest_id,
		table.Contest_name,
		table.Start_time,
		table.PrePare_by,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_team_questions(table Cf_team_problems) error {
	_, err := Pool.Exec(
		"insert into cf_team_problems (Team_contest_id,Team_contest_name, Official_contest_ID, Problem_name, Rating) values (?,?,?,?,?)",
		table.Team_contest_id,
		table.Team_contest_name,
		table.Official_contest_ID,
		table.Problem_name,
		table.Rating,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_cf_official_contests(table Cf_official_contests) error {
	_, err := Pool.Exec(
		"insert into cf_official_contests (Official_contest_id, Official_contest_name, Phase, Start_time) values (?,?,?,?)",
		table.Official_contest_id,
		table.Official_contest_name,
		table.Phase,
		table.Start_time,
	)
	if err != nil {
		return err
	}
	return nil
}
func Insert_luogu_solutions(table Luogu_solutions) error {
	_, err := Pool.Exec(
		"insert into luogu_solutions (Problem_id, Problem_name, Author_name, Author_uid, Solution_md, Creation_time, Link) values (?,?,?,?,?,?,?)",
		table.Problem_id,
		table.Problem_name,
		table.Author_name,
		table.Author_uid,
		table.Solution_md,
		table.Creation_time,
		table.Link,
	)
	if err != nil {
		return err
	}
	return nil
}
