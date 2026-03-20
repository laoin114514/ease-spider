package models

// NiukeUser 牛客用户信息
type NiukeUser struct {
	RealName string `json:"real_name"`
	Uid      string `json:"uid"`
}

// NiukeUserDeliver 传递给爬虫的用户信息
type NiukeUserDeliver struct {
	RealName   string          `json:"real_name"`
	Uid        string          `json:"uid"`
	Count      int             `json:"count"`
	OldDataSet map[string]bool `json:"old_data_set"` // 已爬取的提交ID集合，用于去重
}

// NiukeSubmission 牛客提交记录
type NiukeSubmission struct {
	SubId        string `json:"sub_id"`
	Uid          string `json:"uid"`
	UserName     string `json:"user_name"`
	ProblemId    string `json:"problem_id"`
	ProblemName  string `json:"problem_name"`
	SubmitTime   string `json:"submit_time"`
	Status       string `json:"status"`
	Language     string `json:"language"`
}
