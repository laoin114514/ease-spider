package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
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

type Luogu struct {
	LuoguSubmissionDetail *LuoguSubmissionDetail
	LuoguCookie           *LuoguUpdateCookie
	LuoguSolution         *LuoguSolution
	LuoguRecords          *LuoguRecords
}

func NewLuogu() *Luogu {
	return &Luogu{
		LuoguSubmissionDetail: NewLuoguSubmissionDetail(),
		LuoguCookie:           NewLuoguUpdateCookie(),
		LuoguSolution:         NewLuoguSolution(),
		LuoguRecords:          NewLuoguRecords(),
	}
}

// ================================获取洛谷用户提交记录===============================================
func NewLuoguRecords() *LuoguRecords {
	return &LuoguRecords{
		LogService: NewLogService("logs/luogu.log", "logs/luogu.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
	}
}

// ================================更新洛谷Cookie===============================================
func NewLuoguUpdateCookie() *LuoguUpdateCookie {
	return &LuoguUpdateCookie{
		LogService: NewLogService("logs/luoguUpdateCookie.log", "logs/luoguUpdateCookie.err.log"),
		cookiePool: make(map[string]string),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
	}
}

// ================================获取洛谷题解===============================================
func NewLuoguSolution() *LuoguSolution {
	return &LuoguSolution{
		LogService: NewLogService("logs/luoguSolution.log", "logs/luoguSolution.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		totalPage:  0,
		solutions:  []models.SolutionContent{},
	}
}

// ================================获取洛谷提交记录源代码===============================================
func NewLuoguSubmissionDetail() *LuoguSubmissionDetail {
	return &LuoguSubmissionDetail{
		LogService: NewLogService("logs/luoguSubmissionDetail.log", "logs/luoguSubmissionDetail.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		count:      0,
	}
}

// ================================获取洛谷团队成员===============================================
func NewLuoguTeam() *LuoguTeam {
	return &LuoguTeam{
		LogService: NewLogService("logs/luoguTeam.log", "logs/luoguTeam.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
		count:      0,
		results:    make(map[string]bool),
	}
}

// ================================获取洛谷题单===============================================
func NewLuoguProblemList() *LuoguProblemList {
	return &LuoguProblemList{
		LogService: NewLogService("logs/luoguProblemList.log", "logs/luoguProblemList.err.log"),
		repo:       repository.NewLuoguRepository(),
		debug:      utils.NewDebug(config.AppConfig.DebugConfig.All),
	}
}

// //
// //
// //
// //
// //
// //
// //
// //
// //
// //
// ================================获取洛谷用户提交记录===============================================
type LuoguRecords struct {
	*LogService
	repo  *repository.LuoguRepository
	debug *utils.Debug
}

// GetLuoguUsersRecords 获取洛谷用户提交记录
func (l *LuoguRecords) GetAndStore() error {
	conCurrenter := utils.NewConCurrenter[models.LuoguUserDeliver](config.AppConfig.Luogu.LuoguRecordsConcurrency)
	luoguUserDelivers, err := l.repo.GetUserNameMap()
	if err != nil {
		return err
	}

	// 通过并发器来获取洛谷用户提交记录
	conCurrenter.Run(luoguUserDelivers, func(luoguUser models.LuoguUserDeliver) error {
		return l.processUserRecords(luoguUser)
	})
	return nil
}

// ChangePrivateProblem 修改私有题目难度为unknown - 对外提供的接口
func (l *LuoguRecords) ChangePrivateProblem() error {
	count, err := l.repo.ChangePrivateProblem()
	if err != nil {
		l.debug.Debug(fmt.Sprintf("修改私有题目难度为unknow失败 %s", err.Error()))
		l.AddErr(fmt.Sprintf("修改私有题目难度为unknow失败 %s", err.Error()))
		return err
	}
	if count == 0 {
		return nil
	}
	l.debug.Debug(fmt.Sprintf("修改私有题目难度为unknow完成 %d", count))
	l.AddLog(fmt.Sprintf("修改私有题目难度为unknow完成 %d", count))
	return nil
}

// processUserRecords 处理单个用户的提交记录 - 私有方法
func (l *LuoguRecords) processUserRecords(luoguUser models.LuoguUserDeliver) error {

	// 获取初始化数据：总数和每页数量
	initData, err := l.fetchInitialData(luoguUser.Uid)
	if err != nil {
		l.logError(luoguUser.RealName, "获取提交记录失败", err)
		return err
	}
	if initData.Code == http.StatusNotFound {
		l.logError(luoguUser.RealName, "的uid不存在", nil)
		return fmt.Errorf("用户uid不存在")
	}

	// 计算页数
	page := l.calculatePage(initData)

	// 增量爬取
	err = l.fetchRecordsIncrementally(&luoguUser, page)
	if err != nil {
		l.logError(luoguUser.RealName, "获取提交记录失败", err)
	}

	// 如果已爬取数据与数据库已爬取数据数量一致，则不进行全量爬取
	if l.isDataConsistent(luoguUser, initData.CurrentData.Records.Count) {
		if luoguUser.Count == 0 {
			return nil
		}
		l.logSuccess(luoguUser.RealName, "获取提交记录完成", luoguUser.Count)
		return err
	}

	// 爬取数据与数据库已爬取数据不一致，进行全量爬取
	l.logDataInconsistency(luoguUser, initData.CurrentData.Records.Count)
	err = l.fetchRecordsFully(&luoguUser, page)
	if err != nil {
		l.logError(luoguUser.RealName, "重新获取提交记录失败", err)
		return err
	}

	if luoguUser.Count == 0 {
		l.debug.Debug(fmt.Sprintf("%s重新获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count))
		return nil
	}
	l.logSuccess(luoguUser.RealName, "重新获取提交记录完成", luoguUser.Count)
	return err
}

// fetchInitialData 获取初始数据 - 私有方法
func (l *LuoguRecords) fetchInitialData(uid string) (*models.LuoguRecordsResponse, error) {
	req := utils.NewRequest[models.LuoguRecordsResponse](true)
	data, err := req.SetCookie(utils.JsonDB.Get("Cookie").(string)).Get(
		"https://www.luogu.com.cn/record/list",
		map[string]string{
			"user":         uid,
			"page":         "1",
			"_contentOnly": "1",
		},
	)
	return &data, err
}

// fetchRecordsIncrementally 增量爬取不重复数据 - 私有方法
func (l *LuoguRecords) fetchRecordsIncrementally(luoguUser *models.LuoguUserDeliver, page int) error {
	luoguUser.Count = 0
	cookie := utils.JsonDB.Get("Cookie").(string)
	req := utils.NewRequest[models.LuoguRecordsResponse](true)

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

		err = l.processPageRecords(data.CurrentData.Records.Result, luoguUser, i, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// fetchRecordsFully 全量爬取所有页数的数据
func (l *LuoguRecords) fetchRecordsFully(luoguUser *models.LuoguUserDeliver, page int) error {
	luoguUser.Count = 0
	cookie := utils.JsonDB.Get("Cookie").(string)
	req := utils.NewRequest[models.LuoguRecordsResponse](true)

	for i := 1; i <= page; i++ {
		data, _ := req.SetCookie(cookie).Get(
			"https://www.luogu.com.cn/record/list",
			map[string]string{
				"user":         luoguUser.Uid,
				"page":         strconv.Itoa(i),
				"_contentOnly": "1",
			},
		)
		l.processPageRecords(data.CurrentData.Records.Result, luoguUser, i, false)
	}
	return nil
}

// processPageRecords 处理单页记录
func (l *LuoguRecords) processPageRecords(records []models.LuoguRecord, luoguUser *models.LuoguUserDeliver, page int, strictMode bool) error {
	for _, record := range records {
		if l.shouldSkipRecord(&record, luoguUser, page) {
			continue
		}

		if luoguUser.OldDataSet[strconv.Itoa(int(record.ID))] {
			if strictMode {
				return fmt.Errorf("%s第%d页提交记录已存在 %s", luoguUser.RealName, page, strconv.Itoa(int(record.ID)))
			}
			continue
		}

		table := l.buildSubmissionTable(&record)
		err := db.Insert_luogu_sub(table)
		if err != nil {
			if strictMode {
				return fmt.Errorf("%s处理第%d页提交记录失败 %s", luoguUser.RealName, page, err.Error())
			}
			l.debug.Debug(fmt.Sprintf("%s第%d页提交记录插入失败 %s", luoguUser.RealName, page, err.Error()))
			continue
		}
		luoguUser.Count++
	}
	return nil
}

// shouldSkipRecord 判断是否应该跳过该提交记录
func (l *LuoguRecords) shouldSkipRecord(record *models.LuoguRecord, luoguUser *models.LuoguUserDeliver, page int) bool {
	// 如果提交记录状态为还在测评，则跳过
	nowTime := time.Now().Unix()
	if record.Status == 0 && nowTime-record.SubmitTime < 60*5 {
		l.debug.Debug(fmt.Sprintf("%s第%d页提交记录还在测评 %s||https://www.luogu.com.cn/record/%d", luoguUser.RealName, page, record.Problem.PID, record.ID))
		return true
	}
	// 如果洛谷UID不匹配，则更新洛谷UID
	if int(record.User.UID) == 0 {
		record.User.UID, _ = strconv.ParseInt(luoguUser.Uid, 10, 64)
	}
	return false
}

// buildSubmissionTable 构建提交记录表
func (l *LuoguRecords) buildSubmissionTable(record *models.LuoguRecord) db.Luogu_all_submissions {
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

// calculatePage 计算页数
func (l *LuoguRecords) calculatePage(luoguRecordsResponse *models.LuoguRecordsResponse) int {
	return int(math.Ceil(float64(luoguRecordsResponse.CurrentData.Records.Count) / float64(luoguRecordsResponse.CurrentData.Records.PerPage)))
}

// isDataConsistent 检查数据一致性
func (l *LuoguRecords) isDataConsistent(luoguUser models.LuoguUserDeliver, totalCount int) bool {
	return len(luoguUser.OldDataSet)+luoguUser.Count == totalCount
}

// logError 记录错误日志
func (l *LuoguRecords) logError(userName, message string, err error) {
	errorMsg := fmt.Sprintf("%s%s", userName, message)
	if err != nil {
		errorMsg += fmt.Sprintf(" %s", err.Error())
	}
	l.AddErr(errorMsg)
	l.debug.Debug(errorMsg)
}

// logSuccess 记录成功日志
func (l *LuoguRecords) logSuccess(userName, message string, count int) {
	successMsg := fmt.Sprintf("%s%s %d", userName, message, count)
	l.AddLog(successMsg)
	l.debug.Debug(successMsg)
}

// logDataInconsistency 记录数据不一致日志
func (l *LuoguRecords) logDataInconsistency(luoguUser models.LuoguUserDeliver, actualCount int) {
	msg := fmt.Sprintf("%s爬取实际数量%d，数据库已爬取数量%d", luoguUser.RealName, actualCount, len(luoguUser.OldDataSet)+luoguUser.Count)
	l.debug.Debug(msg)
	l.AddErr(msg)
}

// ================================更新洛谷Cookie===============================================
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

// =============================================爬取洛谷题解===============================================
type LuoguSolution struct {
	*LogService
	repo      *repository.LuoguRepository
	totalPage int
	debug     *utils.Debug
	solutions []models.SolutionContent
	count     int
}

// 公共方法
func (s *LuoguSolution) GetAndStore() error {
	problemIds, err := s.repo.GetProblemIdHasSourceCode()
	if err != nil {
		s.log.AddErr(fmt.Sprintf("获取题解ID失败 %s", err.Error()))
		s.debug.Debug(fmt.Sprintf("获取题解ID失败 %s", err.Error()))
		return err
	}
	//不采取并发，因为题解数量较少，且获取题解时间较长
	for _, problemId := range problemIds {
		solutions, err := s.GetSolutionListByProblemId(problemId)
		s.debug.Debug(fmt.Sprintf("获取题解 题目ID：%s 目前题解数量：%d", problemId, len(solutions)))
		if err != nil {
			s.log.AddErr(fmt.Sprintf("获取题解失败 %s", err.Error()))
			s.debug.Debug(fmt.Sprintf("获取题解失败 %s", err.Error()))
			continue
		}
		err = s.storeSolution(solutions)
		if err != nil {
			s.log.AddErr(fmt.Sprintf("存储题解失败 %s", err.Error()))
			s.debug.Debug(fmt.Sprintf("存储题解失败 %s", err.Error()))
			continue
		}
	}
	s.debug.Debug(fmt.Sprintf("获取题解完成 %d", s.count))
	s.log.AddLog(fmt.Sprintf("获取题解完成 %d", s.count))
	return nil
}

// 获取全部题解
func (s *LuoguSolution) GetSolutionListByProblemId(problemID string) ([]models.SolutionContent, error) {
	_, err := s.analyzeSolution(problemID, 1)
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
		solutions, err := s.analyzeSolution(problemID, page)
		if err != nil {
			return err
		}
		s.solutions = append(s.solutions, solutions...)
		return nil
	})
	return s.solutions, nil
}
func (s *LuoguSolution) storeSolution(solutions []models.SolutionContent) error {
	for _, solution := range solutions {
		table := db.Luogu_solutions{
			Problem_id:    solution.SolutionFor.PID,
			Problem_name:  solution.SolutionFor.Title,
			Author_name:   solution.Author.Name,
			Author_uid:    solution.Author.UID,
			Solution_md:   solution.Content,
			Creation_time: time.Unix(solution.Time, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
			Link:          fmt.Sprintf("https://www.luogu.com.cn/problem/solution/%s?page=%d", solution.SolutionFor.PID, solution.Collection["page"].(int)),
		}
		err := db.Insert_luogu_solutions(table)
		if err != nil {
			continue
		}
		s.count++
	}
	return nil
}

// 获取一页题解及详细信息
func (s *LuoguSolution) analyzeSolution(problemID string, page int) ([]models.SolutionContent, error) {
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
	for i := range response.Data.Solutions.Result {
		response.Data.Solutions.Result[i].Collection = map[string]any{
			"page": page,
		}
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
	repo   *repository.LuoguRepository
	debug  *utils.Debug
	cookie string
	count  int
}

// 获取源代码
func (s *LuoguSubmissionDetail) GetAndStoreSourceCode() error {
	cookie := utils.JsonDB.Get("Cookie")
	if cookie == nil {
		return errors.New("cookie不存在")
	}
	s.cookie = cookie.(string)
	luoguTeam := NewLuoguTeam()
	err := luoguTeam.GetMembers()
	if err != nil {
		s.AddErr(fmt.Sprintf("获取洛谷团队成员失败 %s", err.Error()))
		s.debug.Debug("获取洛谷团队成员失败")
		return err
	}

	subids, err := s.repo.GetSubidNoSourceCode(luoguTeam.results)
	if err != nil {
		s.AddErr(fmt.Sprintf("获取提交记录ID失败 %s", err.Error()))
		s.debug.Debug("获取提交记录ID失败")
		return err
	}
	s.debug.Debug(fmt.Sprintf("总共 %d个提交记录需要获取源代码\n团队成员：%v", len(subids), luoguTeam.results))

	// 创建并发器
	conCurrenter := utils.NewConCurrenter[string](config.AppConfig.Luogu.LuoguSubmissionDetailConcurrency)

	// 并发获取提交记录源代码
	err = conCurrenter.Run(subids, func(subid string) error {
		name, err := s.repo.GetNameBySubid(subid)
		if err != nil {
			s.AddErr(fmt.Sprintf("获取提交记录名称失败 %s", err.Error()))
			s.debug.Debug(subid + name + "获取失败" + err.Error())
			return err
		}

		respJson, err := s.getSourceCodeBySubid(subid)
		if err != nil {
			s.AddErr(fmt.Sprintf("解析提交记录源代码失败 %s", err.Error()))
			s.debug.Debug(subid + name + "解析失败" + err.Error())
			luoUpdateCookieService := NewLuoguUpdateCookie()
			luoUpdateCookieService.Update()
			return err
		}

		//源代码长度小于5，则填充无
		if len(respJson.CurrentData.Record.SourceCode) <= 5 {
			s.AddErr(fmt.Sprintf("提交记录源代码为空 %s", subid))
			respJson.CurrentData.Record.SourceCode = "无"
			s.debug.Debug(subid + name + "源代码为空")
			if respJson.CurrentData.Record.Problem.Type == "P" && respJson.CurrentData.Record.Status == constants.LuoguStatusAccepted {
				s.debug.Debug(subid + name + "AC且无源代码，跳过")
				return nil
			}
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

func (s *LuoguSubmissionDetail) getSourceCodeBySubid(subid string) (models.LuoguSubmissionDetailResponse, error) {
	req := utils.NewRequest[models.LuoguSubmissionDetailResponse](false)
	url := fmt.Sprintf("https://www.luogu.com.cn/record/%s?_contentOnly=1", subid)
	req.SetCookie(s.cookie)
	resp, err := req.Get(url, map[string]string{})
	if err != nil {
		return models.LuoguSubmissionDetailResponse{}, err
	}
	if resp.Code != 200 {
		return models.LuoguSubmissionDetailResponse{}, errors.New("获取提交记录源代码失败" + strconv.Itoa(resp.Code))
	}
	return resp, nil
}

// ============================================获取洛谷团队成员==========================================
type LuoguTeam struct {
	*LogService
	repo    *repository.LuoguRepository
	debug   *utils.Debug
	results map[string]bool
	count   int
}

func (l *LuoguTeam) GetMembers() error {
	req := utils.NewRequest[models.LuoguTeamResponse](true)
	url := fmt.Sprintf("https://www.luogu.com.cn/api/team/members/%d", config.AppConfig.Luogu.LuoguTeamID)
	resp, err := req.Get(url, map[string]string{})
	if err != nil {
		return err
	}
	for _, member := range resp.Members.Result {
		l.results[strconv.FormatInt(member.User.UID, 10)] = true
		l.count++
	}
	l.debug.Debug(fmt.Sprintf("获取洛谷团队成员完成 %d", l.count))
	l.AddLog(fmt.Sprintf("获取洛谷团队成员完成 %d", l.count))
	return nil
}

// ============================================获取洛谷题单列表==========================================
type LuoguProblemList struct {
	*LogService
	repo  *repository.LuoguRepository
	debug *utils.Debug
}
