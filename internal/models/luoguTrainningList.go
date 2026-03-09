package models

//===============================洛谷题单列表响应===============================
type LuoguTrainningListResponse struct {
	Code            int                           `json:"code"`
	CurrentData     LuoguTrainningListCurrentData `json:"currentData"`
	CurrentTemplate string                        `json:"currentTemplate"`
	CurrentTheme    interface{}                   `json:"currentTheme"` // null
	CurrentTime     int64                         `json:"currentTime"`
	CurrentTitle    string                        `json:"currentTitle"`
}

// LuoguCurrentData 当前数据
type LuoguTrainningListCurrentData struct {
	AcceptedCounts interface{}    `json:"acceptedCounts"` // null
	Trainings      LuoguTrainings `json:"trainings"`
}

// LuoguTrainings 训练列表
type LuoguTrainings struct {
	Count   int              `json:"count"`
	PerPage int              `json:"perPage"`
	Result  []LuoguTrainning `json:"result"`
}

// LuoguPractice 单个练习
type LuoguTrainning struct {
	CreateTime   int64         `json:"createTime"`
	Deadline     interface{}   `json:"deadline"` // null
	ID           int64         `json:"id"`
	MarkCount    int           `json:"markCount"`
	Marked       bool          `json:"marked"`
	Name         string        `json:"name"`
	ProblemCount int           `json:"problemCount"`
	Provider     LuoguProvider `json:"provider"`
	Title        string        `json:"title"`
	Type         int           `json:"type"`
}

// LuoguProvider 提供者信息
type LuoguProvider struct {
	Avatar     string  `json:"avatar"`
	Background string  `json:"background"`
	Badge      *string `json:"badge"` // null | string
	CCFLevel   int     `json:"ccfLevel"`
	Color      string  `json:"color"`
	IsAdmin    bool    `json:"isAdmin"`
	IsBanned   bool    `json:"isBanned"`
	Name       string  `json:"name"`
	Slogan     string  `json:"slogan"`
	UID        int64   `json:"uid"`
	XCPCLevel  int     `json:"xcpcLevel"`
}
