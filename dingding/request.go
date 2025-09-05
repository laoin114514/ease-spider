package dingding

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func request(token string, from string, to string) []any {
	userObject := tempDB.Get("dingUserId").(map[string]any)
	userId := []string{}
	for k, _ := range userObject {
		userId = append(userId, k)
	}
	body := map[string]any{
		"userIds":       userId,
		"checkDateFrom": from,
		"checkDateTo":   to,
	}
	c := resty.New()
	var result map[string]any
	response, err := c.R().
		SetBody(body).
		SetResult(&result).
		Post("https://oapi.dingtalk.com/attendance/listRecord?access_token=" + token)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if response.StatusCode() != 200 {
		fmt.Println("请求失败，状态码", response.StatusCode())
	}
	return result["recordresult"].([]any)
}
