package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit_IsRepository(t *testing.T) {
	// Test current directory
	git := New(".")
	isRepo := git.IsRepository()

	// Just check that the method works, don't assert on the result
	// since we don't know if we're in a git repo
	t.Logf("Current directory is git repository: %v", isRepo)

	// Test non-git directory
	nonGitDir := t.TempDir()
	gitNonRepo := New(nonGitDir)
	assert.False(t, gitNonRepo.IsRepository(), "Directory without .git should not be detected as repository")
}

func TestGit_GetDiff(t *testing.T) {
	git := New(".")

	if !git.IsRepository() {
		t.Skip("Not in a git repository")
	}

	_, err := git.GetDiff()
	assert.NoError(t, err)
}
