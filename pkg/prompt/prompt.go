package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

// Conventional Commits v1.0.0 summary format:
// <type>[optional scope][optional !]: <description>
//
// Notes:
// - The spec says types are not case sensitive for implementors.
// - We allow a conservative subset for type/scope characters to prevent malformed subjects.
var conventionalSubjectPattern = regexp.MustCompile(`(?i)^[a-z][a-z0-9-]*(\([^\s)]+\))?(!)?: .+`)

func ValidateCommitMessage(message string) error {
	message = strings.TrimSpace(normalizeNewlines(message))
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	lines := strings.Split(message, "\n")
	subject := strings.TrimRight(lines[0], " \t")
	if subject == "" {
		return fmt.Errorf("commit subject cannot be empty")
	}

	if hasBody(lines[1:]) {
		if len(lines) < 2 || lines[1] != "" {
			return fmt.Errorf("invalid commit message format: separate subject and body with a blank line")
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
		return fmt.Errorf("commit subject must use Conventional Commits v1.0.0 summary format: <type>(<scope>)?!: <description>")
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
