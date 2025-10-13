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

	diff, err := git.GetDiff()
	
	// It's OK if there are no staged changes in CI environment
	if err != nil && err.Error() == "no staged changes found" {
		t.Log("No staged changes found - this is expected in CI")
		return
	}
	
	// For other errors or success, check the results
	assert.NoError(t, err)
	if err == nil {
		assert.NotEmpty(t, diff, "Diff should not be empty when there are staged changes")
	}
}
