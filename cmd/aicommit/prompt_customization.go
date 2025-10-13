package main

import (
	"github.com/aicommit/aicommit/pkg/prompt"
)

// SetCustomPrompt 设置自定义prompt模板
func SetCustomPrompt(systemPrompt, userPrompt string) {
	customTemplate := prompt.NewCustomTemplate(systemPrompt, userPrompt)
	prompt.SetGlobalTemplate(customTemplate)
}

// SetDefaultPrompt 重置为默认prompt
func SetDefaultPrompt() {
	prompt.ResetToDefault()
}

// SetChinesePrompt 设置为中文prompt（示例）
func SetChinesePrompt() {
	systemPrompt := "你是一个专业的代码审查助手，擅长根据代码变更生成符合规范的提交信息。"
	userPrompt := `基于以下git diff，生成一个符合常规提交格式的中文提交信息：

<diff>
%s
</diff>

要求：
1. 使用常规提交格式（类型: 描述）
2. 控制在72个字符以内
3. 使用现在时态
4. 具体说明变更内容
5. 如果有多个变更，专注于最重要的一个

有效类型：feat（功能）, fix（修复）, docs（文档）, style（格式）, refactor（重构）, test（测试）, chore（杂务）, perf（性能）, build（构建）, ci（CI）

只返回提交信息，不要额外的文字。`

	SetCustomPrompt(systemPrompt, userPrompt)
}

// SetDetailedPrompt 设置为更详细的prompt（示例）
func SetDetailedPrompt() {
	systemPrompt := "You are an expert software engineer who writes comprehensive and descriptive commit messages."
	userPrompt := `Based on the following git diff, generate a detailed commit message following the conventional commits format:

<diff>
%s
</diff>

Requirements:
1. Use conventional commits format (type: description)
2. Keep the first line under 72 characters
3. Use present tense
4. Be specific about what changed and why
5. If there are multiple changes, focus on the most important one
6. Consider the impact and reasoning behind the changes
7. Use technical but clear language

Valid types: feat, fix, docs, style, refactor, test, chore, perf, build, ci

Return only the commit message, no additional text.`

	SetCustomPrompt(systemPrompt, userPrompt)
}

// SetMinimalPrompt 设置为简洁的prompt（示例）
func SetMinimalPrompt() {
	systemPrompt := "Generate concise git commit messages."
	userPrompt := `From this diff: %s

Write a short commit message (max 50 chars) in format "type: what changed".
Types: feat, fix, docs, style, refactor, test, chore, perf, build, ci`

	SetCustomPrompt(systemPrompt, userPrompt)
}

// GetCurrentPromptInfo 获取当前prompt信息（调试用）
func GetCurrentPromptInfo() (systemPrompt string, userPrompt string) {
	template := prompt.GetGlobalTemplate()
	return template.GetSystemPrompt(), "Based on the following git diff, generate a concise commit message..."
}
