package prompt

import (
	"strings"
	"testing"
)

func TestTagMessageTemplate_GeneratePrompt(t *testing.T) {
	tpl := NewTagTemplate()
	p := tpl.GeneratePrompt("Release version: v1.2.3\nPrevious tag: v1.2.2\nRange: v1.2.2..HEAD\n")

	if p == "" {
		t.Fatal("prompt should not be empty")
	}
	if !containsAll(p,
		"<context>",
		"Release <version>",
		"Added",
		"Changed",
		"Fixed",
		"Breaking Changes",
	) {
		t.Fatalf("prompt missing expected content:\n%s", p)
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
