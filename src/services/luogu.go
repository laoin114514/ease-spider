package services

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"spider/config/db"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	LuoguStatusAccepted = 12
)

type Luogu struct {
	repo    *repository.LuoguRepository
	req     *utils.Request[models.LuoguRecordsResponse]
	log     *utils.LogContainer
	logPath string
	errPath string
}

func NewLuoguService() *Luogu {
	return &Luogu{
		repo:    repository.NewLuoguRepository(),
		req:     utils.NewRequest[models.LuoguRecordsResponse](),
		log:     utils.NewLogContainer(),
		logPath: "logs/luogu.log",
		errPath: "logs/luogu.err.log",
	}
}

// 核心函数，获取洛谷用户提交记录
func (l *Luogu) GetLuoguUsersRecords(concurrency int) error {
	conCurrenter := utils.NewConCurrenter[models.LuoguUserDeliver](concurrency)
	luoguUserDelivers, err := l.repo.GetUserNameMap()
	if err != nil {
		return err
	}

	//通过并发器来获取洛谷用户提交记录
	conCurrenter.Run(luoguUserDelivers, func(luoguUser models.LuoguUserDeliver) error {

		//获取初始化数据：总数和每页数量
		initData, err := l.req.Get("https://www.luogu.com.cn/record/list", map[string]string{"user": luoguUser.Uid, "page": "1", "_contentOnly": "1"})
		if err != nil {
			l.log.AddErr(fmt.Sprintf("%s获取提交记录失败 %s", luoguUser.RealName, err.Error()))
			return err
		}
		if initData.Code == 404 {
			l.log.AddErr(fmt.Sprintf("%s的uid不存在", luoguUser.RealName))
			return err
		}
		//计算页数
		page := l.calculatePage(&initData)

		//增量爬取
		err = l.loopRequestIncrement(&luoguUser, page)
		if err != nil {
			l.log.AddErr(fmt.Sprintf("%s获取提交记录失败 %s", luoguUser.RealName, err.Error()))
		}

		if len(luoguUser.OldDataSet)+luoguUser.Count == initData.CurrentData.Records.Count {
			l.log.AddLog(fmt.Sprintf("%s获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count))
			return err
		}

		//爬取数据与数据库已爬取数据不一致，进行全量爬取
		l.log.AddErr(fmt.Sprintf("%s爬取实际数量%d，数据库已爬取数量%d", luoguUser.RealName, initData.CurrentData.Records.Count, len(luoguUser.OldDataSet)+luoguUser.Count))
		err = l.loopRequestAll(&luoguUser, page)
		if err != nil {
			l.log.AddErr(fmt.Sprintf("%s重新获取提交记录失败 %s", luoguUser.RealName, err.Error()))
			return err
		}

		l.log.AddLog(fmt.Sprintf("%s重新获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count))
		return err
	})
	return nil
}

// 增量爬取不重复数据
func (l *Luogu) loopRequestIncrement(luoguUser *models.LuoguUserDeliver, page int) error {
	luoguUser.Count = 0
	for i := 1; i <= page; i++ {
		data, err := l.req.Get(
			"https://www.luogu.com.cn/record/list",
			map[string]string{
				"user":         luoguUser.Uid,
				"page":         strconv.Itoa(i),
				"_contentOnly": "1",
			},
		)
		if err != nil {
			return fmt.Errorf("%s获取第%d页提交记录失败 %s", luoguUser.RealName, i, err.Error())
		}
		if data.Code != 200 {
			return fmt.Errorf("%s获取第%d页提交记录失败，状态码： %d", luoguUser.RealName, i, data.Code)
		}

		for _, record := range data.CurrentData.Records.Result {
			// 如果洛谷UID不匹配，则更新洛谷UID
			if strconv.Itoa(int(record.User.UID)) != luoguUser.Uid {
				record.User.UID, _ = strconv.ParseInt(luoguUser.Uid, 10, 64)
			}

			if luoguUser.OldDataSet[strconv.Itoa(int(record.ID))] {
				return fmt.Errorf("%s第%d页提交记录已存在 %s", luoguUser.RealName, i, strconv.Itoa(int(record.ID)))
			}
			table := l.buildTable(&record)
			err = db.Insert_luogu_sub(table)
			if err != nil {
				return fmt.Errorf("%s处理第%d页提交记录失败 %s", luoguUser.RealName, i, err.Error())
			}
			luoguUser.Count++
		}
	}
	return nil
}

// 全量爬取所有页数的数据
func (l *Luogu) loopRequestAll(luoguUser *models.LuoguUserDeliver, page int) error {
	luoguUser.Count = 0
	cookie := utils.JsonDB.Get("Cookie").(string)
	for i := 1; i <= page; i++ {
		data, _ := l.req.
			SetCookie(cookie).
			Get(
				"https://www.luogu.com.cn/record/list",
				map[string]string{
					"user":         luoguUser.Uid,
					"page":         strconv.Itoa(i),
					"_contentOnly": "1",
				},
			)

		for _, record := range data.CurrentData.Records.Result {
			// 如果洛谷UID不匹配，则更新洛谷UID
			if strconv.Itoa(int(record.User.UID)) != luoguUser.Uid {
				record.User.UID, _ = strconv.ParseInt(luoguUser.Uid, 10, 64)
			}
			if luoguUser.OldDataSet[record.Problem.PID] {
				return fmt.Errorf("%s第%d页提交记录已存在 %s", luoguUser.RealName, i, record.Problem.PID)
			}
			table := l.buildTable(&record)
			db.Insert_luogu_sub(table)
			luoguUser.Count++
		}
	}
	return nil
}

// 构建提交记录表
func (l *Luogu) buildTable(record *models.LuoguRecord) db.Luogu_all_submissions {

	var difficulty = []string{"grey", "red", "brown", "yellow", "green", "blue", "purple", "black"}
	return db.Luogu_all_submissions{
		Sub_id:        fmt.Sprintf("%d", record.ID),
		Uid:           fmt.Sprintf("%d", record.User.UID),
		Problem_id:    record.Problem.PID,
		Problem_name:  record.Problem.Title,
		Difficulty:    difficulty[record.Problem.Difficulty],
		Is_pass:       record.Status == LuoguStatusAccepted,
		Creation_time: time.Unix(record.SubmitTime, 0).Add(8 * time.Hour),
	}
}

// 计算页数
func (l *Luogu) calculatePage(luoguRecordsResponse *models.LuoguRecordsResponse) int {
	return int(math.Ceil(float64(luoguRecordsResponse.CurrentData.Records.Count) / float64(luoguRecordsResponse.CurrentData.Records.PerPage)))
}

// 打印日志
func (l *Luogu) GetLog() []string {
	return l.log.GetLog()
}
func (l *Luogu) GetErr() []string {
	return l.log.GetErr()
}
func (l *Luogu) SaveLog() error {
	return SaveLog(l.logPath, l.log.GetLog())
}
func (l *Luogu) SaveErr() error {
	return SaveErr(l.errPath, l.log.GetErr())
}
func (l *Luogu) Clear() error {
	l.log.ClearLog()
	l.log.ClearErr()
	return nil
}

// ================================更新洛谷Cookie(屎山代码，勿动)===============================================
type LuoguUpdateCookie struct {
	cookiePool map[string]string
}

func NewLuoguUpdateCookie() *LuoguUpdateCookie {
	return &LuoguUpdateCookie{
		cookiePool: make(map[string]string),
	}
}
func (l *LuoguUpdateCookie) UpdateLuoguCookie() error {
	l.init()
	l.initRedirect()
	l.GetCaptcha()
	l.RedirCaptcha()
	captcha, err := l.Identify(false)
	if err != nil {
		return err
	}
	fmt.Println(captcha)
	l.Login(captcha)
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
		SetHeader("User-Agent", utils.JsonDB.Get("user_agent").(string))
	return c
}
func (l *LuoguUpdateCookie) GetCaptcha() {
	c := l.restyInit()
	now := time.Now()
	stamp := float64(now.UnixMicro()) / 1000
	resp, _ := c.R().
		Get(fmt.Sprintf("https://www.luogu.com.cn/lg4/captcha?_t=%f", stamp))
	newCookie := resp.Cookies()
	l.cookiePool["cookie2"] = newCookie[0].Name + "=" + newCookie[0].Value
}
func (l *LuoguUpdateCookie) RedirCaptcha() {
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
	file, err := os.OpenFile(fmt.Sprintf("ocr/captcha.jpg"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write(content)
}
func (l *LuoguUpdateCookie) Identify(isInServer bool) (string, error) {
	if isInServer {
		cmd := exec.Command("python3", "./ocr/main.py")
		outPut, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("请检查是否是服务器环境！")
			return "", err
		}
		str := string(outPut)
		str = str[len(str)-5 : len(str)-1]
		return str, nil
	} else {
		cmd := exec.Command("python", "./ocr/main.py")
		outPut, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("请检查是否是测试环境以及虚拟环境是否激活，关键包：ddddocr\n激活指令：./ocr/.venv/Scripts/activate")
			return "", err
		}
		str := string(outPut)
		str = str[len(str)-6 : len(str)-2]
		return str, nil
	}
}
func (l *LuoguUpdateCookie) Login(captcha string) {
	c := l.restyInit()
	cookie := l.cookiePool["cookie1"] + "; " + l.cookiePool["cookie2"]
	resp, err := c.R().
		SetBody(map[string]any{
			"username": utils.JsonDB.Get("username_luogu").(string),
			"password": utils.JsonDB.Get("password_luogu").(string),
			"captcha":  captcha,
		}).
		SetHeader("Cookie", cookie).
		Post("https://www.luogu.com.cn/do-auth/password")
	fmt.Println(resp.Status())
	if resp.StatusCode() != 200 {
		fmt.Println("验证码错误")
		return
	}
	if err != nil {
		fmt.Println(err)
	}
	arr := resp.Cookies()
	uid := arr[0].Name + "=" + arr[0].Value
	Cookie := cookie + ";" + uid
	utils.JsonDB.Set("Cookie", Cookie)
	fmt.Println("登录成功")
}
