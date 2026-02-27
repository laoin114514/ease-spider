package handler

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// 调试脚本（测试用）
func Test() {
	url := "https://www.luogu.com.cn/problem/list?type=all&keyword=717D&page=1"
	var result map[string]any
	c := resty.New()
	c.R().
		SetResult(&result).
		SetHeader("Cookie", "__client_id=bea22041f2bcb2a5690b508231c59c1dd75faf59; C3VK=eaef32;_uid=1851093").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36 Edg/143.0.0.0").
		Get(url)
	fmt.Println(result)
}
