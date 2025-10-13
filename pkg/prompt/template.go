package prompt

import "fmt"

type Template interface {
	GeneratePrompt(diff string) string
	GetSystemPrompt() string
}

type CommitMessageTemplate struct {
	systemPrompt string
	userPrompt   string
}

func NewDefaultTemplate() *CommitMessageTemplate {
	return &CommitMessageTemplate{
		systemPrompt: `You are a senior software engineer writing Git commit messages. Write concise, specific messages that precisely describe the code changes. Never use generic or vague language.`,
		userPrompt: `Generate a commit message for this diff. Be extremely specific about what functionality was added or changed.

<diff>
%s
</diff>

RULES:
1. Format: <type>(<scope>): <precise description>
2. Types: feat|fix|refactor|perf|docs|style|test|chore|ci|build
3. Scope: file/module name (auth, api, db, cli, config, etc.)
4. Description: imperative mood, lowercase, no period
5. CRITICAL: Total message must be under 80 characters
6. Be SPECIFIC about functionality, not files

ANALYZE:
- What new capability was added? (feat)
- What bug was fixed? (fix) 
- What was restructured? (refactor)
- What configuration/tooling changed? (chore)
- Read the CODE CONTENT, not just filenames

FORBIDDEN PATTERNS:
❌ "add/update files" - too vague
❌ "initial commit" - describe what was created
❌ "various changes" - be specific
❌ "fix issues" - specify what was fixed
❌ Listing files - describe functionality instead

EXAMPLES:
✓ feat(auth): implement JWT token authentication
✓ fix(validator): handle null email addresses
✓ refactor(http): extract request handling to middleware
✓ chore(deps): upgrade go modules to latest versions
✓ feat(cli): add dry-run flag for commit preview

If this is an initial commit with project setup:
✓ feat: initialize aicommit CLI tool
✓ feat: create git commit message generator
✓ chore: setup Go project with dependencies

Return ONLY the commit message.`,
	}
}

func (t *CommitMessageTemplate) GeneratePrompt(diff string) string {
	return fmt.Sprintf(t.userPrompt, diff)
}

func (t *CommitMessageTemplate) GetSystemPrompt() string {
	return t.systemPrompt
}

func (t *CommitMessageTemplate) SetSystemPrompt(prompt string) {
	t.systemPrompt = prompt
}

func (t *CommitMessageTemplate) SetUserPrompt(prompt string) {
	t.userPrompt = prompt
}

type CustomTemplate struct {
	systemPrompt string
	userPrompt   string
}

func NewCustomTemplate(systemPrompt, userPrompt string) *CustomTemplate {
	return &CustomTemplate{
		systemPrompt: systemPrompt,
		userPrompt:   userPrompt,
	}
}

func (t *CustomTemplate) GeneratePrompt(diff string) string {
	return fmt.Sprintf(t.userPrompt, diff)
}

func (t *CustomTemplate) GetSystemPrompt() string {
	return t.systemPrompt
}

var globalTemplate Template

func init() {
	globalTemplate = NewDefaultTemplate()
}

func GetGlobalTemplate() Template {
	return globalTemplate
}

func SetGlobalTemplate(template Template) {
	globalTemplate = template
}

func ResetToDefault() {
	globalTemplate = NewDefaultTemplate()
}
