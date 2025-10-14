package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateCommitMessage(message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	if len(message) > 80 {
		return fmt.Errorf("commit message too long: %d characters (max 80)", len(message))
	}

	if !isValidFormat(message) {
		return fmt.Errorf("invalid commit message format. Expected: type: description")
	}

	return nil
}

func isValidFormat(message string) bool {
	conventionalPattern := `^(feat|fix|docs|style|refactor|test|chore|perf|build|ci)(\(.+\))?: .+`
	matched, err := regexp.MatchString(conventionalPattern, message)
	if err != nil {
		return false
	}
	return matched
}

func CleanCommitMessage(message string) string {
	message = strings.TrimSpace(message)
	message = strings.Trim(message, "`")
	message = strings.Trim(message, "\"'")

	lines := strings.Split(message, "\n")
	if len(lines) > 0 {
		message = lines[0]
	}

	return strings.TrimSpace(message)
}
