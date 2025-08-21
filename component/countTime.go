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
	fmt.Printf("总耗时:%ds/%dms   ", c.Duration/1000, c.Duration)
}
