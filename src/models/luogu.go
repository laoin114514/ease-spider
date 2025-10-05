package models

type LuoguUserDeliver struct {
	Uid        string
	RealName   string
	OldDataSet map[string]bool
	Count      int
}
