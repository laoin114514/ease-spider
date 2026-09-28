package models

// CfUserData 是 CF 提交记录任务中单个用户的处理上下文。
type CfUserData struct {
	Account     string
	RealName    string
	OldDataSet  map[int]bool
	InsertCount int
}
