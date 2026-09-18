package models

type LuoguUserDeliver struct {
	Uid        string
	RealName   string
	OldDataSet map[string]bool
	Count      int
	// Pending 本轮因"还在测评"被跳过的记录数：这些记录还没入库，但也不能算作缺失，
	// 否则每次有人刚提交完都会误判成"数据不一致"而触发一次全量爬取
	Pending int
}
