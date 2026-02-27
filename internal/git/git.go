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

func validateTagName(tag string) error {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if strings.HasPrefix(tag, "-") {
		return fmt.Errorf("invalid tag name: cannot start with '-'")
	}
	if strings.ContainsAny(tag, " \t\n\r") {
		return fmt.Errorf("invalid tag name: cannot contain whitespace")
	}
	if strings.Contains(tag, "..") || strings.Contains(tag, "@{") || strings.Contains(tag, "//") {
		return fmt.Errorf("invalid tag name: contains invalid sequence")
	}
	// Conservative subset of git refname rules.
	if strings.ContainsAny(tag, "~^:?*[\\") {
		return fmt.Errorf("invalid tag name: contains invalid character")
	}
	if strings.HasSuffix(tag, "/") || strings.HasSuffix(tag, ".") || strings.HasSuffix(tag, ".lock") {
		return fmt.Errorf("invalid tag name: invalid suffix")
	}
	if strings.HasPrefix(tag, "/") {
		return fmt.Errorf("invalid tag name: cannot start with '/'")
	}
	return nil
}

func (g *Git) runGit(args ...string) (string, error) {
	// #nosec G204 -- We execute the git binary with explicit arguments (no shell).
	cmd := exec.Command("git", args...)
	cmd.Dir = g.workDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errText := strings.TrimSpace(stderr.String())
		if errText == "" {
			errText = strings.TrimSpace(stdout.String())
		}
		if errText != "" {
			return "", fmt.Errorf("git %s failed: %w: %s", strings.Join(args, " "), err, errText)
		}
		return "", fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}

	return stdout.String(), nil
}

func (g *Git) GetDiff() (string, error) {
	diff, err := g.runGit("diff", "--staged")
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no staged changes found")
	}

	// 如果diff太长，截取一部分（限制在256KiB以内）
	const maxDiffBytes = 256 * 1024
	if len(diff) > maxDiffBytes {
		diff = diff[:maxDiffBytes] + "\n... (diff truncated due to length)"
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

	// #nosec G204 -- We execute the git binary with explicit arguments (no shell).
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
	out, err := g.runGit("rev-parse", "--git-dir")
	if err != nil {
		return false
	}

	gitDir := strings.TrimSpace(out)
	return gitDir != ""
}

func (g *Git) LatestTag() (tag string, ok bool, err error) {
	// #nosec G204 -- We execute the git binary with explicit arguments (no shell).
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = g.workDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errText := strings.ToLower(stderr.String())
		// Common messages:
		// - "fatal: No names found, cannot describe anything."
		// - "fatal: No tags can describe ..."
		if strings.Contains(errText, "no names found") || strings.Contains(errText, "no tags") {
			return "", false, nil
		}
		return "", false, fmt.Errorf("failed to get latest tag: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	tag = strings.TrimSpace(stdout.String())
	if tag == "" {
		return "", false, nil
	}
	return tag, true, nil
}

func (g *Git) TagExists(tag string) (bool, error) {
	if err := validateTagName(tag); err != nil {
		return false, err
	}

	// #nosec G204 -- We execute the git binary with explicit arguments (no shell); tag is validated.
	cmd := exec.Command("git", "rev-parse", "-q", "--verify", "refs/tags/"+strings.TrimSpace(tag))
	cmd.Dir = g.workDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Not found -> exit status 1; treat as not existing.
		return false, nil
	}

	return strings.TrimSpace(stdout.String()) != "", nil
}

func (g *Git) CommitSubjects(rangeSpec string, max int) ([]string, bool, error) {
	args := []string{"log", "--pretty=format:%s"}
	rangeSpec = strings.TrimSpace(rangeSpec)
	if rangeSpec != "" {
		args = append(args, rangeSpec)
	}

	out, err := g.runGit(args...)
	if err != nil {
		return nil, false, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "" {
		lines = nil
	}

	truncated := false
	if max > 0 && len(lines) > max {
		lines = lines[:max]
		truncated = true
	}

	return lines, truncated, nil
}

func (g *Git) DiffStat(rangeSpec string) (string, error) {
	rangeSpec = strings.TrimSpace(rangeSpec)
	if rangeSpec == "" {
		return "", fmt.Errorf("rangeSpec cannot be empty")
	}
	return g.runGit("diff", "--stat", rangeSpec)
}

func (g *Git) DiffNameStatus(rangeSpec string) (string, error) {
	rangeSpec = strings.TrimSpace(rangeSpec)
	if rangeSpec == "" {
		return "", fmt.Errorf("rangeSpec cannot be empty")
	}
	return g.runGit("diff", "--name-status", rangeSpec)
}

func (g *Git) CreateAnnotatedTag(tag string, message string) error {
	if err := validateTagName(tag); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp("", "aicommit-tag-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp tag message file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	message = strings.TrimSpace(message)
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}

	if _, err := tmpFile.WriteString(message); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write tag message to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp tag message file: %w", err)
	}

	// #nosec G204 -- We execute the git binary with explicit arguments (no shell); tag is validated.
	cmd := exec.Command("git", "tag", "-a", strings.TrimSpace(tag), "-F", tmpFile.Name())
	cmd.Dir = g.workDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create annotated tag: %w", err)
	}

	return nil
}
