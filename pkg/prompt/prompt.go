package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	maxCommitMessageLength = 5000
	maxSubjectLength       = 100
	maxBodyLineLength      = 120
	maxTrailerLineLength   = 200
)

var trailerPattern = regexp.MustCompile(`^[A-Za-z-]+: `)
var conventionalSubjectPattern = regexp.MustCompile(`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([^)]+\))?(!)?: .+`)

func ValidateCommitMessage(message string) error {
	message = strings.TrimSpace(normalizeNewlines(message))
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	if len(message) > maxCommitMessageLength {
		return fmt.Errorf("commit message too long: %d characters (max %d)", len(message), maxCommitMessageLength)
	}

	lines := strings.Split(message, "\n")
	subject := strings.TrimRight(lines[0], " \t")
	if subject == "" {
		return fmt.Errorf("commit subject cannot be empty")
	}
	if len(subject) > maxSubjectLength {
		return fmt.Errorf("commit subject too long: %d characters (max %d)", len(subject), maxSubjectLength)
	}

	if hasBody(lines[1:]) {
		if len(lines) < 2 || lines[1] != "" {
			return fmt.Errorf("invalid commit message format: separate subject and body with a blank line")
		}

		for i, line := range lines[2:] {
			if line == "" {
				continue
			}
			limit := maxBodyLineLength
			if isTrailerLine(line) {
				limit = maxTrailerLineLength
			}
			if len(line) > limit {
				return fmt.Errorf("commit body line too long at line %d: %d characters (max %d)", i+3, len(line), limit)
			}
		}
	}

	return nil
}

func ValidateConventionalCommitMessage(message string) error {
	message = strings.TrimSpace(normalizeNewlines(message))
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	lines := strings.Split(message, "\n")
	subject := strings.TrimRight(lines[0], " \t")
	if subject == "" {
		return fmt.Errorf("commit subject cannot be empty")
	}
	if !conventionalSubjectPattern.MatchString(subject) {
		return fmt.Errorf("commit subject must use Conventional Commits format: <type>(<scope>)?: <summary>")
	}
	return nil
}

func hasBody(lines []string) bool {
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			return true
		}
	}
	return false
}

func isTrailerLine(line string) bool {
	return trailerPattern.MatchString(line)
}

func CleanCommitMessage(message string) string {
	return CleanAIText(message)
}

func CleanAIText(text string) string {
	text = strings.TrimSpace(normalizeNewlines(text))

	if extracted, ok := extractFirstFencedCodeBlock(text); ok {
		text = extracted
	}

	text = strings.TrimSpace(text)
	text = trimMatchingWrapper(text, '`')
	text = trimMatchingWrapper(text, '"')
	text = trimMatchingWrapper(text, '\'')

	return strings.TrimSpace(text)
}

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func extractFirstFencedCodeBlock(s string) (string, bool) {
	start := strings.Index(s, "```")
	if start == -1 {
		return "", false
	}

	afterFence := s[start+3:]
	newline := strings.Index(afterFence, "\n")
	if newline == -1 {
		return "", false
	}

	contentStart := start + 3 + newline + 1
	endRel := strings.Index(s[contentStart:], "```")
	if endRel == -1 {
		return "", false
	}

	contentEnd := contentStart + endRel
	return strings.TrimSpace(s[contentStart:contentEnd]), true
}

func trimMatchingWrapper(s string, wrapper byte) string {
	if len(s) >= 2 && s[0] == wrapper && s[len(s)-1] == wrapper {
		return s[1 : len(s)-1]
	}
	return s
}
