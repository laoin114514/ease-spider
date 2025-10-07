package main

import (
	"spider/config/db"
	"spider/src/handler"
	"spider/src/utils"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	db.Init()
	godotenv.Load()
	utils.InitGlobalJSONDB("tempDB.json")
	utils.InitGenerateCFurl()

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	handler.TimeTask()
}
