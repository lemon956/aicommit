package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCommitMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "valid subject only",
			message: "Add JWT auth to CLI login",
			wantErr: false,
		},
		{
			name:    "valid subject and body",
			message: "Subject line\n\nBody line 1\nBody line 2",
			wantErr: false,
		},
		{
			name:    "conventional commit still accepted as subject",
			message: "feat(auth): add JWT validation",
			wantErr: false,
		},
		{
			name:    "trailer line can exceed body limit",
			message: "Subject\n\nCo-authored-by: Very Long Name <email@domain.com>",
			wantErr: false,
		},
		{
			name:    "empty message",
			message: "",
			wantErr: true,
		},
		{
			name:    "subject too long",
			message: "This subject line is intentionally made longer than seventy-two characters to fail",
			wantErr: true,
		},
		{
			name:    "missing blank line between subject and body",
			message: "Subject line\nBody line 1",
			wantErr: true,
		},
		{
			name:    "body line too long",
			message: "Subject\n\nThis body line is intentionally made longer than seventy-two characters to fail",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommitMessage(tt.message)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCleanCommitMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "message with backticks",
			message:  "`feat: add user authentication`",
			expected: "feat: add user authentication",
		},
		{
			name:     "message in fenced code block",
			message:  "```text\nSubject line\n\nBody line\n```",
			expected: "Subject line\n\nBody line",
		},
		{
			name:     "message with quotes",
			message:  "\"feat: add user authentication\"",
			expected: "feat: add user authentication",
		},
		{
			name:     "message with newlines is preserved",
			message:  "Subject line\n\nAdditional details here",
			expected: "Subject line\n\nAdditional details here",
		},
		{
			name:     "windows newlines are normalized",
			message:  "Subject line\r\n\r\nBody line\r\n",
			expected: "Subject line\n\nBody line",
		},
		{
			name:     "message with whitespace",
			message:  "  feat: add user authentication  ",
			expected: "feat: add user authentication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanCommitMessage(tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}
