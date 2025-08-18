package genUrl

import (
	"errors"
	"fmt"
	"spider/component"
	"time"
)

const randomKey = "123456"
const Baseurl = "https://codeforces.com/api/"

func apiKeyUrl(method string, pararms map[string]any) (string, error) {
	ok, apiKey := findApikey(pararms["handle"].(string))
	if !ok {
		return "", errors.New(fmt.Sprintf("%s的apiKey不存在", pararms["handle"].(string)))
	}
	pararms["apiKey"] = apiKey.ApiKey
	pararms["time"] = time.Now().Unix()
	tail := buildPararms(pararms)
	hashCode := component.Hash(fmt.Sprintf("%v/%v?%v#%v", randomKey, method, tail, apiKey.Secret))
	return fmt.Sprintf("https://codeforces.com/api/%v?%v&apiSig=%v%v", method, tail, randomKey, hashCode), nil
}
func defaultUrl(method string, pararms map[string]any) (string, error) {
	tail := buildPararms(pararms)
	return fmt.Sprintf("https://codeforces.com/api/%v?%v", method, tail), nil
}
