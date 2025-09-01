package component

import (
	"fmt"
	"time"
)

type CountTime struct {
	StartTime int64
	EndTime   int64
	Duration  int64 `ms`
}

func (c *CountTime) Start() {
	c.StartTime = time.Now().UnixMilli()
}
func (c *CountTime) End() {
	c.EndTime = time.Now().UnixMilli()
	c.Duration = c.EndTime - c.StartTime
	fmt.Printf("总耗时:%ds/%dms   \n", c.Duration/1000, c.Duration)
}
func DateTime(rawStamp int64) string {
	stamp := time.Unix(rawStamp, 1)
	str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", stamp.Year(), stamp.Month(), stamp.Day(), stamp.Hour(), stamp.Minute(), stamp.Second())
	return str
}
func NowDateTime() string {
	now := time.Now()
	str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	return str
}
func NowDate() string {
	now := time.Now()
	str := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	return str
}
func BeforDateTime(day int64) string {
	now := time.Now()
	befor := time.Unix(now.Unix()-day*24*3600, 1)
	str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", befor.Year(), befor.Month(), befor.Day(), befor.Hour(), befor.Minute(), befor.Second())
	return str
}
