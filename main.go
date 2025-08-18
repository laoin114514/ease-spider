package main

import (
	"spider/cf/cfBetterSub"
	"spider/component"
)

func main() {
	var countTime component.CountTime
	countTime.Start()
	cfBetterSub.Use()
	countTime.End()
}
