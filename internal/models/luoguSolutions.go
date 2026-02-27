package models

type LuoguSolutionResponse struct {
	Instance string `json:"instance"`
	Template string `json:"template"`
	Status   int    `json:"status"`
	Locale   string `json:"locale"`
	Data     struct {
		Solutions LuoguSolution `json:"solutions"`
		Problem   struct {
			PID        string `json:"pid"`
			Type       string `json:"type"`
			Title      string `json:"title"`
			Difficulty int    `json:"difficulty"`
			Submitted  bool   `json:"submitted"`
			Accepted   bool   `json:"accepted"`
		} `json:"problem"`
		AcceptSolution bool `json:"acceptSolution"`
	} `json:"data"`
	User struct {
		UID                int    `json:"uid"`
		Avatar             string `json:"avatar"`
		Name               string `json:"name"`
		Slogan             string `json:"slogan"`
		Badge              string `json:"badge"`
		IsAdmin            bool   `json:"isAdmin"`
		IsBanned           bool   `json:"isBanned"`
		Color              string `json:"color"`
		CCFLevel           int    `json:"ccfLevel"`
		XCPCLevel          int    `json:"xcpcLevel"`
		Background         string `json:"background"`
		Verified           bool   `json:"verified"`
		EloValue           *int   `json:"eloValue"`
		FollowingCount     int    `json:"followingCount"`
		FollowerCount      int    `json:"followerCount"`
		Ranking            *int   `json:"ranking"`
		UnreadMessageCount int    `json:"unreadMessageCount"`
		UnreadNoticeCount  int    `json:"unreadNoticeCount"`
	} `json:"user"`
	Time  float64 `json:"time"`
	Theme struct {
		ID     int `json:"id"`
		Header struct {
			Type       int     `json:"type"`
			Color      [][]int `json:"color"`
			Degree     int     `json:"degree"`
			Blur       int     `json:"blur"`
			Brightness int     `json:"brightness"`
			Position   []int   `json:"position"`
			Repeat     int     `json:"repeat"`
			ClassName  string  `json:"__CLASS_NAME"`
		} `json:"header"`
		SideNav struct {
			LogoBackgroundColor []int  `json:"logoBackgroundColor"`
			Color               []int  `json:"color"`
			InvertColor         bool   `json:"invertColor"`
			ClassName           string `json:"__CLASS_NAME"`
		} `json:"sideNav"`
		Footer struct {
			Type       int     `json:"type"`
			Color      [][]int `json:"color"`
			Degree     int     `json:"degree"`
			Blur       int     `json:"blur"`
			Brightness int     `json:"brightness"`
			ClassName  string  `json:"__CLASS_NAME"`
		} `json:"footer"`
	} `json:"theme"`
}
type LuoguSolution struct {
	PerPage int               `json:"perPage"`
	Count   int               `json:"count"`
	Result  []SolutionContent `json:"result"` //md格式的题解
}

// Solution 题解结构
type SolutionContent struct {
	LID           string          `json:"lid"`
	Title         string          `json:"title"`
	Time          int64           `json:"time"`
	Author        SolutionAuthor  `json:"author"`
	Upvote        int             `json:"upvote"`
	ReplyCount    int             `json:"replyCount"`
	FavorCount    int             `json:"favorCount"`
	Category      int             `json:"category"`
	Status        int             `json:"status"`
	SolutionFor   SolutionProblem `json:"solutionFor"`
	PromoteStatus int             `json:"promoteStatus"`
	Collection    map[string]any  `json:"collection"`
	Content       string          `json:"content"`
	CategoryOld   string          `json:"categoryOld"`
	ContentFull   bool            `json:"contentFull"`
	AdminNote     *string         `json:"adminNote"`
	Voted         *int            `json:"voted"`
	CanReply      bool            `json:"canReply"`
	CanEdit       bool            `json:"canEdit"`
}

// Author 作者信息结构
type SolutionAuthor struct {
	UID        int     `json:"uid"`
	Avatar     string  `json:"avatar"`
	Name       string  `json:"name"`
	Slogan     string  `json:"slogan"`
	Badge      *string `json:"badge"`
	IsAdmin    bool    `json:"isAdmin"`
	IsBanned   bool    `json:"isBanned"`
	Color      string  `json:"color"`
	CCFLevel   int     `json:"ccfLevel"`
	XCPCLevel  int     `json:"xcpcLevel"`
	Background string  `json:"background"`
	IsRoot     *bool   `json:"isRoot,omitempty"`
}

// Problem 题目信息结构
type SolutionProblem struct {
	PID        string `json:"pid"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Difficulty int    `json:"difficulty"`
	Submitted  bool   `json:"submitted"`
	Accepted   bool   `json:"accepted"`
}
