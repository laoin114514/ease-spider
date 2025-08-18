package main

import (
	"fmt"
	"spider/cf/cfBetterSub"
	"time"
)

func main() {
	start := time.Now().Unix()
	cfBetterSub.Use()
	end := time.Now().Unix()
	fmt.Printf("耗时:%ds", (end - start))
}
