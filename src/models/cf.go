package models

type CfResponse[T any] struct {
	Status string `json:"status"` // 响应状态：OK
	Result T      `json:"result"` // 响应结果
}

// Codeforces题目信息结构体
type CfProblem struct {
	ContestId int64    `json:"contestId"` // 比赛ID
	Index     string   `json:"index"`     // 题目索引（如A、B、C）
	Name      string   `json:"name"`      // 题目名称
	Type      string   `json:"type"`      // 题目类型：PROGRAMMING
	Tags      []string `json:"tags"`      // 题目标签
	Points    float64  `json:"points"`    // 题目分数
	Rating    int64    `json:"rating"`    // 题目难度
}

// Codeforces请求参数结构体
