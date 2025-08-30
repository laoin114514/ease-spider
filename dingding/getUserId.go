package dingding

import (
	"encoding/json"
	"fmt"
	"os"
)

func getUserId() map[string]any {
	content, err := os.ReadFile("dingUser_id.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var result map[string]any
	json.Unmarshal(content, &result)
	return result
}
