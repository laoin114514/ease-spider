package main

import (
	"fmt"
	"spider/config/db"
	"spider/src/models"
	"spider/src/utils"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db.Init()
	utils.InitGlobalJSONDB("tempDB.json")
	structFunc := utils.NewStructFunc(models.ContestStandingsParams{
		Handles:   "tourist",
		ContestID: 1000,
		AsManager: true,
		From:      1,
		Count:     10,
	})
	str, err := structFunc.StructToOrderParams()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
}
