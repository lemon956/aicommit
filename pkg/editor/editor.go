package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Open opens the content in the specified editor or system default.
// It returns the edited content.
func Open(content string, editorCmd string) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "aicommit-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to temp file
	if _, err := tmpFile.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	// Determine editor
	if editorCmd == "" {
		editorCmd = os.Getenv("EDITOR")
		if editorCmd == "" {
			editorCmd = os.Getenv("VISUAL")
			if editorCmd == "" {
				// Default fallbacks
				if _, err := exec.LookPath("nvim"); err == nil {
					editorCmd = "nvim"
				} else if _, err := exec.LookPath("vim"); err == nil {
					editorCmd = "vim"
				} else if _, err := exec.LookPath("nano"); err == nil {
					editorCmd = "nano"
				} else if _, err := exec.LookPath("vi"); err == nil {
					editorCmd = "vi"
				} else {
					return "", fmt.Errorf("no editor found. Please set EDITOR environment variable or configure 'editor' in config file")
				}
			}
		}
	}

	// Open editor
	parts := strings.Fields(editorCmd)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid editor command")
	}

	path := parts[0]
	args := parts[1:]
	args = append(args, tmpFile.Name())

	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	// Read content back
	newContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read back temp file: %w", err)
	}

	return string(newContent), nil
}
