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

type Cf_official_problems struct {
	Problem_id string
	Title      string
	Points     int
	Rating     int
	Tags       string
}
type Cf_team_problems struct {
	Team_contest_id     int
	Team_contest_name   string
	Official_contest_ID string
	Problem_name        string
	Rating              int
}

type Cf_team_contests struct {
	Contest_id   int
	Contest_name string
	Start_time   time.Time
	PrePare_by   string
}
type Cf_never_pass struct {
	Account      string
	Problem_id   string
	Problem_name string
	Rating       int
}

type Luogu_all_submissions struct {
	Sub_id        string
	Uid           string
	Problem_id    string
	Problem_name  string
	Difficulty    string
	Is_pass       bool
	Creation_time time.Time
}
type Ding_checkUp struct {
	Name       string
	Ding_id    string
	Time       time.Time
	Check_type string
}
type Cf_official_contests struct {
	Official_contest_id   int
	Official_contest_name string
	Phase                 string
	Start_time            time.Time
}
