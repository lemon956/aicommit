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
			name:    "valid conventional commit",
			message: "feat: add user authentication",
			wantErr: false,
		},
		{
			name:    "valid conventional commit with scope",
			message: "fix(auth): resolve login issue",
			wantErr: false,
		},
		{
			name:    "empty message",
			message: "",
			wantErr: true,
		},
		{
			name:    "too long message",
			message: "feat: this is a very long commit message that exceeds the 250 character limit and should be rejected because it contains way too much text and goes on and on with unnecessary details that make the commit message extremely verbose and difficult to read in git logs",
			wantErr: true,
		},
		{
			name:    "valid long message under limit",
			message: "feat: this is a reasonably long commit message that is under the 250 character limit so it should be accepted even though it contains quite a bit of text to describe the changes being made",
			wantErr: false,
		},
		{
			name:    "invalid format",
			message: "add user authentication",
			wantErr: true,
		},
		{
			name:    "invalid type",
			message: "invalid: add user authentication",
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
			name:     "message with quotes",
			message:  "\"feat: add user authentication\"",
			expected: "feat: add user authentication",
		},
		{
			name:     "message with newlines",
			message:  "feat: add user authentication\n\nAdditional details here",
			expected: "feat: add user authentication",
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
