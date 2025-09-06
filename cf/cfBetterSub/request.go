package cfBetterSub

import (
	"errors"
	"fmt"
	"spider/cf/genUrl"

	"github.com/go-resty/resty/v2"
)

func request(handle string) ([]any, error) {
	c := resty.New()
	var user genUrl.User
	url, genErr := user.Status(true, genUrl.User_status{
		Handle:         handle,
		From:           1,
		Count:          50000,
		IncludeSources: false,
	})
	if genErr != nil {
		return nil, genErr
	}

	var response map[string]any
	_, err := c.R().
		SetResult(&response).
		Get(url)
	if err != nil {
		fmt.Println("请求错误", err)
		return nil, err
	}

	results, ok := response["result"].([]any)
	if !ok {
		return []any{}, errors.New("handle不正确")
	}
	return results, nil
}
