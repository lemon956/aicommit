package prompt

import "fmt"

// Template defines the interface for generating prompts
type Template interface {
	GeneratePrompt(diff string) string
	GetSystemPrompt() string
}

// CommitMessageTemplate is the hardcoded template for commit messages
type CommitMessageTemplate struct {
	systemPrompt string
	userPrompt   string
}

// NewDefaultTemplate creates a new instance of the default template
func NewDefaultTemplate() *CommitMessageTemplate {
	return &CommitMessageTemplate{
		systemPrompt: `You are a senior software engineer. You write Git commit messages that follow the gitcommit(5) guidelines AND Conventional Commits for the subject line. Be concise, concrete, and accurate.`,
		userPrompt: `Write a Git commit message for this diff. Focus on what changed and why it matters (not filenames).

<diff>
%s
</diff>

RULES:
1. Output ONLY the commit message (no Markdown, no quotes, no code fences).
2. Use this structure:
   - Line 1: subject (summary)
   - Optional: blank line
   - Optional: body (one or more lines)
3. Subject:
   - MUST use Conventional Commits format: <type>(<scope>)?: <summary>
   - Allowed types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
   - Use an optional scope when it helps (e.g. feat(auth): ...)
   - English, imperative mood if possible
   - No trailing period
   - Aim for <= 50 characters; hard limit: 100 characters
4. Body (only if needed to explain why/impact/behavior change):
   - MUST be separated from the subject by a blank line
   - Wrap lines to <= 120 characters
   - Explain WHAT and WHY; avoid implementation details unless necessary
5. If no body is needed, return a single-line subject.

ANALYZE:
- What new capability was added?
- What bug was fixed?
- What behavior changed and why?
- What risks or compatibility notes matter?
- Read the CODE CONTENT, not just filenames

EXAMPLES:
✓ feat(cli): add dry-run flag to preview generated commit message
✓ fix(git): avoid panic when staged diff is empty
✓ refactor(config): support XDG_CONFIG_HOME
✓ chore(deps): update dependencies to address security advisory

✓ feat(editor): add editor support for reviewing commit message

This lets users edit the generated message before committing and reduces
incorrect commits caused by prompt misunderstandings.`,
	}
}

func (t *CommitMessageTemplate) GeneratePrompt(diff string) string {
	return fmt.Sprintf(t.userPrompt, diff)
}

func (t *CommitMessageTemplate) GetSystemPrompt() string {
	return t.systemPrompt
}

var globalTemplate Template = NewDefaultTemplate()

// GetGlobalTemplate returns the global template instance
func GetGlobalTemplate() Template {
	return globalTemplate
}
