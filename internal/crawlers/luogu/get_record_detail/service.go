package getrecorddetail

import (
	"errors"
	"fmt"
	"spider/config"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/internal/repository"
	"spider/internal/utils"
	"strconv"
)

// 获取源代码
func (g *GetRecordDetail) GetAndStoreSourceCode() error {
	cookie := utils.JsonDB.Get("Cookie")
	if cookie == nil {
		return errors.New("cookie不存在")
	}
	g.cookie = cookie.(string)
	// 注意: 这里应该通过依赖注入获取 LuoguTeam，而不是直接创建
	// 当前设计下，LuoguTeam 应该在 Luogu 结构体中统一管理
	luoguTeam := NewLuoguTeam()
	err := luoguTeam.GetMembers()
	if err != nil {
		g.log.Errorf("获取洛谷团队成员失败 %s", err.Error())
		return err
	}

	subids, err := g.repo.GetSubidNoSourceCode(luoguTeam.results)
	if err != nil {
		g.log.Errorf("获取提交记录ID失败 %s", err.Error())
		return err
	}
	g.log.Printf("总共 %d个提交记录需要获取源代码\n团队成员：%v", len(subids), luoguTeam.results)

	// 创建并发器
	conCurrenter := utils.NewConCurrenter[string](config.AppConfig.Luogu.LuoguSubmissionDetailConcurrency)

	// 并发获取提交记录源代码
	err = conCurrenter.Run(subids, func(subid string) error {
		name, err := g.repo.GetNameBySubid(subid)
		if err != nil {
			g.log.Errorf("获取提交记录名称失败 %s", err.Error())
			return err
		}

		respJson, err := g.getSourceCodeBySubid(subid)
		if err != nil {
			g.log.Errorf("解析提交记录源代码失败 %s", err.Error())
			return err
		}

		//源代码长度小于5，则填充无
		if len(respJson.CurrentData.Record.SourceCode) <= 5 {
			g.log.Printf("提交记录源代码为空 %s", subid)
			respJson.CurrentData.Record.SourceCode = "无"
			g.log.Printf("%s源代码为空", subid+name)
			if respJson.CurrentData.Record.Problem.Type == "P" && respJson.CurrentData.Record.Status == constants.LuoguStatusAccepted {
				g.log.Printf("%sAC且无源代码，跳过", subid+name)
				return nil
			}
		}

		err = g.repo.InsertSourceCode(subid, respJson.CurrentData.Record.SourceCode)
		if err != nil {
			g.log.Errorf("插入提交记录源代码失败 %s", err.Error())

			return err
		}

		g.log.Printf("%s插入成功", subid+name)
		g.count++
		return nil
	})
	g.log.Printf("插入提交记录源代码完成 %d", g.count)
	return err
}

func (g *GetRecordDetail) getSourceCodeBySubid(subid string) (models.LuoguSubmissionDetailResponse, error) {
	req := utils.NewRequest[models.LuoguSubmissionDetailResponse](false)
	url := fmt.Sprintf("https://www.luogu.com.cn/record/%s?_contentOnly=1", subid)
	req.SetCookie(g.cookie)
	resp, err := req.Get(url, map[string]string{})
	if err != nil {
		return models.LuoguSubmissionDetailResponse{}, err
	}
	if resp.Code != 200 {
		return models.LuoguSubmissionDetailResponse{}, errors.New("获取提交记录源代码失败" + strconv.Itoa(resp.Code))
	}
	return resp, nil
}

type LuoguTeam struct {
	repo    *repository.LuoguRepository
	debug   *utils.Debug
	results map[string]bool
	count   int
}

func NewLuoguTeam() *LuoguTeam {
	return &LuoguTeam{
		repo:    repository.NewLuoguRepository(),
		debug:   utils.NewDebug(config.AppConfig.DebugConfig.All),
		count:   0,
		results: make(map[string]bool),
	}
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
	return nil
}
