package services

import (
	"fmt"
	"spider/config"
	"spider/internal/models"
	"spider/internal/repository"
	"spider/internal/utils"
	"strconv"
)

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
