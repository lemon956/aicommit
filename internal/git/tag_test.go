package git

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGit_TagHelpers(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	dir := t.TempDir()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")

	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello\n"), 0o644))
	runGit(t, dir, "add", "a.txt")
	runGit(t, dir, "commit", "-m", "feat: init")
	runGit(t, dir, "tag", "-a", "v0.1.0", "-m", "Release v0.1.0")

	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello world\n"), 0o644))
	runGit(t, dir, "add", "a.txt")
	runGit(t, dir, "commit", "-m", "fix: bug")

	g := New(dir)

	latest, ok, err := g.LatestTag()
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "v0.1.0", latest)

	exists, err := g.TagExists("v0.1.0")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = g.TagExists("v9.9.9")
	require.NoError(t, err)
	assert.False(t, exists)

	rangeSpec := "v0.1.0..HEAD"
	subjects, truncated, err := g.CommitSubjects(rangeSpec, 50)
	require.NoError(t, err)
	assert.False(t, truncated)
	assert.Contains(t, subjects, "fix: bug")
	assert.NotContains(t, subjects, "feat: init")

	stat, err := g.DiffStat(rangeSpec)
	require.NoError(t, err)
	assert.Contains(t, stat, "a.txt")

	nameStatus, err := g.DiffNameStatus(rangeSpec)
	require.NoError(t, err)
	assert.Contains(t, nameStatus, "a.txt")

	err = g.CreateAnnotatedTag("v0.2.0", "Release v0.2.0\n\nAdded\n- Something\n")
	require.NoError(t, err)

	exists, err = g.TagExists("v0.2.0")
	require.NoError(t, err)
	assert.True(t, exists)

	contents := runGit(t, dir, "for-each-ref", "refs/tags/v0.2.0", "--format=%(contents)")
	assert.Contains(t, contents, "Release v0.2.0")
	assert.Contains(t, contents, "Added")
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("git %s failed: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}

	return stdout.String()
}
