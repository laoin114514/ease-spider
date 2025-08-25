package main

import (
	"fmt"
	"spider/cf/cfBetterSub"
	"spider/component"
	"time"
)

func main() {
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()
	// luogu.Use()
	fmt.Println("开始执行")
	var countTime component.CountTime
	countTime.Start()
	cfBetterSub.Use()
	countTime.End()
	for range ticker.C {
		countTime.Start()
		cfBetterSub.Use()
		countTime.End()
		fmt.Println("")
	}
}
