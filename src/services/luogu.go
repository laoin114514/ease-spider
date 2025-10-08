package services

import (
	"fmt"
	"math"
	"spider/config/db"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	updatecookie "spider/src/utils/updateCookie"
	"strconv"
	"time"
)

const (
	LuoguStatusAccepted = 12
)

type Luogu struct {
	repo *repository.LuoguRepository
	req  *utils.Request[models.LuoguRecordsResponse]
	log  *utils.LogContainer
}

func NewLuoguService() *Luogu {
	return &Luogu{
		repo: repository.NewLuoguRepository(),
		req:  utils.NewRequest[models.LuoguRecordsResponse](),
		log:  utils.NewLogContainer(),
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
	for i := 1; i <= page; i++ {
		data, _ := l.req.Get(
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
func (l *Luogu) Clear() error {
	l.log.ClearLog()
	l.log.ClearErr()
	return nil
}

// ================================更新洛谷Cookie===============================================
func (l *Luogu) UpdateLuoguCookie() error {
	updatecookie.Use(false)
	return nil
}
