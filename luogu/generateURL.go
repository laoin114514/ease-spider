package luogu

import "fmt"

func url(uid string) string {
	return fmt.Sprintf("https://www.luogu.com.cn/user/%s?_contentOnly=1", uid)
}
