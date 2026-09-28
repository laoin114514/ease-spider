package getofficialproblem

import (
	"fmt"
	"spider/config/db"
	"strings"

	cf "github.com/laoin114514/codeforcesClient"
)

// GetCfOfficialProblems 获取CF官方题目 - 对外提供的主要接口
func (g *GetOfifcialProblems) GetCfOfficialProblems() error {
	g.count = 0

	// problemset.problems 是公共接口
	resp, err := g.client.ProblemsetProblems(&cf.ProblemsetProblemsParams{})
	if err != nil {
		g.log.Errorf("获取数据失败 %s", err.Error())
		return err
	}

	// 处理官方题目数据
	err = g.processOfficialProblems(resp)
	if err != nil {
		g.log.Errorf("处理数据失败 %s", err.Error())
		return err
	}

	g.log.Printf("插入官方题目 %d", g.count)
	return nil
}

// processOfficialProblems 处理官方题目数据 - 私有方法
func (g *GetOfifcialProblems) processOfficialProblems(resp *cf.ProblemsetProblemsResponse) error {
	if resp.Result == nil {
		return nil
	}
	for _, cfProblem := range resp.Result.Problems {
		if cfProblem == nil {
			continue
		}
		table := g.buildProblemTable(cfProblem)
		err := db.Insert_cf_official(table)
		if err != nil {
			continue
		}
		g.count++
	}
	return nil
}

// buildProblemTable 构建题目表 - 私有方法
//
// 标签的拼装方式与迁移前完全一致（包括 Tags 为空时会得到 [""] 这一既有行为）。
func (g *GetOfifcialProblems) buildProblemTable(cfProblem *cf.Problem) db.Cf_official_problems {
	tags := strings.Join(cfProblem.Tags, "\",\"")
	tagsArr := fmt.Sprintf("[\"%s\"]", tags)
	rating := cfProblem.Rating
	if rating == 0 {
		rating = -1
	}
	return db.Cf_official_problems{
		Problem_id: fmt.Sprintf("%d%s", cfProblem.ContestID, cfProblem.Index),
		Title:      cfProblem.Name,
		Points:     int(cfProblem.Points),
		Rating:     rating,
		Tags:       tagsArr,
	}
}
