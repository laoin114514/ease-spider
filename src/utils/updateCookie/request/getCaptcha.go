package request

import (
	"fmt"
	"os"
	"time"
)

func GetCaptcha() {
	c := RestyInit()
	now := time.Now()
	stamp := float64(now.UnixMicro()) / 1000
	// cookie := tempDB.Get("cookie1").(string) + "; " + tempDB.Get("cookie2").(string)
	resp, _ := c.R().
		// SetHeader("Cookie", cookie).
		Get(fmt.Sprintf("https://www.luogu.com.cn/lg4/captcha?_t=%f", stamp))
	newCookie := resp.Cookies()
	cookie2 := newCookie[0].Name + "=" + newCookie[0].Value
	JsonDB.Set("cookie2", cookie2)
}
func RedirCaptcha() {
	c := RestyInit()
	now := time.Now()
	stamp := float64(now.UnixMicro()) / 1000
	cookie := JsonDB.Get("cookie1").(string) + "; " + JsonDB.Get("cookie2").(string)
	resp, _ := c.R().
		SetHeader("Cookie", cookie).
		Get(fmt.Sprintf("https://www.luogu.com.cn/lg4/captcha?_t=%f", stamp))
	newCookie := resp.Cookies()
	cookie2 := newCookie[0].Name + "=" + newCookie[0].Value
	JsonDB.Set("cookie2", cookie2)
	saveImage(resp.Body())
}
func saveImage(content []byte) {
	file, err := os.OpenFile(fmt.Sprintf("ocr/captcha.jpg"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write(content)
}
