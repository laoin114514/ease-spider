package getcookie

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"spider/config"
	"spider/internal/utils"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func (l *GetCookie) Update() error {
	l.init()
	l.initRedirect()
	l.getCaptcha()
	l.redirCaptcha()
	captcha, err := l.identify(config.AppConfig.Luogu.IsInServer)
	if err != nil {
		l.log.Errorf("验证码识别失败 %s", err.Error())
		return err
	}
	l.log.Printf("验证码：%s", captcha)
	err = l.login(captcha)
	if err != nil {
		l.log.Errorf("登录失败 %s", err.Error())
		return err
	}
	l.log.Println("登录成功")
	return nil
}
func (l *GetCookie) init() {
	c := l.restyInit()
	resp, _ := c.R().Get("https://www.luogu.com.cn/auth/login")
	cookie := resp.Cookies()
	l.cookiePool["cookie2"] = cookie[0].Name + "=" + cookie[0].Value
}

func (l *GetCookie) initRedirect() {
	c := l.restyInit()
	resp, _ := c.R().
		SetHeader("Cookie", l.cookiePool["cookie2"]).
		Get("https://www.luogu.com.cn/auth/login")
	cookie := resp.Cookies()
	l.cookiePool["cookie1"] = cookie[0].Name + "=" + cookie[0].Value
}

func (l *GetCookie) restyInit() *resty.Client {
	c := resty.New()
	c.SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetHeader("User-Agent", config.AppConfig.Luogu.UserAgent)
	return c
}
func (l *GetCookie) getCaptcha() {
	c := l.restyInit()
	now := time.Now()
	stamp := float64(now.UnixMicro()) / 1000
	resp, _ := c.R().
		Get(fmt.Sprintf("https://www.luogu.com.cn/lg4/captcha?_t=%f", stamp))
	newCookie := resp.Cookies()
	l.cookiePool["cookie2"] = newCookie[0].Name + "=" + newCookie[0].Value
}
func (l *GetCookie) redirCaptcha() {
	c := l.restyInit()
	now := time.Now()
	stamp := float64(now.UnixMicro()) / 1000
	cookie := l.cookiePool["cookie1"] + "; " + l.cookiePool["cookie2"]
	resp, _ := c.R().
		SetHeader("Cookie", cookie).
		Get(fmt.Sprintf("https://www.luogu.com.cn/lg4/captcha?_t=%f", stamp))
	newCookie := resp.Cookies()
	l.cookiePool["cookie2"] = newCookie[0].Name + "=" + newCookie[0].Value
	l.saveImage(resp.Body())
}
func (l *GetCookie) saveImage(content []byte) {
	file, err := os.OpenFile("data/captcha.jpg", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		l.log.Println(err)
		return
	}
	defer file.Close()
	file.Write(content)
}
func (l *GetCookie) login(captcha string) error {
	c := l.restyInit()
	cookie := l.cookiePool["cookie1"] + "; " + l.cookiePool["cookie2"]
	resp, err := c.R().
		SetBody(map[string]any{
			"username": config.AppConfig.Luogu.Username,
			"password": config.AppConfig.Luogu.Password,
			"captcha":  captcha,
		}).
		SetHeader("Cookie", cookie).
		Post("https://www.luogu.com.cn/do-auth/password")
	if resp.StatusCode() != 200 {
		return errors.New("验证码错误" + strconv.Itoa(resp.StatusCode()))
	}
	if err != nil {
		return err
	}
	arr := resp.Cookies()
	uid := arr[0].Name + "=" + arr[0].Value
	Cookie := cookie + ";" + uid
	utils.JsonDB.Set("Cookie", Cookie)
	return nil
}

type ocrServerOut struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Text string `json:"text"`
	} `json:"data"`
}

func (l *GetCookie) identify(isInServer bool) (string, error) {
	c := resty.New()
	var result ocrServerOut
	resp, err := c.R().
		SetFile("file", "data/captcha.jpg").
		Post("http://127.0.0.1:8000/ocr/classify")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", errors.New("验证码错误" + strconv.Itoa(resp.StatusCode()))
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}
	return result.Data.Text, nil
}
