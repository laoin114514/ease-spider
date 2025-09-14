package dingding

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func getToken() string {
	godotenv.Load()
	appKey := os.Getenv("ding_accessToken")
	appSecret := os.Getenv("ding_appSecret")

	url := "https://api.dingtalk.com/v1.0/oauth2/accessToken"
	c := resty.New()
	var result map[string]any
	_, err := c.R().
		SetBody(map[string]any{
			"appKey":    appKey,
			"appSecret": appSecret,
		}).
		SetResult(&result).
		Post(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return result["accessToken"].(string)
}
