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
	SimilarityDetails *SimilarityDetails
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
func (c *CheatJudge) JudgeWithCode(code string, problemID string, threshold float64) error {
	luoguSolutionService := NewLuoguSolution()
	solutions, err := luoguSolutionService.GetSolutionList(problemID)
	if err != nil {
		return err
	}

	for _, solution := range solutions {
		codeDetectorService := NewCodeDetectionService(threshold)
		markdownParser := utils.NewMarkdownParser(solution.Content)

		//提取代码块
		codeBlocks := markdownParser.ExtractCodeBlocks()
		for _, codeBlock := range codeBlocks {
			similarityDetails := codeDetectorService.CompareCodeSimilarityDetailed(code, codeBlock.Content, codeBlock.Language)
			if similarityDetails.IsSuspicious {
				c.suspiciousResults = append(c.suspiciousResults, suspiciousResult{
					ProblemID:         problemID,
					SimilarityDetails: similarityDetails,
					Code:              code,
					SimilarCode:       codeBlock.Content,
				})
			}
		}
	}
	return nil
}
