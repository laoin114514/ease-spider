package cfBetterSub

import (
	"errors"
	"fmt"
	genUrl2 "spider/src/cf/genUrl"

	"github.com/go-resty/resty/v2"
)

func request(parm chanParm) ([]any, error) {
	if parm.handle == "" {
		return []any{}, errors.New(fmt.Sprintf("%s的cf账号不存在", parm.username))
	}
	c := resty.New()
	var user genUrl2.User
	url, genErr := user.Status(true, genUrl2.User_status{
		Handle:         parm.handle,
		From:           1,
		Count:          50000,
		IncludeSources: false,
	})
	if genErr != nil {
		return nil, genErr
	}

	var response map[string]any
	resp, err := c.R().
		SetResult(&response).
		Get(url)
	code := resp.StatusCode()
	if code == 429 {
		errCount++
		return nil, errors.New(fmt.Sprintf("%s请求错误 %d", parm.username, code))
	}
	if err != nil {
		fmt.Println("请求错误", err)
		return nil, err
	}

	results, ok := response["result"].([]any)
	if !ok {
		return []any{}, errors.New(fmt.Sprintf("%s的cf账号不存在", parm.username))
	}
	return results, nil
}
