package db

import "errors"

func Insert_luogu_sub(table Luogu_all_submissions) error {
	_, err := Pool.Exec("insert into luogu_all_submissions (subId, username, uid, isPass, subTime, problemName, difficulty, pid) values (?,?,?,?,?,?,?,?)",
		table.SubId,
		table.Username,
		table.Uid,
		table.IsPass,
		table.SubTime,
		table.ProblemName,
		table.Difficulty,
		table.Pid,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_checkup(table DingCheckUp) error {
	_, err := Pool.Exec("insert into checkup (name,dingID,time,checkType) values (?,?,?,?)",
		table.Name,
		table.UserId,
		table.Time,
		table.CheckType,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_cf_all_sub(table Cf_all_submissions) error {
	_, err := Pool.Exec(
		"insert into cf_all_submissions (subId,problemId,handle,problemName,rating,verdict,creationTime) values(?,?,?,?,?,?,?)",
		table.SubId,
		table.ProblemId,
		table.Handle,
		table.ProblemName,
		table.Rating,
		table.Verdict,
		table.CreationTime,
	)
	if err != nil {
		return errors.New("数据更新完毕")
	}
	return nil
}

func Insert_never_pass(table Cf_never_pass) error {
	_, err := Pool.Exec(
		"insert into cf_never_pass (ProblemId,Handle,ProblemName,Rating ) values (?,?,?,?)",
		table.ProblemId,
		table.Handle,
		table.ProblemName,
		table.Rating,
	)
	if err != nil {
		return err
	}
	return nil
}

func Insert_cf_official(table Cf_contest_official) error {
	_, err := Pool.Exec(
		"insert into cf_contest_official (id, title, points, rating, tags) values (?,?,?,?,?)",
		table.Id,
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

func Insert_team_trainning(table Cf_team_trainning) error {
	_, err := Pool.Exec(
		"insert into cf_team_training (id, name, startTime, prepareBy) values (?,?,?,?)",
		table.Id,
		table.Name,
		table.StartTime,
		table.PrePareBy,
	)
	if err != nil {
		return err
	}
	return nil
}
