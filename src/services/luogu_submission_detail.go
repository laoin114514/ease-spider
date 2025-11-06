package services

import (
	"errors"
	"fmt"
	"spider/config"
	"spider/src/constants"
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
	"strconv"
)

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
	// 注意: 这里应该通过依赖注入获取 LuoguTeam，而不是直接创建
	// 当前设计下，LuoguTeam 应该在 Luogu 结构体中统一管理
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
