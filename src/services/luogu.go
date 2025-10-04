package services

import (
	"fmt"
	"log"
	"math"
	"spider/config/db"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"strconv"
	"time"
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
	conCurrenter.Run(luoguUserDelivers, func(luoguUser models.LuoguUserDeliver) {

		//获取初始化数据：总数和每页数量
		initData, err := l.req.Get("https://www.luogu.com.cn/record/list", map[string]string{"user": luoguUser.Uid, "page": "1", "_contentOnly": "1"})
		if err != nil {
			l.log.AddErr(fmt.Sprintf("%s获取提交记录失败 %s", luoguUser.RealName, err.Error()))
			return
		}
		page := int(math.Ceil(float64(initData.CurrentData.Records.Count) / float64(initData.CurrentData.Records.PerPage)))

		err = l.loopRequest(&luoguUser, page, false)
		if err != nil {
			l.log.AddErr(fmt.Sprintf("%s获取提交记录失败 %s", luoguUser.RealName, err.Error()))
		}
		recordsLengthInDb, err := l.repo.GetUserRecordsLength(luoguUser.Uid)
		if err != nil {
			l.log.AddErr(fmt.Sprintf("%s获取提交记录数量失败 %s", luoguUser.RealName, err.Error()))
			return
		}
		if recordsLengthInDb < initData.CurrentData.Records.Count {
			l.log.AddErr(fmt.Sprintf("%s少爬取实际数量%d，数据库数量%d，重新获取中。。。", luoguUser.RealName, initData.CurrentData.Records.Count, recordsLengthInDb))
			err = l.loopRequest(&luoguUser, page, true)
			if err != nil {
				l.log.AddErr(fmt.Sprintf("%s重新获取提交记录失败 %s", luoguUser.RealName, err.Error()))
				return
			}
		}
		l.log.AddLog(fmt.Sprintf("%s获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count))
	})
	return nil
}

// 循环获取所有页数的数据
func (l *Luogu) loopRequest(luoguUser *models.LuoguUserDeliver, page int, isAllcatch bool) error {
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
		if isAllcatch {
			if err != nil {
				return fmt.Errorf("%s获取第%d页提交记录失败 %s", luoguUser.RealName, i, err.Error())
			}
			if data.Code != 200 {
				return fmt.Errorf("%s获取第%d页提交记录失败，状态码： %d", luoguUser.RealName, i, data.Code)
			}
		}

		for _, record := range data.CurrentData.Records.Result {
			// 如果洛谷UID不匹配，则更新洛谷UID
			if strconv.Itoa(int(record.User.UID)) != luoguUser.Uid {
				record.User.UID, _ = strconv.ParseInt(luoguUser.Uid, 10, 64)
			}
			table := l.buildTable(&record)
			err = db.Insert_luogu_sub(table)
			if err != nil && !isAllcatch {
				return fmt.Errorf("%s处理第%d页提交记录失败 %s", luoguUser.RealName, i, err.Error())
			} else if err != nil && isAllcatch {
				continue
			}
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
		Is_pass:       record.Status == 12,
		Creation_time: time.Unix(record.SubmitTime, 0).Add(8 * time.Hour),
	}
}

// 打印日志
func (l *Luogu) Log() {
	log.Println("日志", l.log.GetLog())
	log.Println("错误", l.log.GetErr())
}
