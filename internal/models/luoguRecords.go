package models

// LuoguRecordsResponse 洛谷记录列表响应
// LuoguRecordsResponse 洛谷记录列表响应
type LuoguRecordsResponse struct {
	Code            int              `json:"code"`
	CurrentTemplate string           `json:"currentTemplate"`
	CurrentData     LuoguRecordsData `json:"currentData"` //主要数据
	CurrentTitle    string           `json:"currentTitle"`
	CurrentTheme    LuoguTheme       `json:"currentTheme"`
	CurrentTime     int64            `json:"currentTime"`
	CurrentUser     LuoguUser        `json:"currentUser"`
}

// LuoguRecordsData 记录列表数据
type LuoguRecordsData struct {
	Records LuoguRecords `json:"records"`
}

// LuoguRecords 记录列表
type LuoguRecords struct {
	Result  []LuoguRecord `json:"result"`
	PerPage int           `json:"perPage"`
	Count   int           `json:"count"`
}

// LuoguRecord 单条提交记录
type LuoguRecord struct {
	Time             int           `json:"time"`
	Memory           int           `json:"memory"`
	Problem          LuoguProblem  `json:"problem"`
	Contest          *LuoguContest `json:"contest"`
	SourceCodeLength int           `json:"sourceCodeLength"`
	SubmitTime       int64         `json:"submitTime"`
	Language         int           `json:"language"`
	User             LuoguUser     `json:"user"`
	ID               int64         `json:"id"`
	Status           int           `json:"status"`
	EnableO2         bool          `json:"enableO2"`
	Score            int           `json:"score"`
}

// LuoguProblem 题目信息
type LuoguProblem struct {
	PID        string `json:"pid"`
	Title      string `json:"title"`
	Difficulty int    `json:"difficulty"`
	FullScore  int    `json:"fullScore"`
	Type       string `json:"type"`
}

// LuoguContest 比赛信息
type LuoguContest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	StartTime   int64  `json:"startTime"`
	EndTime     int64  `json:"endTime"`
	Type        int    `json:"type"`
	Rated       bool   `json:"rated"`
	Description string `json:"description"`
}

//============================================以下均为无关紧要的数据============================================//
// LuoguUser 用户信息
type LuoguUser struct {
	FollowingCount     int     `json:"followingCount"`
	FollowerCount      int     `json:"followerCount"`
	Ranking            int     `json:"ranking"`
	EloValue           *int    `json:"eloValue"`
	BlogAddress        *string `json:"blogAddress"`
	UnreadMessageCount int     `json:"unreadMessageCount"`
	UnreadNoticeCount  int     `json:"unreadNoticeCount"`
	UID                int64   `json:"uid"`
	Name               string  `json:"name"`
	Avatar             string  `json:"avatar"`
	Slogan             string  `json:"slogan"`
	Badge              *string `json:"badge"`
	IsAdmin            bool    `json:"isAdmin"`
	IsBanned           bool    `json:"isBanned"`
	Color              string  `json:"color"`
	CCFLevel           int     `json:"ccfLevel"`
	XCPCLevel          int     `json:"xcpcLevel"`
	Background         string  `json:"background"`
	Verified           bool    `json:"verified"`
}

// LuoguTheme 主题配置
type LuoguTheme struct {
	ID      int          `json:"id"`
	Header  LuoguHeader  `json:"header"`
	SideNav LuoguSideNav `json:"sideNav"`
	Footer  LuoguFooter  `json:"footer"`
}

// LuoguHeader 头部配置
type LuoguHeader struct {
	ImagePath  *string `json:"imagePath"`
	Color      [][]int `json:"color"`
	Blur       int     `json:"blur"`
	Brightness int     `json:"brightness"`
	Degree     int     `json:"degree"`
	Repeat     int     `json:"repeat"`
	Position   []int   `json:"position"`
	Size       []int   `json:"size"`
	Type       int     `json:"type"`
	ClassName  string  `json:"__CLASS_NAME"`
}

// LuoguSideNav 侧边导航配置
type LuoguSideNav struct {
	LogoBackgroundColor []int  `json:"logoBackgroundColor"`
	Color               []int  `json:"color"`
	InvertColor         bool   `json:"invertColor"`
	ClassName           string `json:"__CLASS_NAME"`
}

// LuoguFooter 底部配置
type LuoguFooter struct {
	ImagePath  *string `json:"imagePath"`
	Color      [][]int `json:"color"`
	Blur       int     `json:"blur"`
	Brightness int     `json:"brightness"`
	Degree     int     `json:"degree"`
	Repeat     int     `json:"repeat"`
	Position   []int   `json:"position"`
	Size       []int   `json:"size"`
	Type       int     `json:"type"`
	ClassName  string  `json:"__CLASS_NAME"`
}
