package services

import (
	"fmt"
	"os"
	"spider/config/db"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"time"

	"github.com/joho/godotenv"
)

type Dingding struct {
	repo        *repository.DingdingRepository
	log         *utils.LogContainer
	token       string
	dingUserMap map[string]any
	insertCount int
}

// 创建钉钉服务
func NewDingdingService() *Dingding {
	JsonDB := utils.JsonDB
	return &Dingding{
		repo:        repository.NewDingdingRepository(),
		log:         utils.NewLogContainer(),
		token:       "",
		dingUserMap: JsonDB.Get("dingUserId").(map[string]any),
		insertCount: 0,
	}
}

// 获取钉钉打卡数据
func (d *Dingding) GetDingdingCheckUpData() error {
	d.insertCount = 0
	// 获取钉钉token
	var err error
	d.token, err = d.getDingdingToken()
	if err != nil {
		return err
	}

	// 获取钉钉打卡数据（周）
	err = d.getDataWithWeek(60)
	if err != nil {
		return err
	}
	d.log.AddLog(fmt.Sprintf("插入钉钉打卡数据：%d", d.insertCount))
	return nil
}

// 获取钉钉token
func (d *Dingding) getDingdingToken() (string, error) {
	godotenv.Load()

	appKey := os.Getenv("ding_accessToken")
	appSecret := os.Getenv("ding_appSecret")

	req := utils.NewRequest[models.DingdingTokenResponse]()

	url := "https://api.dingtalk.com/v1.0/oauth2/accessToken"
	token, err := req.Post(url,
		map[string]any{
			"appKey":    appKey,
			"appSecret": appSecret,
		})
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// 获取钉钉打卡数据（周）
func (d *Dingding) getDataWithWeek(week int) error {
	dateFormat := utils.NewDateFormat()
	req := utils.NewRequest[models.DingdingCheckUpData]()

	for i := int64(0); i < int64(week); i++ {
		from := dateFormat.BeforDateTimeWithDay(i + 7)
		to := dateFormat.BeforDateTimeWithDay(i)
		userIds := []string{}
		for k, _ := range d.dingUserMap {
			userIds = append(userIds, k)
		}
		body := map[string]any{
			"userIds":       userIds,
			"checkDateFrom": from,
			"checkDateTo":   to,
		}
		dingdingCheckUpData, err := req.Post("https://oapi.dingtalk.com/attendance/listRecord?access_token="+d.token, body)
		if err != nil {
			return err
		}
		for _, v := range dingdingCheckUpData.Records {
			table := d.buildDingdingCheckUpTable(v)
			err = db.Insert_checkup(table)
			if err != nil {
				continue
			}
			d.log.AddLog(fmt.Sprintf("插入钉钉打卡数据：%v", table))
			d.insertCount++
		}
	}
	return nil
}

// 构建钉钉打卡表格
func (d *Dingding) buildDingdingCheckUpTable(checkUpData models.DingdingCheckRecord) db.Ding_checkUp {
	var table db.Ding_checkUp
	table.Ding_id = checkUpData.UserId
	table.Time = time.UnixMilli(checkUpData.UserCheckTime).Add(8 * time.Hour)
	table.Name = d.dingUserMap[checkUpData.UserId].(string)
	table.Check_type = checkUpData.CheckType
	return table
}
func (d *Dingding) GetLog() []string {
	return d.log.GetLog()
}
func (d *Dingding) GetErr() []string {
	return d.log.GetErr()
}
func (d *Dingding) Clear() error {
	d.log.ClearLog()
	d.log.ClearErr()
	d.insertCount = 0
	return nil
}
