package main

import (
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/src/handler"
	"spider/src/utils"
	"sync"
)

func main() {
	err := config.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.InitGlobalJSONDB("config.json")
	utils.InitGenerateCFurl()

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	handler.TimeTask()
}
