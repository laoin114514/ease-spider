package utils

import (
	"regexp"
	"strings"
)

type MarkdownParser struct {
	content string
}

func NewMarkdownParser(content string) *MarkdownParser {
	return &MarkdownParser{
		content: content,
	}
}

type CodeBlock struct {
	Language string // 编程语言
	Content  string // 代码内容
	Index    int    // 在文档中的索引位置
}

// ExtractCodeBlocks 提取所有代码块
func (mp *MarkdownParser) ExtractCodeBlocks() []CodeBlock {
	var codeBlocks []CodeBlock

	// 匹配 ```language 格式的代码块
	codeBlockRegex := regexp.MustCompile("```(\\w*)\\n([\\s\\S]*?)```")
	matches := codeBlockRegex.FindAllStringSubmatch(mp.content, -1)

	for i, match := range matches {
		if len(match) >= 3 {
			language := match[1]
			content := strings.TrimSpace(match[2])

			codeBlocks = append(codeBlocks, CodeBlock{
				Language: language,
				Content:  content,
				Index:    i,
			})
		}
	}

	return codeBlocks
}

// ExtractCodeBlocksByLanguage 按编程语言提取代码块
func (mp *MarkdownParser) ExtractCodeBlocksByLanguage(targetLanguage string) []CodeBlock {
	var filteredBlocks []CodeBlock
	allBlocks := mp.ExtractCodeBlocks()

	for _, block := range allBlocks {
		if strings.EqualFold(block.Language, targetLanguage) {
			filteredBlocks = append(filteredBlocks, block)
		}
	}

	return filteredBlocks
}

// GetCodeContentArray 获取所有代码块的内容数组
func (mp *MarkdownParser) GetCodeContentArray() []string {
	var contents []string
	codeBlocks := mp.ExtractCodeBlocks()

	for _, block := range codeBlocks {
		contents = append(contents, block.Content)
	}

	return contents
}

// GetCodeContentArrayByLanguage 按语言获取代码块内容数组
func (mp *MarkdownParser) GetCodeContentArrayByLanguage(language string) []string {
	var contents []string
	codeBlocks := mp.ExtractCodeBlocksByLanguage(language)

	for _, block := range codeBlocks {
		contents = append(contents, block.Content)
	}

	return contents
}

// CountCodeBlocks 统计代码块数量
func (mp *MarkdownParser) CountCodeBlocks() int {
	return len(mp.ExtractCodeBlocks())
}

// CountCodeBlocksByLanguage 按语言统计代码块数量
func (mp *MarkdownParser) CountCodeBlocksByLanguage(language string) int {
	return len(mp.ExtractCodeBlocksByLanguage(language))
}

// GetCodeBlockInfo 获取代码块详细信息
func (mp *MarkdownParser) GetCodeBlockInfo() map[string]interface{} {
	codeBlocks := mp.ExtractCodeBlocks()

	// 统计各语言的代码块数量
	languageCount := make(map[string]int)
	for _, block := range codeBlocks {
		languageCount[block.Language]++
	}

	return map[string]interface{}{
		"total_blocks":   len(codeBlocks),
		"language_count": languageCount,
		"code_blocks":    codeBlocks,
	}
}
