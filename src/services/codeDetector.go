package services

import (
	"fmt"
	"math"
	"sort"
	"spider/src/utils/detector"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// CodeDetectionService 代码检测服务
type CodeDetectionService struct {
	similarityThreshold float64
}

// NewCodeDetectionService 创建代码检测服务
func NewCodeDetectionService(threshold float64) *CodeDetectionService {
	return &CodeDetectionService{
		similarityThreshold: threshold,
	}
}

// CompareCodeSimilarity 比较两个代码的相似度
func (s *CodeDetectionService) CompareCodeSimilarity(code1, code2, language string) float64 {
	return s.compareCodeSimilarity(code1, code2, language)
}

// CompareCodeSimilarityDetailed 比较两个代码的详细相似度信息
func (s *CodeDetectionService) CompareCodeSimilarityDetailed(code1, code2, language string) *SimilarityDetails {
	return s.compareCodeSimilarityDetailed(code1, code2, language)
}

// QuickCompare 快速比较函数，返回相似度百分比字符串
func (s *CodeDetectionService) QuickCompare(code1, code2, language string) string {
	return s.quickCompare(code1, code2, language)
}

// BatchCompare 批量比较代码相似度
func (s *CodeDetectionService) BatchCompare(targetCode, targetLanguage string, codeList []string, languageList []string) []SimilarityResult {
	return s.batchCompare(targetCode, targetLanguage, codeList, languageList)
}

// GetTopSimilar 获取最相似的代码
func (s *CodeDetectionService) GetTopNSimilar(targetCode, targetLanguage string, codeList []string, languageList []string, topN int) []SimilarityResult {
	return s.getTopNSimilar(targetCode, targetLanguage, codeList, languageList, topN)
}
func (s *CodeDetectionService) GetTopSimilar(targetCode, targetLanguage string, codeList []string, languageList []string) SimilarityResult {
	result := s.getTopNSimilar(targetCode, targetLanguage, codeList, languageList, 1)
	if len(result) > 0 {
		return result[0]
	}
	return SimilarityResult{}
}

// CompareCodeSimilarity 比较两个代码的相似度
// 参数:
//   - subCode: 提交的代码
//   - similarCode: 要比较的代码
//   - language: 编程语言 (go, java, python, javascript, cpp)
//
// 返回:
//   - rate: 相似度 (0.0-1.0, 1.0表示完全相同)
func (s *CodeDetectionService) compareCodeSimilarity(subCode, similarCode, language string) float64 {
	// 创建代码分析器
	analyzer1 := detector.NewCodeAnalyzer(language)
	analyzer2 := detector.NewCodeAnalyzer(language)

	// 分析两个代码的特征
	features1 := analyzer1.Analyze(subCode)
	features2 := analyzer2.Analyze(similarCode)

	// 创建相似度计算器
	similarityCalc := detector.NewSimilarityCalculator(0.0) // 阈值为0，返回所有相似度

	// 计算相似度
	result := similarityCalc.CalculateSimilarity(features1, features2)

	return result.OverallSimilarity
}

// CompareCodeSimilarityDetailed 比较两个代码的详细相似度信息
// 参数:
//   - subCode: 提交的代码
//   - similarCode: 要比较的代码
//   - language: 编程语言
//
// 返回:
//   - SimilarityDetails: 详细的相似度信息
func (s *CodeDetectionService) compareCodeSimilarityDetailed(subCode, similarCode, language string) *SimilarityDetails {
	// 创建代码分析器
	analyzer1 := detector.NewCodeAnalyzer(language)
	analyzer2 := detector.NewCodeAnalyzer(language)

	// 分析两个代码的特征
	features1 := analyzer1.Analyze(subCode)
	features2 := analyzer2.Analyze(similarCode)

	// 创建相似度计算器
	similarityCalc := detector.NewSimilarityCalculator(0.9)

	// 计算相似度
	result := similarityCalc.CalculateSimilarity(features1, features2)

	return &SimilarityDetails{
		OverallSimilarity:   result.OverallSimilarity,
		TextSimilarity:      result.TextSimilarity,
		SyntaxSimilarity:    result.SyntaxSimilarity,
		SemanticSimilarity:  result.SemanticSimilarity,
		StructureSimilarity: result.StructureSimilarity,
		IsSuspicious:        result.IsSuspicious,
		SimilarityType:      s.determineSimilarityType(result),
	}
}

// SimilarityDetails 详细相似度信息
type SimilarityDetails struct {
	OverallSimilarity   float64 `json:"overall_similarity"`   // 综合相似度
	TextSimilarity      float64 `json:"text_similarity"`      // 文本相似度
	SyntaxSimilarity    float64 `json:"syntax_similarity"`    // 语法相似度
	SemanticSimilarity  float64 `json:"semantic_similarity"`  // 语义相似度
	StructureSimilarity float64 `json:"structure_similarity"` // 结构相似度
	IsSuspicious        bool    `json:"is_suspicious"`        // 是否可疑
	SimilarityType      string  `json:"similarity_type"`      // 相似度类型
}

// determineSimilarityType 确定相似度类型
func (s *CodeDetectionService) determineSimilarityType(result *detector.SimilarityResult) string {
	if result.TextSimilarity > 0.9 {
		return "文本"
	} else if result.SyntaxSimilarity > 0.8 {
		return "语法"
	} else if result.SemanticSimilarity > 0.7 {
		return "语义"
	} else if result.StructureSimilarity > 0.6 {
		return "结构"
	}
	return "混合"
}

// QuickCompare 快速比较函数，返回相似度百分比字符串
// 参数:
//   - subCode: 提交的代码
//   - similarCode: 要比较的代码
//   - language: 编程语言
//
// 返回:
//   - string: 相似度百分比 (如 "85.6%")
func (s *CodeDetectionService) quickCompare(subCode, similarCode, language string) string {
	rate := s.compareCodeSimilarity(subCode, similarCode, language)
	return s.formatSimilarity(rate)
}

// FormatSimilarity 格式化相似度为百分比字符串
func (s *CodeDetectionService) formatSimilarity(similarity float64) string {
	return fmt.Sprintf("%.2f%%", similarity*100)
}

// SimpleTextSimilarity 简单的文本相似度计算（基于编辑距离）
// 参数:
//   - text1: 第一个文本
//   - text2: 第二个文本
//
// 返回:
//   - float64: 相似度 (0.0-1.0)
func (s *CodeDetectionService) simpleTextSimilarity(text1, text2 string) float64 {
	if text1 == text2 {
		return 1.0
	}

	if len(text1) == 0 || len(text2) == 0 {
		return 0.0
	}

	// 使用编辑距离算法
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	levenshtein := dmp.DiffLevenshtein(diffs)

	maxLen := math.Max(float64(len(text1)), float64(len(text2)))
	return 1.0 - (float64(levenshtein) / maxLen)
}

// SimpleTokenSimilarity 简单的token相似度计算（基于Jaccard相似度）
// 参数:
//   - tokens1: 第一个token列表
//   - tokens2: 第二个token列表
//
// 返回:
//   - float64: 相似度 (0.0-1.0)
func (s *CodeDetectionService) simpleTokenSimilarity(tokens1, tokens2 []string) float64 {
	if len(tokens1) == 0 && len(tokens2) == 0 {
		return 1.0
	}

	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	// 转换为map
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, token := range tokens1 {
		set1[strings.ToLower(token)] = true
	}

	for _, token := range tokens2 {
		set2[strings.ToLower(token)] = true
	}

	// 计算交集和并集
	intersection := 0
	for token := range set1 {
		if set2[token] {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// BatchCompare 批量比较代码相似度
// 参数:
//   - targetCode: 目标代码
//   - targetLanguage: 目标代码语言
//   - codeList: 要比较的代码列表
//   - languageList: 对应的语言列表
//
// 返回:
//   - []SimilarityResult: 相似度结果列表
func (s *CodeDetectionService) batchCompare(targetCode, targetLanguage string, codeList []string, languageList []string) []SimilarityResult {
	if len(codeList) != len(languageList) {
		return nil
	}

	var results []SimilarityResult
	targetAnalyzer := detector.NewCodeAnalyzer(targetLanguage)
	targetFeatures := targetAnalyzer.Analyze(targetCode)

	for i, code := range codeList {
		analyzer := detector.NewCodeAnalyzer(languageList[i])
		features := analyzer.Analyze(code)

		similarityCalc := detector.NewSimilarityCalculator(0.0)
		result := similarityCalc.CalculateSimilarity(targetFeatures, features)

		results = append(results, SimilarityResult{
			Index:               i,
			Code:                code,
			Language:            languageList[i],
			OverallSimilarity:   result.OverallSimilarity,
			TextSimilarity:      result.TextSimilarity,
			SyntaxSimilarity:    result.SyntaxSimilarity,
			SemanticSimilarity:  result.SemanticSimilarity,
			StructureSimilarity: result.StructureSimilarity,
			IsSuspicious:        result.IsSuspicious,
		})
	}

	// 按相似度排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].OverallSimilarity > results[j].OverallSimilarity
	})

	return results
}

// SimilarityResult 相似度结果
type SimilarityResult struct {
	Index               int     `json:"index"`                // 在列表中的索引
	Code                string  `json:"code"`                 // 代码内容
	Language            string  `json:"language"`             // 编程语言
	OverallSimilarity   float64 `json:"overall_similarity"`   // 综合相似度
	TextSimilarity      float64 `json:"text_similarity"`      // 文本相似度
	SyntaxSimilarity    float64 `json:"syntax_similarity"`    // 语法相似度
	SemanticSimilarity  float64 `json:"semantic_similarity"`  // 语义相似度
	StructureSimilarity float64 `json:"structure_similarity"` // 结构相似度
	IsSuspicious        bool    `json:"is_suspicious"`        // 是否可疑
}

// GetTopSimilar 获取最相似的代码
// 参数:
//   - targetCode: 目标代码
//   - targetLanguage: 目标代码语言
//   - codeList: 要比较的代码列表
//   - languageList: 对应的语言列表
//   - topN: 返回前N个结果
//
// 返回:
//   - []SimilarityResult: 前N个最相似的结果
func (s *CodeDetectionService) getTopNSimilar(targetCode, targetLanguage string, codeList []string, languageList []string, topN int) []SimilarityResult {
	results := s.batchCompare(targetCode, targetLanguage, codeList, languageList)

	if topN <= 0 || topN > len(results) {
		topN = len(results)
	}

	return results[:topN]
}
