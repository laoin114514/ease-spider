package models

// TrainningDetailResponse 训练详情响应结构
type LuoguTrainningDetailResponse struct {
	Code            int                             `json:"code"`
	CurrentData     LuoguTrainningDetailCurrentData `json:"currentData"`
	CurrentTemplate string                          `json:"currentTemplate"`
	CurrentTheme    interface{}                     `json:"currentTheme"`
	CurrentTime     int64                           `json:"currentTime"`
	CurrentTitle    string                          `json:"currentTitle"`
}

// CurrentData 当前数据结构
type LuoguTrainningDetailCurrentData struct {
	CanEdit          bool                         `json:"canEdit"`
	PrivilegedTeams  []string                     `json:"privilegedTeams"`
	Training         LuoguTrainningDetailTraining `json:"training"`
	TrainingProblems TrainingProblems             `json:"trainingProblems"`
}

// Training 训练信息
type LuoguTrainningDetailTraining struct {
	CreateTime   int64            `json:"createTime"`
	Deadline     interface{}      `json:"deadline"`
	Description  string           `json:"description"`
	ID           int              `json:"id"`
	MarkCount    int              `json:"markCount"`
	Marked       bool             `json:"marked"`
	Name         string           `json:"name"`
	ProblemCount int              `json:"problemCount"`
	Problems     []ProblemElement `json:"problems"`
	Provider     Provider         `json:"provider"`
	Title        string           `json:"title"`
	Type         int              `json:"type"`
	UserScore    interface{}      `json:"userScore"`
}

// ProblemElement 题目元素
type ProblemElement struct {
	Problem ProblemProblem `json:"problem"`
}

// ProblemProblem 题目信息
type ProblemProblem struct {
	Difficulty       int    `json:"difficulty"`
	Flag             int    `json:"flag"`
	FullScore        int    `json:"fullScore"`
	PID              string `json:"pid"`
	Tags             []int  `json:"tags"`
	Title            string `json:"title"`
	TotalAccepted    int    `json:"totalAccepted"`
	TotalSubmit      int    `json:"totalSubmit"`
	Type             string `json:"type"`
	WantsTranslation bool   `json:"wantsTranslation"`
}

// Provider 提供者信息
type Provider struct {
	Avatar     string      `json:"avatar"`
	Background string      `json:"background"`
	Badge      interface{} `json:"badge"`
	CCFLevel   int         `json:"ccfLevel"`
	Color      string      `json:"color"`
	IsAdmin    bool        `json:"isAdmin"`
	IsBanned   bool        `json:"isBanned"`
	Name       string      `json:"name"`
	Slogan     string      `json:"slogan"`
	UID        int         `json:"uid"`
	XCPCLevel  int         `json:"xcpcLevel"`
}

// TrainingProblems 训练题目列表
type TrainingProblems struct {
	Count   int         `json:"count"`
	PerPage interface{} `json:"perPage"`
	Result  [][]string  `json:"result"`
}
