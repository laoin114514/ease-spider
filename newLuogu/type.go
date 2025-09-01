package newLuogu

import "time"

type luogu_all_submissions struct {
	subId       string
	pid         string
	username    string
	uid         string
	isPass      bool
	subTime     time.Time
	problemName string
	difficulty  string
}
