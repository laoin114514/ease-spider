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
	structFunc := utils.NewStructFunc(models.User_info{
		Handles:              []string{"tourist", "test123"},
		CheckHistoricHandles: true,
	})
	str, err := structFunc.StructToOrderParams()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
}
