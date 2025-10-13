package handler

import (
	"fmt"
	"spider/src/services"
)

func Test() {
	luoguSolutionService := services.NewLuoguSolution()
	err := luoguSolutionService.GetSolutionHasSourceCode()
	if err != nil {
		fmt.Println(err)
	}
	luoguSolutionService.LogService.Clear()
}
