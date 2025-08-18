package db

type User struct {
	Id          int
	Account     string
	Password    string
	Username    string
	Email       string
	Ranking     int
	Role_id     int
	Create_time string
	Update_time string
	School      string
}

type Cf_all_submissions struct {
	SubId        int
	ProblemId    string
	Handle       string
	ProblemName  string
	Rating       int
	Verdict      string
	CreationTime string
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
	StartTime string
	PrePareBy string
}
