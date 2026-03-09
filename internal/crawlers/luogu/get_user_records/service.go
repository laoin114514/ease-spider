package getuserrecords

import (
	"fmt"
	"math"
	"net/http"
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/internal/utils"
	easecrawler "spider/pkg/ease-crawler"
	"strconv"
	"time"
)

// GetLuoguUsersRecords 获取洛谷用户提交记录
func (g *GetUserRecords) GetAndStore() error {
	conCurrenter := easecrawler.NewConCurrenter[models.LuoguUserDeliver](config.AppConfig.Luogu.LuoguRecordsConcurrency)
	conCurrenter.SetLogger(g.log)
	luoguUserDelivers, err := g.repo.GetUserNameMap()
	if err != nil {
		return err
	}

	// 通过并发器来获取洛谷用户提交记录
	conCurrenter.Run(luoguUserDelivers, func(luoguUser models.LuoguUserDeliver) error {
		return g.processUserRecords(luoguUser)
	})

	return nil
}

// ChangePrivateProblem 修改私有题目难度为unknown - 对外提供的接口
func (g *GetUserRecords) ChangePrivateProblem() error {
	count, err := g.repo.ChangePrivateProblem()
	if err != nil {
		g.log.Errorf("修改私有题目难度为unknow失败 %s", err.Error())
		return err
	}
	if count == 0 {
		return nil
	}
	g.log.Printf("修改私有题目难度为unknow完成 %d", count)
	return nil
}

// processUserRecords 处理单个用户的提交记录 - 私有方法
func (g *GetUserRecords) processUserRecords(luoguUser models.LuoguUserDeliver) error {

	// 获取初始化数据：总数和每页数量
	initData, err := g.fetchInitialData(luoguUser.Uid)
	if err != nil {
		g.log.Errorf("%s获取提交记录失败 %s", luoguUser.RealName, err.Error())
		return err
	}
	if initData.Code == http.StatusNotFound {
		g.log.Errorf("%s的uid不存在", luoguUser.RealName)
		return fmt.Errorf("用户uid不存在")
	}

	// 计算页数
	page := g.calculatePage(initData)

	// 增量爬取
	err = g.fetchRecordsIncrementally(&luoguUser, page)
	// if err != nil {
	// 	l.logError(luoguUser.RealName, "获取提交记录失败", err)
	// }

	// 如果已爬取数据与数据库已爬取数据数量一致，则不进行全量爬取
	if g.isDataConsistent(luoguUser, initData.CurrentData.Records.Count) {
		if luoguUser.Count == 0 {
			return nil
		}
		g.log.Printf("%s获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count)
		return err
	}

	// 爬取数据与数据库已爬取数据不一致，进行全量爬取
	g.logDataInconsistency(luoguUser, initData.CurrentData.Records.Count)
	err = g.fetchRecordsFully(&luoguUser, page)
	if err != nil {
		g.log.Errorf("%s重新获取提交记录失败 %s", luoguUser.RealName, err.Error())
		return err
	}

	if luoguUser.Count == 0 {
		g.log.Printf("%s重新获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count)
		return nil
	}
	g.log.Printf("%s重新获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count)
	return err
}

// fetchInitialData 获取初始数据 - 私有方法
func (g *GetUserRecords) fetchInitialData(uid string) (*models.LuoguRecordsResponse, error) {
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
func (g *GetUserRecords) fetchRecordsIncrementally(luoguUser *models.LuoguUserDeliver, page int) error {
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

		err = g.processPageRecords(data.CurrentData.Records.Result, luoguUser, i, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// fetchRecordsFully 全量爬取所有页数的数据
func (g *GetUserRecords) fetchRecordsFully(luoguUser *models.LuoguUserDeliver, page int) error {
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
		g.processPageRecords(data.CurrentData.Records.Result, luoguUser, i, false)
	}
	return nil
}

// processPageRecords 处理单页记录
func (g *GetUserRecords) processPageRecords(records []models.LuoguRecord, luoguUser *models.LuoguUserDeliver, page int, strictMode bool) error {
	for _, record := range records {
		if g.shouldSkipRecord(&record, luoguUser, page) {
			continue
		}

		if luoguUser.OldDataSet[strconv.Itoa(int(record.ID))] {
			if strictMode {
				return fmt.Errorf("%s第%d页提交记录已存在 %s", luoguUser.RealName, page, strconv.Itoa(int(record.ID)))
			}
			continue
		}

		table := g.buildSubmissionTable(&record)
		err := db.Insert_luogu_sub(table)
		if err != nil {
			if strictMode {
				return fmt.Errorf("%s处理第%d页提交记录失败 %s", luoguUser.RealName, page, err.Error())
			}
			g.log.Errorf("%s第%d页提交记录插入失败 %s", luoguUser.RealName, page, err.Error())
			continue
		}
		luoguUser.Count++
	}
	return nil
}

// shouldSkipRecord 判断是否应该跳过该提交记录
func (g *GetUserRecords) shouldSkipRecord(record *models.LuoguRecord, luoguUser *models.LuoguUserDeliver, page int) bool {
	// 如果提交记录状态为还在测评，则跳过
	nowTime := time.Now().Unix()
	if record.Status == 0 && nowTime-record.SubmitTime < 60*5 {
		g.log.Warnf("%s第%d页提交记录还在测评 %s||https://www.luogu.com.cn/record/%d", luoguUser.RealName, page, record.Problem.PID, record.ID)
		return true
	}
	// 如果洛谷UID不匹配，则更新洛谷UID
	if int(record.User.UID) == 0 {
		record.User.UID, _ = strconv.ParseInt(luoguUser.Uid, 10, 64)
	}
	return false
}

// buildSubmissionTable 构建提交记录表
func (g *GetUserRecords) buildSubmissionTable(record *models.LuoguRecord) db.Luogu_all_submissions {
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
func (g *GetUserRecords) calculatePage(luoguRecordsResponse *models.LuoguRecordsResponse) int {
	return int(math.Ceil(float64(luoguRecordsResponse.CurrentData.Records.Count) / float64(luoguRecordsResponse.CurrentData.Records.PerPage)))
}

// isDataConsistent 检查数据一致性
func (g *GetUserRecords) isDataConsistent(luoguUser models.LuoguUserDeliver, totalCount int) bool {
	return len(luoguUser.OldDataSet)+luoguUser.Count == totalCount
}

// logDataInconsistency 记录数据不一致日志
func (g *GetUserRecords) logDataInconsistency(luoguUser models.LuoguUserDeliver, actualCount int) {
	msg := fmt.Sprintf("%s爬取实际数量%d，数据库已爬取数量%d", luoguUser.RealName, actualCount, len(luoguUser.OldDataSet)+luoguUser.Count)
	g.log.Printf(msg)
}
