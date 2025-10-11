package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"spider/config"
	"spider/config/db"
	"spider/src/constants"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type LuoguRecords struct {
	*LogService
	repo *repository.LuoguRepository
}

func NewLuoguRecordsService() *LuoguRecords {
	return &LuoguRecords{
		LogService: NewLogService("logs/luogu.log", "logs/luogu.err.log"),
		repo:       repository.NewLuoguRepository(),
	}
}

// 核心函数，获取洛谷用户提交记录
func (l *LuoguRecords) GetLuoguUsersRecords() error {
	conCurrenter := utils.NewConCurrenter[models.LuoguUserDeliver](config.AppConfig.Luogu.LuoguRecordsConcurrency)
	luoguUserDelivers, err := l.repo.GetUserNameMap()
	if err != nil {
		return err
	}

	//通过并发器来获取洛谷用户提交记录
	conCurrenter.Run(luoguUserDelivers, func(luoguUser models.LuoguUserDeliver) error {
		req := utils.NewRequest[models.LuoguRecordsResponse]()
		//获取初始化数据：总数和每页数量
		initData, err := req.SetCookie(utils.JsonDB.Get("Cookie").(string)).Get("https://www.luogu.com.cn/record/list", map[string]string{"user": luoguUser.Uid, "page": "1", "_contentOnly": "1"})
		if err != nil {
			l.AddErr(fmt.Sprintf("%s获取提交记录失败 %s", luoguUser.RealName, err.Error()))
			return err
		}
		if initData.Code == http.StatusNotFound {
			l.AddErr(fmt.Sprintf("%s的uid不存在", luoguUser.RealName))
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
			l.AddLog(fmt.Sprintf("%s获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count))
			return err
		}

		//爬取数据与数据库已爬取数据不一致，进行全量爬取
		l.AddErr(fmt.Sprintf("%s爬取实际数量%d，数据库已爬取数量%d", luoguUser.RealName, initData.CurrentData.Records.Count, len(luoguUser.OldDataSet)+luoguUser.Count))
		err = l.loopRequestAll(&luoguUser, page)
		if err != nil {
			l.AddErr(fmt.Sprintf("%s重新获取提交记录失败 %s", luoguUser.RealName, err.Error()))
			return err
		}

		l.AddLog(fmt.Sprintf("%s重新获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count))
		return err
	})
	return nil
}

// 增量爬取不重复数据
func (l *LuoguRecords) loopRequestIncrement(luoguUser *models.LuoguUserDeliver, page int) error {
	luoguUser.Count = 0
	cookie := utils.JsonDB.Get("Cookie").(string)
	req := utils.NewRequest[models.LuoguRecordsResponse]()
	for i := 1; i <= page; i++ {
		data, err := req.SetCookie(cookie).Get(
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
func (l *LuoguRecords) loopRequestAll(luoguUser *models.LuoguUserDeliver, page int) error {
	luoguUser.Count = 0
	cookie := utils.JsonDB.Get("Cookie").(string)
	req := utils.NewRequest[models.LuoguRecordsResponse]()
	for i := 1; i <= page; i++ {
		data, _ := req.
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
func (l *LuoguRecords) buildTable(record *models.LuoguRecord) db.Luogu_all_submissions {
	difficulty := constants.LuoguDifficultyMap
	if record.Problem.Difficulty >= len(difficulty) {
		difficulty = append(difficulty, "unknown")
	}

	return db.Luogu_all_submissions{
		Sub_id:        fmt.Sprintf("%d", record.ID),
		Uid:           fmt.Sprintf("%d", record.User.UID),
		Problem_id:    record.Problem.PID,
		Problem_name:  record.Problem.Title,
		Difficulty:    difficulty[record.Problem.Difficulty],
		Is_pass:       record.Status == constants.LuoguStatusAccepted,
		Creation_time: time.Unix(record.SubmitTime, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}

// 计算页数
func (l *LuoguRecords) calculatePage(luoguRecordsResponse *models.LuoguRecordsResponse) int {
	return int(math.Ceil(float64(luoguRecordsResponse.CurrentData.Records.Count) / float64(luoguRecordsResponse.CurrentData.Records.PerPage)))
}

// ================================更新洛谷Cookie===============================================
type LuoguUpdateCookie struct {
	*LogService
	cookiePool map[string]string
}

func NewLuoguUpdateCookie() *LuoguUpdateCookie {
	return &LuoguUpdateCookie{
		LogService: NewLogService("logs/luoguUpdateCookie.log", "logs/luoguUpdateCookie.err.log"),
		cookiePool: make(map[string]string),
	}
}
func (l *LuoguUpdateCookie) UpdateLuoguCookie() error {
	l.init()
	l.initRedirect()
	l.getCaptcha()
	l.redirCaptcha()
	captcha, err := l.identify(config.AppConfig.Luogu.IsInServer)
	if err != nil {
		l.AddLog(fmt.Sprintf("验证码识别失败 %s", err.Error()))
		return err
	}
	l.AddLog(fmt.Sprintf("验证码：%s", captcha))
	err = l.login(captcha)
	if err != nil {
		l.AddLog(fmt.Sprintf("登录失败 %s", err.Error()))
		return err
	}
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

// =============================================爬取洛谷题解===============================================
type LuoguSolution struct {
	*LogService
	repo      *repository.LuoguRepository
	totalPage int
	solutions []models.SolutionContent
}

func NewLuoguSolution() *LuoguSolution {
	return &LuoguSolution{
		LogService: NewLogService("logs/luoguSolution.log", "logs/luoguSolution.err.log"),
		repo:       repository.NewLuoguRepository(),
		totalPage:  0,
		solutions:  []models.SolutionContent{},
	}
}

// 获取全部题解
func (s *LuoguSolution) GetSolutionList(problemID string) ([]models.SolutionContent, error) {
	_, err := s.GetSolution(problemID, 1)
	if err != nil {
		return nil, err
	}
	pageRange := []int{}
	for i := 1; i <= s.totalPage; i++ {
		pageRange = append(pageRange, i)
	}
	//初始化并发器
	conCurrenter := utils.NewConCurrenter[int](config.AppConfig.Luogu.LuoguSolutionConcurrency)
	conCurrenter.Run(pageRange, func(page int) error {
		solutions, err := s.GetSolution(problemID, page)
		if err != nil {
			return err
		}
		s.solutions = append(s.solutions, solutions...)
		return nil
	})
	return s.solutions, nil
}

// 获取一页题解及详细信息
func (s *LuoguSolution) GetSolution(problemID string, page int) ([]models.SolutionContent, error) {
	//如果页数大于总页数，并且不是第一页，则返回错误
	if page > s.totalPage && page != 1 {
		return nil, errors.New("页数超出范围")
	}
	client := resty.New()
	cookie := utils.JsonDB.Get("Cookie")
	if cookie == nil {
		return nil, errors.New("cookie不存在")
	}
	url := fmt.Sprintf("https://www.luogu.com.cn/problem/solution/%s?page=%d", problemID, page)
	resp, err := client.R().
		SetHeader("Cookie", cookie.(string)).
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("http code错误" + strconv.Itoa(resp.StatusCode()))
	}
	//从html提取json数据
	response, err := s.parseHTMLJSON(string(resp.Body()))
	if err != nil {
		return nil, err
	}
	totalPage := int(math.Ceil(float64(response.Data.Solutions.Count) / float64(response.Data.Solutions.PerPage)))
	s.totalPage = totalPage
	return response.Data.Solutions.Result, nil
}

// 解析html
func (s *LuoguSolution) parseHTMLJSON(htmlContent string) (models.LuoguSolutionResponse, error) {
	// 使用正则表达式提取JSON数据
	re := regexp.MustCompile(`<script id="lentille-context" type="application/json">\s*(\{[\s\S]*?\})\s*</script>`)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) < 2 {
		return models.LuoguSolutionResponse{}, fmt.Errorf("未找到JSON数据")
	}

	jsonStr := matches[1]

	// 解析JSON
	var pageData models.LuoguSolutionResponse
	err := json.Unmarshal([]byte(jsonStr), &pageData)
	if err != nil {
		return models.LuoguSolutionResponse{}, fmt.Errorf("解析JSON失败: %v", err)
	}
	return pageData, nil
}

// ============================================爬取洛谷提交记录源代码===============================================
type LuoguSubmissionDetail struct {
	*LogService
	repo  *repository.LuoguRepository
	debug *utils.Debug
	count int
}

func NewLuoguSubmissionDetail() *LuoguSubmissionDetail {
	return &LuoguSubmissionDetail{
		LogService: NewLogService("logs/luoguSubmissionDetail.log", "logs/luoguSubmissionDetail.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		count:      0,
	}
}

// 获取源代码
func (s *LuoguSubmissionDetail) GetRecordSourceCode() error {
	subids, err := s.repo.GetSubidNoSourceCode()
	if err != nil {
		s.AddErr(fmt.Sprintf("获取提交记录ID失败 %s", err.Error()))
		s.debug.Debug("获取提交记录ID失败")
		return err
	}
	s.debug.Debug(fmt.Sprintf("总共 %d个提交记录需要获取源代码", len(subids)))
	conCurrenter := utils.NewConCurrenter[string](config.AppConfig.Luogu.LuoguSubmissionDetailConcurrency)
	err = conCurrenter.Run(subids, func(subid string) error {
		name, err := s.repo.GetNameBySubid(subid)
		if err != nil {
			s.AddErr(fmt.Sprintf("获取提交记录名称失败 %s", err.Error()))
			s.debug.Debug(subid + name + "获取失败" + err.Error())
			return err
		}
		html, err := s.getRecordSourceCodeHTML(subid)
		if err != nil {
			s.AddErr(fmt.Sprintf("获取提交记录源代码失败 %s", err.Error()))
			s.debug.Debug(subid + name + "获取失败" + err.Error())

			return err
		}
		// 解析HTML中的JSON数据
		respJson, err := s.parseSourceCodeFromHTML(html)
		if err != nil {
			s.AddErr(fmt.Sprintf("解析提交记录源代码失败 %s", err.Error()))
			s.debug.Debug(subid + name + "解析失败" + err.Error())
			luoUpdateCookieService := NewLuoguUpdateCookie()
			luoUpdateCookieService.UpdateLuoguCookie()
			return err
		}
		//源代码长度小于5，则填充无
		if len(respJson.CurrentData.Record.SourceCode) <= 5 {
			s.AddErr(fmt.Sprintf("提交记录源代码为空 %s", subid))
			respJson.CurrentData.Record.SourceCode = "无"
			s.debug.Debug(subid + name + "源代码为空")
		}
		err = s.repo.InsertSourceCode(subid, respJson.CurrentData.Record.SourceCode)
		if err != nil {
			s.AddErr(fmt.Sprintf("插入提交记录源代码失败 %s", err.Error()))
			s.debug.Debug(subid + name + "插入失败" + err.Error())

			return err
		}
		s.debug.Debug(subid + name + "插入成功")
		s.count++
		return nil
	})
	s.debug.Debug(fmt.Sprintf("插入提交记录源代码完成 %d", s.count))
	s.AddLog(fmt.Sprintf("插入提交记录源代码完成 %d", s.count))
	return err
}

// 获取提交记录源代码html
func (s *LuoguSubmissionDetail) getRecordSourceCodeHTML(recordID string) (string, error) {
	client := resty.New()
	cookie := utils.JsonDB.Get("Cookie")
	if cookie == nil {
		return "", errors.New("cookie不存在")
	}
	url := fmt.Sprintf("https://www.luogu.com.cn/record/%s", recordID)
	resp, err := client.R().
		SetHeader("Cookie", cookie.(string)).
		SetQueryParams(map[string]string{}).Get(url)

	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", errors.New("http code错误" + strconv.Itoa(resp.StatusCode()))
	}
	return string(resp.Body()), nil
}

// 解析html中的json数据
func (s *LuoguSubmissionDetail) parseSourceCodeFromHTML(htmlContent string) (models.LuoguSubmissionDetailResponse, error) {
	// 使用正则表达式提取URL编码的JSON字符串
	re := regexp.MustCompile(`decodeURIComponent\("([^"]+)"\)`)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) < 2 {
		return models.LuoguSubmissionDetailResponse{}, errors.New("未找到URL编码的JSON数据")
	}
	// 获取URL编码的JSON字符串
	encodedJSON := matches[1]
	// URL解码JSON
	decodedJSON, err := url.QueryUnescape(encodedJSON)
	if err != nil {
		return models.LuoguSubmissionDetailResponse{}, errors.New("JSON URL解码失败: " + err.Error())
	}
	var realJson models.LuoguSubmissionDetailResponse
	json.Unmarshal([]byte(decodedJSON), &realJson)
	return realJson, nil
}
