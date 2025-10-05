package models

type CfUserData struct {
	Account  string
	RealName string
}
type CfRecordData struct {
	Status string `json:"status"`
	Result []any  `json:"result"`
}
