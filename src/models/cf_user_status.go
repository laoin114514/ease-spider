package models

// Codeforces用户基本信息
type CfUserData struct {
	Account     string
	RealName    string
	OldDataSet  map[int]bool
	InsertCount int
}

// Codeforces用户状态响应结构体
type CfUserStatusResponse struct {
	CfResponse[[]CfSubmission]
}

// Codeforces提交记录结构体
type CfSubmission struct {
	Id                  int64     `json:"id"`                  // 提交ID
	ContestId           int64     `json:"contestId"`           // 比赛ID
	CreationTimeSeconds int64     `json:"creationTimeSeconds"` // 创建时间（秒）
	RelativeTimeSeconds int64     `json:"relativeTimeSeconds"` // 相对时间（秒）
	Problem             CfProblem `json:"problem"`             // 题目信息
	Author              CfAuthor  `json:"author"`              // 作者信息
	ProgrammingLanguage string    `json:"programmingLanguage"` // 编程语言
	Verdict             string    `json:"verdict"`             // 判题结果
	Testset             string    `json:"testset"`             // 测试集类型
	PassedTestCount     int       `json:"passedTestCount"`     // 通过的测试用例数
	TimeConsumedMillis  int64     `json:"timeConsumedMillis"`  // 时间消耗（毫秒）
	MemoryConsumedBytes int64     `json:"memoryConsumedBytes"` // 内存消耗（字节）
}

// Codeforces作者信息结构体
type CfAuthor struct {
	ContestId        int64      `json:"contestId"`        // 比赛ID
	ParticipantId    int64      `json:"participantId"`    // 参与者ID
	Members          []CfMember `json:"members"`          // 团队成员
	ParticipantType  string     `json:"participantType"`  // 参与者类型：VIRTUAL、CONTESTANT等
	TeamId           int64      `json:"teamId"`           // 团队ID
	TeamName         string     `json:"teamName"`         // 团队名称
	Ghost            bool       `json:"ghost"`            // 是否为幽灵用户
	StartTimeSeconds int64      `json:"startTimeSeconds"` // 开始时间（秒）
}

// Codeforces团队成员结构体
type CfMember struct {
	Handle string `json:"handle"` // 用户句柄
}

// Codeforces提交记录统计信息
type CfSubmissionStats struct {
	TotalSubmissions    int                       `json:"totalSubmissions"`    // 总提交数
	AcceptedSubmissions int                       `json:"acceptedSubmissions"` // 通过提交数
	UserStats           map[string]CfUserStats    `json:"userStats"`           // 用户统计
	ProblemStats        map[string]CfProblemStats `json:"problemStats"`        // 题目统计
	LanguageStats       map[string]int            `json:"languageStats"`       // 语言统计
	VerdictStats        map[string]int            `json:"verdictStats"`        // 判题结果统计
}

// Codeforces用户统计信息
type CfUserStats struct {
	Handle              string `json:"handle"`              // 用户句柄
	TotalSubmissions    int    `json:"totalSubmissions"`    // 总提交数
	AcceptedSubmissions int    `json:"acceptedSubmissions"` // 通过提交数
	LastSubmissionTime  int64  `json:"lastSubmissionTime"`  // 最后提交时间
}

// Codeforces题目统计信息
type CfProblemStats struct {
	ProblemIndex        string `json:"problemIndex"`        // 题目索引
	ProblemName         string `json:"problemName"`         // 题目名称
	TotalSubmissions    int    `json:"totalSubmissions"`    // 总提交数
	AcceptedSubmissions int    `json:"acceptedSubmissions"` // 通过提交数
	UniqueUsers         int    `json:"uniqueUsers"`         // 唯一用户数
}
