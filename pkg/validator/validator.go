package validator

import (
	"fmt"
	"os"
	"path/filepath"
)

func ValidateRepository(path string) error {
	if path == "" {
		path = "."
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	gitPath := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", absPath)
	}

	return nil
}
