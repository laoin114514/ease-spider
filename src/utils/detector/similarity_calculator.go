package detector

import (
	"math"
	"sort"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// SimilarityCalculator 相似度计算器
type SimilarityCalculator struct {
	threshold float64 // 相似度阈值
}

// NewSimilarityCalculator 创建相似度计算器
func NewSimilarityCalculator(threshold float64) *SimilarityCalculator {
	return &SimilarityCalculator{
		threshold: threshold,
	}
}

// CalculateSimilarity 计算两个代码特征的相似度
func (sc *SimilarityCalculator) CalculateSimilarity(features1, features2 *CodeFeatures) *SimilarityResult {
	result := &SimilarityResult{
		OverallSimilarity:   0.0,
		TextSimilarity:      0.0,
		SyntaxSimilarity:    0.0,
		SemanticSimilarity:  0.0,
		StructureSimilarity: 0.0,
	}

	// 1. 文本相似度（基于编辑距离）
	result.TextSimilarity = sc.calculateTextSimilarity(features1, features2)

	// 2. 语法相似度（基于token比较）
	result.SyntaxSimilarity = sc.calculateSyntaxSimilarity(features1, features2)

	// 3. 语义相似度（基于关键词和指纹）
	result.SemanticSimilarity = sc.calculateSemanticSimilarity(features1, features2)

	// 4. 结构相似度（基于代码结构）
	result.StructureSimilarity = sc.calculateStructureSimilarity(features1, features2)

	// 5. 综合相似度（加权平均）
	result.OverallSimilarity = sc.calculateOverallSimilarity(result)

	// 6. 判断是否可疑
	result.IsSuspicious = result.OverallSimilarity >= sc.threshold

	return result
}

// SimilarityResult 相似度结果
type SimilarityResult struct {
	OverallSimilarity   float64 `json:"overall_similarity"`
	TextSimilarity      float64 `json:"text_similarity"`
	SyntaxSimilarity    float64 `json:"syntax_similarity"`
	SemanticSimilarity  float64 `json:"semantic_similarity"`
	StructureSimilarity float64 `json:"structure_similarity"`
	IsSuspicious        bool    `json:"is_suspicious"`
}

// calculateTextSimilarity 计算文本相似度
func (sc *SimilarityCalculator) calculateTextSimilarity(features1, features2 *CodeFeatures) float64 {
	// 使用编辑距离算法
	dmp := diffmatchpatch.New()

	// 标准化代码
	text1 := sc.normalizeText(features1.Tokens)
	text2 := sc.normalizeText(features2.Tokens)

	// 计算编辑距离
	diffs := dmp.DiffMain(text1, text2, false)

	// 计算相似度
	similarity := dmp.DiffLevenshtein(diffs)
	maxLen := math.Max(float64(len(text1)), float64(len(text2)))

	if maxLen == 0 {
		return 1.0
	}

	return 1.0 - (float64(similarity) / maxLen)
}

// calculateSyntaxSimilarity 计算语法相似度
func (sc *SimilarityCalculator) calculateSyntaxSimilarity(features1, features2 *CodeFeatures) float64 {
	tokens1 := features1.Tokens
	tokens2 := features2.Tokens

	if len(tokens1) == 0 && len(tokens2) == 0 {
		return 1.0
	}

	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	// 使用Jaccard相似度
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, token := range tokens1 {
		set1[token] = true
	}

	for _, token := range tokens2 {
		set2[token] = true
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

// calculateSemanticSimilarity 计算语义相似度
func (sc *SimilarityCalculator) calculateSemanticSimilarity(features1, features2 *CodeFeatures) float64 {
	// 基于指纹相似度
	fingerprint1 := features1.Fingerprint
	fingerprint2 := features2.Fingerprint

	if fingerprint1 == fingerprint2 {
		return 1.0
	}

	// 基于关键词相似度
	keywords1 := features1.Keywords
	keywords2 := features2.Keywords

	keywordSimilarity := sc.calculateSetSimilarity(keywords1, keywords2)

	// 基于指纹的编辑距离
	fingerprintSimilarity := sc.calculateStringSimilarity(fingerprint1, fingerprint2)

	// 加权平均
	return 0.6*fingerprintSimilarity + 0.4*keywordSimilarity
}

// calculateStructureSimilarity 计算结构相似度
func (sc *SimilarityCalculator) calculateStructureSimilarity(features1, features2 *CodeFeatures) float64 {
	struct1 := features1.Structure
	struct2 := features2.Structure

	if struct1 == nil || struct2 == nil {
		return 0.0
	}

	// 函数数量相似度
	funcSimilarity := sc.calculateNumericSimilarity(struct1.FunctionCount, struct2.FunctionCount)

	// 类数量相似度
	classSimilarity := sc.calculateNumericSimilarity(struct1.ClassCount, struct2.ClassCount)

	// 复杂度相似度
	complexitySimilarity := sc.calculateNumericSimilarity(struct1.Complexity, struct2.Complexity)

	// 导入相似度
	importSimilarity := sc.calculateSetSimilarity(struct1.Imports, struct2.Imports)

	// 加权平均
	return 0.3*funcSimilarity + 0.2*classSimilarity + 0.3*complexitySimilarity + 0.2*importSimilarity
}

// calculateOverallSimilarity 计算综合相似度
func (sc *SimilarityCalculator) calculateOverallSimilarity(result *SimilarityResult) float64 {
	// 加权平均，可以根据需要调整权重
	weights := map[string]float64{
		"text":      0.25,
		"syntax":    0.30,
		"semantic":  0.30,
		"structure": 0.15,
	}

	overall := weights["text"]*result.TextSimilarity +
		weights["syntax"]*result.SyntaxSimilarity +
		weights["semantic"]*result.SemanticSimilarity +
		weights["structure"]*result.StructureSimilarity

	return math.Min(overall, 1.0)
}

// calculateSetSimilarity 计算集合相似度
func (sc *SimilarityCalculator) calculateSetSimilarity(set1, set2 []string) float64 {
	if len(set1) == 0 && len(set2) == 0 {
		return 1.0
	}

	if len(set1) == 0 || len(set2) == 0 {
		return 0.0
	}

	// 转换为map
	map1 := make(map[string]bool)
	map2 := make(map[string]bool)

	for _, item := range set1 {
		map1[strings.ToLower(item)] = true
	}

	for _, item := range set2 {
		map2[strings.ToLower(item)] = true
	}

	// 计算交集和并集
	intersection := 0
	for item := range map1 {
		if map2[item] {
			intersection++
		}
	}

	union := len(map1) + len(map2) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// calculateNumericSimilarity 计算数值相似度
func (sc *SimilarityCalculator) calculateNumericSimilarity(val1, val2 int) float64 {
	if val1 == 0 && val2 == 0 {
		return 1.0
	}

	if val1 == 0 || val2 == 0 {
		return 0.0
	}

	maxVal := math.Max(float64(val1), float64(val2))
	minVal := math.Min(float64(val1), float64(val2))

	return minVal / maxVal
}

// calculateStringSimilarity 计算字符串相似度
func (sc *SimilarityCalculator) calculateStringSimilarity(str1, str2 string) float64 {
	if str1 == str2 {
		return 1.0
	}

	if len(str1) == 0 || len(str2) == 0 {
		return 0.0
	}

	// 使用编辑距离
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(str1, str2, false)
	levenshtein := dmp.DiffLevenshtein(diffs)

	maxLen := math.Max(float64(len(str1)), float64(len(str2)))
	return 1.0 - (float64(levenshtein) / maxLen)
}

// normalizeText 标准化文本
func (sc *SimilarityCalculator) normalizeText(tokens []string) string {
	// 排序并连接
	sort.Strings(tokens)
	return strings.Join(tokens, " ")
}

// FindSimilarSubmissions 查找相似的代码提交
func (sc *SimilarityCalculator) FindSimilarSubmissions(targetFeatures *CodeFeatures, allFeatures []*CodeFeatures) []*SimilarityMatch {
	var matches []*SimilarityMatch

	for _, features := range allFeatures {
		if features.Hash == targetFeatures.Hash {
			continue // 跳过自己
		}

		similarity := sc.CalculateSimilarity(targetFeatures, features)

		if similarity.OverallSimilarity >= sc.threshold {
			matches = append(matches, &SimilarityMatch{
				Features:   features,
				Similarity: similarity,
			})
		}
	}

	// 按相似度排序
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Similarity.OverallSimilarity > matches[j].Similarity.OverallSimilarity
	})

	return matches
}

// SimilarityMatch 相似度匹配结果
type SimilarityMatch struct {
	Features   *CodeFeatures     `json:"features"`
	Similarity *SimilarityResult `json:"similarity"`
}
