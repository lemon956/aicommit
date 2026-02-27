package prompt

import "testing"

func TestValidateConventionalCommitMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "valid feat",
			message: "feat: add thing",
			wantErr: false,
		},
		{
			name:    "valid type not case sensitive",
			message: "Feat: add thing",
			wantErr: false,
		},
		{
			name:    "valid scope",
			message: "fix(auth): handle missing token",
			wantErr: false,
		},
		{
			name:    "invalid scope with whitespace",
			message: "fix(auth api): handle missing token",
			wantErr: true,
		},
		{
			name:    "valid breaking marker",
			message: "refactor(api)!: rename endpoint",
			wantErr: false,
		},
		{
			name:    "missing type",
			message: "Add new thing",
			wantErr: true,
		},
		{
			name:    "custom type allowed by spec",
			message: "feature: add new thing",
			wantErr: false,
		},
		{
			name:    "custom type with hyphen",
			message: "build-system: update pipeline",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConventionalCommitMessage(tt.message)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
