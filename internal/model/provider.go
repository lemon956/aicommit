package model

import (
	"context"

	"github.com/aicommit/aicommit/pkg/prompt"
)

type Provider interface {
	GenerateMessage(ctx context.Context, input string) (string, error)
	SetTemplate(template prompt.Template)
	Name() string
}
