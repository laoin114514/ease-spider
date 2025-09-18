package db

import "errors"

func Insert_luogu_sub(table Luogu_all_submissions) error {
	_, err := Pool.Exec("insert into luogu_all_submissions (Sub_id, Uid,Is_pass, Creation_time, Problem_name, Difficulty, Problem_id) values (?,?,?,?,?,?,?)",
		table.Sub_id,
		table.Uid,
		table.Is_pass,
		table.Creation_time,
		table.Problem_name,
		table.Difficulty,
		table.Problem_id,
	)
	if err != nil {
		return err
	}
	return nil
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
		return errors.New("数据更新完毕")
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

func Insert_team_trainning(table Cf_team_contests) error {
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
