package luogu

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func url(uid string) string {
	return fmt.Sprintf("https://www.luogu.com.cn/user/%s?_contentOnly=1", uid)
}

func request(uid string) []any {
	godotenv.Load()
	cookie := os.Getenv("Cookie")
	c := resty.New()
	var result map[string]any
	c.R().
		SetResult(&result).
		SetHeader("Cookie", cookie).
		Get(url(uid))
	code := result["code"].(float64)
	if code != 200 {
		fmt.Printf("%s请求错误 code:%v\n", uid, code)
		return []any{}
	}
	passedProblems, ok := result["currentData"].(map[string]any)["passedProblems"].([]any)
	if !ok {
		passedProblems = []any{}
		fmt.Printf("%s无权限访问\n", uid)
	}
	return passedProblems
}
