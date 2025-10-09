package detector

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/kljensen/snowball"
)

// CodeAnalyzer 代码分析器
type CodeAnalyzer struct {
	language string
}

// NewCodeAnalyzer 创建代码分析器
func NewCodeAnalyzer(language string) *CodeAnalyzer {
	return &CodeAnalyzer{
		language: language,
	}
}

// Analyze 分析代码并返回特征
func (ca *CodeAnalyzer) Analyze(code string) *CodeFeatures {
	features := &CodeFeatures{
		Language:    ca.language,
		Hash:        ca.generateHash(code),
		Fingerprint: ca.generateFingerprint(code),
		Tokens:      ca.tokenize(code),
		Keywords:    ca.extractKeywords(code),
		Structure:   ca.analyzeStructure(code),
	}
	return features
}

// CodeFeatures 代码特征
type CodeFeatures struct {
	Language    string         `json:"language"`
	Hash        string         `json:"hash"`
	Fingerprint string         `json:"fingerprint"`
	Tokens      []string       `json:"tokens"`
	Keywords    []string       `json:"keywords"`
	Structure   *CodeStructure `json:"structure"`
}

// CodeStructure 代码结构
type CodeStructure struct {
	FunctionCount int            `json:"function_count"`
	ClassCount    int            `json:"class_count"`
	VariableCount int            `json:"variable_count"`
	CommentCount  int            `json:"comment_count"`
	LineCount     int            `json:"line_count"`
	Complexity    int            `json:"complexity"`
	Imports       []string       `json:"imports"`
	Functions     []FunctionInfo `json:"functions"`
}

// FunctionInfo 函数信息
type FunctionInfo struct {
	Name   string `json:"name"`
	Line   int    `json:"line"`
	Params int    `json:"params"`
	Lines  int    `json:"lines"`
}

// generateHash 生成代码哈希
func (ca *CodeAnalyzer) generateHash(code string) string {
	// 标准化代码（移除空白字符和注释）
	normalized := ca.normalizeCode(code)
	hash := md5.Sum([]byte(normalized))
	return fmt.Sprintf("%x", hash)
}

// generateFingerprint 生成代码指纹
func (ca *CodeAnalyzer) generateFingerprint(code string) string {
	tokens := ca.tokenize(code)
	// 使用词干提取标准化关键词
	var stemmedTokens []string
	for _, token := range tokens {
		if stem, err := snowball.Stem(token, "english", true); err == nil {
			stemmedTokens = append(stemmedTokens, stem)
		} else {
			stemmedTokens = append(stemmedTokens, token)
		}
	}

	// 排序并去重
	sort.Strings(stemmedTokens)
	var uniqueTokens []string
	prev := ""
	for _, token := range stemmedTokens {
		if token != prev && len(token) > 2 {
			uniqueTokens = append(uniqueTokens, token)
			prev = token
		}
	}

	return strings.Join(uniqueTokens, "|")
}

// tokenize 代码分词
func (ca *CodeAnalyzer) tokenize(code string) []string {
	// 移除注释
	code = ca.removeComments(code)

	// 分词正则表达式
	re := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*|[0-9]+|[+\-*/=<>!&|^~%]+|[()[\]{}.,;:]`)
	matches := re.FindAllString(code, -1)

	var tokens []string
	for _, match := range matches {
		if len(match) > 1 && !ca.isKeyword(match) {
			tokens = append(tokens, strings.ToLower(match))
		}
	}

	return tokens
}

// extractKeywords 提取关键词
func (ca *CodeAnalyzer) extractKeywords(code string) []string {
	keywords := make(map[string]bool)

	// 根据语言提取关键词
	switch ca.language {
	case "go":
		keywords = ca.extractGoKeywords(code)
	case "java":
		keywords = ca.extractJavaKeywords(code)
	case "python":
		keywords = ca.extractPythonKeywords(code)
	case "javascript":
		keywords = ca.extractJavaScriptKeywords(code)
	case "cpp", "c++":
		keywords = ca.extractCppKeywords(code)
	}

	var result []string
	for keyword := range keywords {
		result = append(result, keyword)
	}

	return result
}

// analyzeStructure 分析代码结构
func (ca *CodeAnalyzer) analyzeStructure(code string) *CodeStructure {
	lines := strings.Split(code, "\n")
	structure := &CodeStructure{
		LineCount: len(lines),
		Imports:   ca.extractImports(code),
		Functions: ca.extractFunctions(code),
	}

	structure.FunctionCount = len(structure.Functions)
	structure.CommentCount = ca.countComments(code)
	structure.Complexity = ca.calculateComplexity(code)

	return structure
}

// normalizeCode 标准化代码
func (ca *CodeAnalyzer) normalizeCode(code string) string {
	// 移除注释
	code = ca.removeComments(code)

	// 移除多余空白
	re := regexp.MustCompile(`\s+`)
	code = re.ReplaceAllString(code, " ")

	// 移除字符串内容
	re = regexp.MustCompile(`"[^"]*"`)
	code = re.ReplaceAllString(code, `""`)

	re = regexp.MustCompile(`'[^']*'`)
	code = re.ReplaceAllString(code, `''`)

	return strings.TrimSpace(code)
}

// removeComments 移除注释
func (ca *CodeAnalyzer) removeComments(code string) string {
	lines := strings.Split(code, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "/*") {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// isKeyword 判断是否为关键字
func (ca *CodeAnalyzer) isKeyword(word string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "for": true, "while": true, "do": true,
		"switch": true, "case": true, "default": true, "break": true,
		"continue": true, "return": true, "function": true, "var": true,
		"let": true, "const": true, "class": true, "interface": true,
		"import": true, "package": true, "public": true, "private": true,
		"protected": true, "static": true, "final": true, "abstract": true,
	}

	return keywords[strings.ToLower(word)]
}

// 各语言关键词提取方法
func (ca *CodeAnalyzer) extractGoKeywords(code string) map[string]bool {
	keywords := make(map[string]bool)
	re := regexp.MustCompile(`\b(func|package|import|var|const|type|interface|struct|if|else|for|range|switch|case|default|return|go|defer|chan|select)\b`)
	matches := re.FindAllString(code, -1)
	for _, match := range matches {
		keywords[match] = true
	}
	return keywords
}

func (ca *CodeAnalyzer) extractJavaKeywords(code string) map[string]bool {
	keywords := make(map[string]bool)
	re := regexp.MustCompile(`\b(public|private|protected|class|interface|extends|implements|static|final|abstract|synchronized|volatile|transient|native|strictfp|if|else|for|while|do|switch|case|default|break|continue|return|try|catch|finally|throw|throws|new|this|super|package|import)\b`)
	matches := re.FindAllString(code, -1)
	for _, match := range matches {
		keywords[match] = true
	}
	return keywords
}

func (ca *CodeAnalyzer) extractPythonKeywords(code string) map[string]bool {
	keywords := make(map[string]bool)
	re := regexp.MustCompile(`\b(def|class|if|elif|else|for|while|try|except|finally|with|as|import|from|lambda|return|yield|global|nonlocal|assert|pass|break|continue|and|or|not|in|is|del|exec|print)\b`)
	matches := re.FindAllString(code, -1)
	for _, match := range matches {
		keywords[match] = true
	}
	return keywords
}

func (ca *CodeAnalyzer) extractJavaScriptKeywords(code string) map[string]bool {
	keywords := make(map[string]bool)
	re := regexp.MustCompile(`\b(function|var|let|const|if|else|for|while|do|switch|case|default|break|continue|return|try|catch|finally|throw|class|extends|import|export|async|await|yield|typeof|instanceof|in|of|new|this|super|delete|void)\b`)
	matches := re.FindAllString(code, -1)
	for _, match := range matches {
		keywords[match] = true
	}
	return keywords
}

func (ca *CodeAnalyzer) extractCppKeywords(code string) map[string]bool {
	keywords := make(map[string]bool)
	re := regexp.MustCompile(`\b(public|private|protected|class|struct|union|enum|namespace|using|template|typename|const|static|extern|volatile|mutable|virtual|override|final|explicit|inline|friend|operator|if|else|for|while|do|switch|case|default|break|continue|return|try|catch|throw|new|delete|this|sizeof|typedef|auto|decltype|nullptr|constexpr|noexcept|static_assert|#include|#define|#ifdef|#ifndef|#endif)\b`)
	matches := re.FindAllString(code, -1)
	for _, match := range matches {
		keywords[match] = true
	}
	return keywords
}

// extractImports 提取导入语句
func (ca *CodeAnalyzer) extractImports(code string) []string {
	var imports []string
	lines := strings.Split(code, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "import") || strings.HasPrefix(line, "from") || strings.HasPrefix(line, "#include") {
			imports = append(imports, line)
		}
	}

	return imports
}

// extractFunctions 提取函数信息
func (ca *CodeAnalyzer) extractFunctions(code string) []FunctionInfo {
	var functions []FunctionInfo
	lines := strings.Split(code, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if ca.isFunctionDeclaration(line) {
			funcInfo := FunctionInfo{
				Name:   ca.extractFunctionName(line),
				Line:   i + 1,
				Params: ca.countFunctionParams(line),
			}
			functions = append(functions, funcInfo)
		}
	}

	return functions
}

// isFunctionDeclaration 判断是否为函数声明
func (ca *CodeAnalyzer) isFunctionDeclaration(line string) bool {
	switch ca.language {
	case "go":
		return strings.Contains(line, "func ") && strings.Contains(line, "(")
	case "java":
		return strings.Contains(line, "(") && strings.Contains(line, ")") && (strings.Contains(line, "public") || strings.Contains(line, "private") || strings.Contains(line, "protected"))
	case "python":
		return strings.HasPrefix(line, "def ")
	case "javascript":
		return strings.Contains(line, "function") || (strings.Contains(line, "(") && strings.Contains(line, ")") && strings.Contains(line, "=>"))
	case "cpp", "c++":
		return strings.Contains(line, "(") && strings.Contains(line, ")") && !strings.Contains(line, "if") && !strings.Contains(line, "while") && !strings.Contains(line, "for")
	}
	return false
}

// extractFunctionName 提取函数名
func (ca *CodeAnalyzer) extractFunctionName(line string) string {
	re := regexp.MustCompile(`\b(\w+)\s*\(`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// countFunctionParams 计算函数参数数量
func (ca *CodeAnalyzer) countFunctionParams(line string) int {
	re := regexp.MustCompile(`\([^)]*\)`)
	matches := re.FindString(line)
	if matches == "" {
		return 0
	}

	// 移除括号
	params := strings.Trim(matches, "()")
	if params == "" {
		return 0
	}

	// 分割参数
	paramList := strings.Split(params, ",")
	return len(paramList)
}

// countComments 计算注释行数
func (ca *CodeAnalyzer) countComments(code string) int {
	lines := strings.Split(code, "\n")
	count := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "/*") {
			count++
		}
	}

	return count
}

// calculateComplexity 计算代码复杂度
func (ca *CodeAnalyzer) calculateComplexity(code string) int {
	complexity := 1 // 基础复杂度

	// 计算控制流复杂度
	complexity += strings.Count(code, "if")
	complexity += strings.Count(code, "else")
	complexity += strings.Count(code, "for")
	complexity += strings.Count(code, "while")
	complexity += strings.Count(code, "switch")
	complexity += strings.Count(code, "case")
	complexity += strings.Count(code, "catch")
	complexity += strings.Count(code, "&&")
	complexity += strings.Count(code, "||")

	return complexity
}
