package dingding

import "time"

type dingCheckUp struct {
	name      string
	userId    string
	time      time.Time
	checkType string
}
