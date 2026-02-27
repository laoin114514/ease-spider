package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"spider/config"
	"spider/internal/constants"
	"spider/internal/utils"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// ocr缺失会导致登录失败，此时可以暂停使用该组件
type LuoguUpdateCookie struct {
	*LogService
	cookiePool map[string]string
	debug      *utils.Debug
}

func (l *LuoguUpdateCookie) Update() error {
	l.init()
	l.initRedirect()
	l.getCaptcha()
	l.redirCaptcha()
	captcha, err := l.identify(config.AppConfig.Luogu.IsInServer)
	if err != nil {
		l.debug.Debug(fmt.Sprintf("验证码识别失败 %s", err.Error()))
		l.AddLog(fmt.Sprintf("验证码识别失败 %s", err.Error()))
		return err
	}
	l.AddLog(fmt.Sprintf("验证码：%s", captcha))
	err = l.login(captcha)
	if err != nil {
		l.debug.Debug(fmt.Sprintf("登录失败 %s", err.Error()))
		l.AddLog(fmt.Sprintf("登录失败 %s", err.Error()))
		return err
	}
	l.debug.Debug("登录成功")
	l.AddLog("登录成功")
	return nil
}
func (l *LuoguUpdateCookie) init() {
	c := l.restyInit()
	resp, _ := c.R().Get("https://www.luogu.com.cn/auth/login")
	cookie := resp.Cookies()
	l.cookiePool["cookie2"] = cookie[0].Name + "=" + cookie[0].Value
}

func (l *LuoguUpdateCookie) initRedirect() {
	c := l.restyInit()
	resp, _ := c.R().
		SetHeader("Cookie", l.cookiePool["cookie2"]).
		Get("https://www.luogu.com.cn/auth/login")
	cookie := resp.Cookies()
	l.cookiePool["cookie1"] = cookie[0].Name + "=" + cookie[0].Value
}

func (l *LuoguUpdateCookie) restyInit() *resty.Client {
	c := resty.New()
	c.SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetHeader("User-Agent", config.AppConfig.Luogu.UserAgent)
	return c
}
func (l *LuoguUpdateCookie) getCaptcha() {
	c := l.restyInit()
	now := time.Now()
	stamp := float64(now.UnixMicro()) / 1000
	resp, _ := c.R().
		Get(fmt.Sprintf("https://www.luogu.com.cn/lg4/captcha?_t=%f", stamp))
	newCookie := resp.Cookies()
	l.cookiePool["cookie2"] = newCookie[0].Name + "=" + newCookie[0].Value
}
func (l *LuoguUpdateCookie) redirCaptcha() {
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
func (l *LuoguUpdateCookie) saveImage(content []byte) {
	file, err := os.OpenFile("ocr/captcha.jpg", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write(content)
}
func (l *LuoguUpdateCookie) identify(isInServer bool) (string, error) {
	var cmd *exec.Cmd
	var expectedLength int

	if isInServer {
		cmd = exec.Command("python3", "./ocr/main.py")
		expectedLength = constants.CaptchaLengthServer
	} else {
		cmd = exec.Command("python", "./ocr/main.py")
		expectedLength = constants.CaptchaLengthLocal
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if isInServer {
			fmt.Println("请检查是否是服务器环境！")
		} else {
			fmt.Println("请检查是否是测试环境以及虚拟环境是否激活，关键包：ddddocr\n激活指令：./ocr/.venv/Scripts/activate")
		}
		return "", err
	}

	str := string(output)
	// str = strings.TrimSpace(str)

	// 安全地提取验证码
	if len(str) >= expectedLength {
		startIndex := len(str) - expectedLength - 1
		if startIndex < 0 {
			startIndex = 0
		}
		endIndex := len(str) - 1
		if endIndex > len(str) {
			endIndex = len(str)
		}
		return str[startIndex:endIndex], nil
	}

	return str, nil
}
func (l *LuoguUpdateCookie) login(captcha string) error {
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
