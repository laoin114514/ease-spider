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

// Codeforces比赛信息结构体
type CfContest struct {
	Id                  int64  `json:"id"`                  //比赛id
	Name                string `json:"name"`                //比赛名称
	Type                string `json:"type"`                //比赛类型
	StartTime           int64  `json:"startTime"`           //开始时间
	PreparedBy          string `json:"preparedBy"`          //准备人
	Phase               string `json:"phase"`               //阶段
	Frozen              bool   `json:"frozen"`              //是否冻结
	DurationSeconds     int64  `json:"durationSeconds"`     //持续时间秒
	StartTimeSeconds    int64  `json:"startTimeSeconds"`    //开始时间秒
	RelativeTimeSeconds int64  `json:"relativeTimeSeconds"` //相对时间秒
	Country             string `json:"country"`             //国家
	City                string `json:"city"`                //城市
	Kind                string `json:"kind"`                //种类
	Season              string `json:"season"`              //赛季
	Difficulty          int    `json:"difficulty"`          //难度
}
