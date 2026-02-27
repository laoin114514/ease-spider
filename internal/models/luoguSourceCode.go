package models

// LuoguSubmissionDetailResponse 洛谷提交记录详情响应
type LuoguSubmissionDetailResponse struct {
	Code            int                 `json:"code"`
	CurrentData     LuoguSubmissionData `json:"currentData"`
	CurrentTemplate string              `json:"currentTemplate"`
	CurrentTheme    LuoguTheme          `json:"currentTheme"`
	CurrentTime     int64               `json:"currentTime"`
	CurrentTitle    string              `json:"currentTitle"`
	CurrentUser     LuoguUser           `json:"currentUser"`
}

// LuoguSubmissionData 提交记录数据
type LuoguSubmissionData struct {
	Record        LuoguSubmissionRecord `json:"record"`
	ShowStatus    bool                  `json:"showStatus"`
	TestCaseGroup [][]int               `json:"testCaseGroup"`
}

// LuoguSubmissionRecord 提交记录详情
type LuoguSubmissionRecord struct {
	Contest          *LuoguContest     `json:"contest"`
	Detail           LuoguRecordDetail `json:"detail"`
	EnableO2         bool              `json:"enableO2"`
	ID               int64             `json:"id"`
	Language         int               `json:"language"`
	Memory           int               `json:"memory"`
	Problem          LuoguProblem      `json:"problem"`
	Score            int               `json:"score"`
	SourceCode       string            `json:"sourceCode"`
	SourceCodeLength int               `json:"sourceCodeLength"`
	Status           int               `json:"status"`
	SubmitTime       int64             `json:"submitTime"`
	Time             int               `json:"time"`
	User             LuoguUser         `json:"user"`
}

// LuoguRecordDetail 记录详情
type LuoguRecordDetail struct {
	ClassName     string             `json:"__CLASS_NAME"`
	CompileResult LuoguCompileResult `json:"compileResult"`
	JudgeResult   LuoguJudgeResult   `json:"judgeResult"`
	Version       int                `json:"version"`
}

// LuoguCompileResult 编译结果
type LuoguCompileResult struct {
	ClassName string  `json:"__CLASS_NAME"`
	Message   *string `json:"message"`
	Opt2      bool    `json:"opt2"`
	Success   bool    `json:"success"`
}

// LuoguJudgeResult 评测结果
type LuoguJudgeResult struct {
	ClassName         string                    `json:"__CLASS_NAME"`
	FinishedCaseCount int                       `json:"finishedCaseCount"`
	Memory            int                       `json:"memory"`
	Score             int                       `json:"score"`
	Status            int                       `json:"status"`
	Subtasks          []LuoguSubtaskJudgeResult `json:"subtasks"`
	Time              int                       `json:"time"`
}

// LuoguSubtaskJudgeResult 子任务评测结果
type LuoguSubtaskJudgeResult struct {
	ClassName string                         `json:"__CLASS_NAME"`
	ID        int                            `json:"id"`
	Judger    string                         `json:"judger"`
	Memory    int                            `json:"memory"`
	Score     int                            `json:"score"`
	Status    int                            `json:"status"`
	TestCases map[string]LuoguTestCaseResult `json:"testCases"`
	Time      int                            `json:"time"`
}

// LuoguTestCaseResult 测试用例结果
type LuoguTestCaseResult struct {
	ClassName   string `json:"__CLASS_NAME"`
	Description string `json:"description"`
	ExitCode    int    `json:"exitCode"`
	ID          int    `json:"id"`
	Memory      int    `json:"memory"`
	Score       int    `json:"score"`
	Signal      int    `json:"signal"`
	Status      int    `json:"status"`
	SubtaskID   int    `json:"subtaskID"`
	Time        int    `json:"time"`
}

// LuoguHeaderFooterConfig 头部/底部配置
type LuoguHeaderFooterConfig struct {
	ClassName  string  `json:"__CLASS_NAME"`
	Blur       int     `json:"blur"`
	Brightness int     `json:"brightness"`
	Color      [][]int `json:"color"`
	Degree     int     `json:"degree"`
	ImagePath  *string `json:"imagePath"`
	Position   []int   `json:"position"`
	Repeat     int     `json:"repeat"`
	Size       []int   `json:"size"`
	Type       int     `json:"type"`
}

// LuoguSideNavConfig 侧边导航配置
type LuoguSideNavConfig struct {
	ClassName           string `json:"__CLASS_NAME"`
	Color               []int  `json:"color"`
	InvertColor         bool   `json:"invertColor"`
	LogoBackgroundColor []int  `json:"logoBackgroundColor"`
}
