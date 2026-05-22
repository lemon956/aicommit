package prompt

import (
	"strings"
	"testing"
)

func TestDefaultTemplateDoesNotUseHardLengthLimits(t *testing.T) {
	tpl := NewDefaultTemplate()
	p := strings.ToLower(tpl.GeneratePrompt("diff --git a/file b/file"))

	forbidden := []string{
		"hard limit",
		"wrap lines to <=",
	}
	for _, phrase := range forbidden {
		if strings.Contains(p, phrase) {
			t.Fatalf("default prompt should not contain hard length limit phrase %q:\n%s", phrase, p)
		}
	}
}
