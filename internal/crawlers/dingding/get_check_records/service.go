package getcheckrecords

import (
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/utils"
	"time"
)

// 获取钉钉打卡数据
func (g *GetCheckRecords) GetDingdingCheckUpData() error {
	g.insertCount = 0
	// 获取钉钉token
	var err error
	g.token, err = g.getDingdingToken()
	if err != nil {
		return err
	}
	dingUserMap, err := g.repo.GetDingUserMap()
	if err != nil {
		g.log.Errorf("获取钉钉用户映射失败 %s", err.Error())
		return err
	}
	g.dingUserMap = dingUserMap
	// 获取钉钉打卡数据（周）
	err = g.getDataWithWeek(config.AppConfig.Dingding.WeekRange)
	if err != nil {
		return err
	}
	g.log.Printf("插入钉钉打卡数据：%d", g.insertCount)
	return nil
}

// 获取钉钉token
func (g *GetCheckRecords) getDingdingToken() (string, error) {

	appKey := config.AppConfig.Dingding.AppKey
	appSecret := config.AppConfig.Dingding.AppSecret

	req := utils.NewRequest[TokenResponse](true)

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
func (g *GetCheckRecords) getDataWithWeek(week int) error {
	dateFormat := utils.NewDateFormat()
	req := utils.NewRequest[CheckUpData](true)

	for i := int64(0); i < int64(week); i++ {
		from := dateFormat.BeforDateTimeWithDay(i + 7)
		to := dateFormat.BeforDateTimeWithDay(i)
		userIds := []string{}
		for k, _ := range g.dingUserMap {
			userIds = append(userIds, k)
		}
		body := map[string]any{
			"userIds":       userIds,
			"checkDateFrom": from,
			"checkDateTo":   to,
		}
		dingdingCheckUpData, err := req.Post("https://oapi.dingtalk.com/attendance/listRecord?access_token="+g.token, body)
		if err != nil {
			return err
		}
		for _, v := range dingdingCheckUpData.Records {
			table := g.buildDingdingCheckUpTable(v)
			err = db.Insert_checkup(table)
			if err != nil {
				continue
			}
			g.log.Printf("插入钉钉打卡数据：%v", table)
			g.insertCount++
		}
	}
	return nil
}

// 构建钉钉打卡表格
func (g *GetCheckRecords) buildDingdingCheckUpTable(checkUpData CheckRecord) db.Ding_checkUp {
	var table db.Ding_checkUp
	table.Ding_id = checkUpData.UserId
	table.Time = time.UnixMilli(checkUpData.UserCheckTime).Add(constants.TimeZoneOffsetHours * time.Hour)
	table.Name = g.dingUserMap[checkUpData.UserId].(string)
	table.Check_type = checkUpData.CheckType
	return table
}
func (g *GetCheckRecords) Clear() error {
	g.insertCount = 0
	return nil
}
