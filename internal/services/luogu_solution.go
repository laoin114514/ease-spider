package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/internal/repository"
	"spider/internal/utils"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

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
		s.AddErr(fmt.Sprintf("获取题解ID失败 %s", err.Error()))
		s.debug.Debug(fmt.Sprintf("获取题解ID失败 %s", err.Error()))
		return err
	}
	//不采取并发，因为题解数量较少，且获取题解时间较长
	for _, problemId := range problemIds {
		solutions, err := s.GetSolutionListByProblemId(problemId)
		s.debug.Debug(fmt.Sprintf("获取题解 题目ID：%s 目前题解数量：%d", problemId, len(solutions)))
		if err != nil {
			s.AddErr(fmt.Sprintf("获取题解失败 %s", err.Error()))
			s.debug.Debug(fmt.Sprintf("获取题解失败 %s", err.Error()))
			continue
		}
		err = s.storeSolution(solutions)
		if err != nil {
			s.AddErr(fmt.Sprintf("存储题解失败 %s", err.Error()))
			s.debug.Debug(fmt.Sprintf("存储题解失败 %s", err.Error()))
			continue
		}
	}
	s.debug.Debug(fmt.Sprintf("获取题解完成 %d", s.count))
	s.AddLog(fmt.Sprintf("获取题解完成 %d", s.count))
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
