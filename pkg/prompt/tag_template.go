package prompt

import "fmt"

// TagMessageTemplate generates prompts for annotated git tag messages (release notes).
type TagMessageTemplate struct {
	systemPrompt string
	userPrompt   string
}

func NewTagTemplate() *TagMessageTemplate {
	return &TagMessageTemplate{
		systemPrompt: `You are a senior software engineer and release manager. You write accurate, concise Git annotated tag messages (release notes) based strictly on the provided context.`,
		userPrompt: `Write an annotated Git tag message (release notes) for the release described below.

<context>
%s
</context>

RULES:
1. Output ONLY the tag message (no Markdown, no quotes, no code fences).
2. Language: English.
3. First line: "Release <version>" (use the version provided in context).
4. Use these sections when applicable (omit empty sections):
   - Added
   - Changed
   - Fixed
   - Breaking Changes
5. Use bullet points under each section. Keep bullets concrete and user-facing.
6. Base the content ONLY on the provided context. Do not invent changes.
7. If the context is insufficient, say so briefly and conservatively.
`,
	}
}

func (t *TagMessageTemplate) GeneratePrompt(input string) string {
	return fmt.Sprintf(t.userPrompt, input)
}

func (t *TagMessageTemplate) GetSystemPrompt() string {
	return t.systemPrompt
}
