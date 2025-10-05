package utils

import "time"

type Timer struct{}

func NewTimer() *Timer {
	return &Timer{}
}
func (t *Timer) CountDurationInMs(f func()) int64 {
	start := time.Now()
	f()
	end := time.Since(start).Milliseconds()
	return end
}
func (t *Timer) CountDurationInS(f func()) int64 {
	start := time.Now()
	f()
	end := int64(time.Since(start).Seconds())
	return end
}
func (t *Timer) CountDurationStr(f func()) string {
	start := time.Now()
	f()
	end := time.Since(start).String()
	return end
}
