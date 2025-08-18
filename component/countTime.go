package component

import (
	"fmt"
	"time"
)

type CountTime struct {
	StartTime int64
	EndTime   int64
}

func (c *CountTime) Start() {
	c.StartTime = time.Now().Unix()
}
func (c *CountTime) End() {
	c.EndTime = time.Now().Unix()
	fmt.Printf("总耗时:%ds", c.EndTime-c.StartTime)
}
