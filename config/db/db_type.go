package db

import "time"

type User struct {
	Id          int
	Account     string
	Password    string
	Real_name   string
	Email       string
	School      string
	Role_id     int
	Enter_time  string
	Retire_time string
	Create_time string
}

type Oj_account struct {
	User_id          int
	Cf_account       string
	Cf_apikey        string
	Cf_secret        string
	Luogu_uid        string
	Nowcoder_account string
	Vjudge_account   string
}

type Cf_all_submissions struct {
	Sub_id        int
	Account       string
	Problem_id    string
	Problem_name  string
	Rating        int
	Verdict       string
	Creation_time time.Time
}

type Cf_contest_official struct {
	Id     string
	Title  string
	Points int
	Rating int
	Tags   string
}

type Cf_never_pass struct {
	ProblemId   string
	Handle      string
	ProblemName string
	Rating      int
}

type Cf_team_question struct {
	In_team_ID    int
	Contest_name  string
	Official_ID   string
	Question_name string
	Rating        int
}

type Cf_team_trainning struct {
	Id        int
	Name      string
	StartTime time.Time
	PrePareBy string
}

type Luogu_all_submissions struct {
	SubId       string
	Pid         string
	Username    string
	Uid         string
	IsPass      bool
	SubTime     time.Time
	ProblemName string
	Difficulty  string
}
type DingCheckUp struct {
	Name      string
	UserId    string
	Time      time.Time
	CheckType string
}
