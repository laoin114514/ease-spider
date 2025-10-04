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
	conCurrenter := utils.NewConCurrenter[map[string]string](concurrency)
	nameMaps, err := l.repo.GetUserNameMap()
	if err != nil {
		return err
	}

	//通过并发器来获取洛谷用户提交记录
	conCurrenter.Run(nameMaps, func(userNameMap map[string]string) {
		uid := userNameMap["luogu_uid"]
		realName := userNameMap["real_name"]

		//获取初始化数据：总数和每页数量
		initData, err := l.req.Get("https://www.luogu.com.cn/record/list", map[string]string{"user": uid, "page": "1", "_contentOnly": "1"})
		if err != nil {
			l.log.AddErr(fmt.Sprintf("%s获取提交记录失败 %s", realName, err.Error()))
			return
		}
		// 检查数据是否为空
		if initData.CurrentData.Records.Count == 0 {
			l.log.AddErr(fmt.Sprintf("%s没有提交记录", realName))
			return
		}
		if initData.Code != 200 {
			l.log.AddErr(fmt.Sprintf("%s获取提交记录失败 %s", realName, initData.Code))
			return
		}

		page := int(math.Ceil(float64(initData.CurrentData.Records.Count) / float64(initData.CurrentData.Records.PerPage)))

		//循环获取所有页数的数据
		for i := 1; i <= page; i++ {
			data, err := l.req.Get(
				"https://www.luogu.com.cn/record/list",
				map[string]string{
					"user":         uid,
					"page":         strconv.Itoa(i),
					"_contentOnly": "1",
				},
			)
			if err != nil {
				l.log.AddErr(fmt.Sprintf("%s获取第%d页提交记录失败 %s", realName, i, err.Error()))
				return
			}

			//处理数据
			err = l.handleRecords(&data)
			if err != nil {
				l.log.AddErr(fmt.Sprintf("%s处理第%d页提交记录失败 %s", realName, i, err.Error()))
				return
			}
		}
	})
	return nil
}

// 处理提交记录
func (l *Luogu) handleRecords(records *models.LuoguRecordsResponse) error {
	var tables []db.Luogu_all_submissions
	for _, record := range records.CurrentData.Records.Result {
		table := l.buildTable(record)
		tables = append(tables, table)
	}

	return db.BatchInsert_luogu_sub(tables)
}

// 构建提交记录表
func (l *Luogu) buildTable(record models.UserRecord) db.Luogu_all_submissions {
	return db.Luogu_all_submissions{
		Sub_id:        fmt.Sprintf("%d", record.ID),
		Uid:           fmt.Sprintf("%d", record.User.UID),
		Problem_id:    record.Problem.PID,
		Problem_name:  record.Problem.Title,
		Difficulty:    fmt.Sprintf("%d", record.Problem.Difficulty),
		Is_pass:       record.Status == 12,
		Creation_time: time.Unix(record.SubmitTime, 0).Add(8 * time.Hour),
	}
}

// 打印日志
func (l *Luogu) Log() {
	log.Println("日志", l.log.GetLog())
	log.Println("错误", l.log.GetErr())
}
