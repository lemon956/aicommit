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
			message: "feat: this is a very long commit message that exceeds the 72 character limit and should be rejected",
			wantErr: true,
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
