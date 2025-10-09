package services

import (
	"spider/src/utils"
)

type CheatJudge struct {
	*LogService
	suspiciousResults []suspiciousResult
}
type suspiciousResult struct {
	ProblemID         string
	SimilarityDetails SimilarityResult
	Code              string
	SimilarCode       string
}

func NewCheatJudge() *CheatJudge {
	return &CheatJudge{
		LogService: NewLogService("logs/cheatJudge.log", "logs/cheatJudge.err.log"),
	}
}
func (c *CheatJudge) GetSuspiciousResults() []suspiciousResult {
	return c.suspiciousResults
}
func (c *CheatJudge) JudgeWithCode(code string, problemID string) error {
	luoguSolutionService := NewLuoguSolution()
	solutions, err := luoguSolutionService.GetSolutionList(problemID)
	if err != nil {
		return err
	}

	for _, solution := range solutions {
		codeDetectorService := NewCodeDetectionService(1)
		markdownParser := utils.NewMarkdownParser(solution.Content)

		//提取代码块
		codeBlocks := markdownParser.ExtractCodeBlocks()

		//将codeBlocks转换为codeBlocksStr和languageList
		codeBlocksStr := make([]string, len(codeBlocks))
		languageList := make([]string, len(codeBlocks))
		for i, codeBlock := range codeBlocks {
			codeBlocksStr[i] = codeBlock.Content
			languageList[i] = codeBlock.Language
		}
		similarityResults := codeDetectorService.GetTopSimilar(code, "markdown", codeBlocksStr, languageList)
		if similarityResults.IsSuspicious {
			c.suspiciousResults = append(c.suspiciousResults, suspiciousResult{
				ProblemID:         problemID,
				SimilarityDetails: similarityResults,
				Code:              code,
				SimilarCode:       similarityResults.Code,
			})
		}
	}
	return nil
}
