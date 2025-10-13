package model

import "context"

type Provider interface {
	GenerateCommitMessage(ctx context.Context, diff string) (string, error)
	Name() string
}
