package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Git struct {
	workDir string
}

func New(workDir string) *Git {
	return &Git{workDir: workDir}
}

func (g *Git) GetDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	cmd.Dir = g.workDir

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	diff := stdout.String()
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no staged changes found")
	}

	// 如果diff太长，截取一部分（限制在10000字符以内）
	if len(diff) > 10000 {
		diff = diff[:10000] + "\n... (diff truncated due to length)"
	}

	return diff, nil
}

func (g *Git) Commit(message string) error {
	tmpFile, err := os.CreateTemp("", "aicommit-commit-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp commit message file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	if _, err := tmpFile.WriteString(message); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write commit message to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp commit message file: %w", err)
	}

	cmd := exec.Command("git", "commit", "-F", tmpFile.Name())
	cmd.Dir = g.workDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (g *Git) IsRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = g.workDir

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	gitDir := strings.TrimSpace(string(output))
	return gitDir != ""
}
