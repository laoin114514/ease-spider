package getofficialproblem

import (
	"fmt"
	"spider/config/db"
	"spider/internal/models"
	"spider/internal/utils"
	cfurlgenerator "spider/pkg/cf-url-generator"
	"strings"
)

// GetCfOfficialProblems 获取CF官方题目 - 对外提供的主要接口
func (g *GetOfifcialProblems) GetCfOfficialProblems() error {
	g.count = 0

	// 构建请求URL
	url, err := g.buildProblemsURL()
	if err != nil {
		g.log.Errorf("构建请求URL失败 %s", err.Error())
		return err
	}

	// 获取数据
	resp, err := g.fetchProblemsData(url)
	if err != nil {
		g.log.Errorf("获取数据失败 %s", err.Error())
		return err
	}

	// 处理官方题目数据
	err = g.processOfficialProblems(&resp)
	if err != nil {
		g.log.Errorf("处理数据失败 %s", err.Error())
		return err
	}

	g.log.Printf("插入官方题目 %d", g.count)
	return nil
}

// buildProblemsURL 构建题目请求URL - 私有方法
func (g *GetOfifcialProblems) buildProblemsURL() (string, error) {
	return g.urlGenerator.ProblemSet.Problems(&cfurlgenerator.ProblemsetProblemsParams{})
}

// fetchProblemsData 获取题目数据 - 私有方法
func (g *GetOfifcialProblems) fetchProblemsData(url string) (models.CfOfficialProblemsResponse, error) {
	req := utils.NewRequest[models.CfOfficialProblemsResponse](true)
	return req.Get(url, map[string]string{})
}

// processOfficialProblems 处理官方题目数据 - 私有方法
func (g *GetOfifcialProblems) processOfficialProblems(resp *models.CfOfficialProblemsResponse) error {
	for _, cfProblem := range resp.Result.Problems {
		table := g.buildProblemTable(&cfProblem)
		err := db.Insert_cf_official(table)
		if err != nil {
			continue
		}
		g.count++
	}
	return nil
}

// buildProblemTable 构建题目表 - 私有方法
func (g *GetOfifcialProblems) buildProblemTable(cfProblem *models.CfProblem) db.Cf_official_problems {
	tags := strings.Join(cfProblem.Tags, "\",\"")
	tagsArr := fmt.Sprintf("[\"%s\"]", tags)
	if cfProblem.Rating == 0 {
		cfProblem.Rating = -1
	}
	return db.Cf_official_problems{
		Problem_id: fmt.Sprintf("%d%s", cfProblem.ContestId, cfProblem.Index),
		Title:      cfProblem.Name,
		Points:     int(cfProblem.Points),
		Rating:     int(cfProblem.Rating),
		Tags:       tagsArr,
	}
}
