package service

import (
	"context"
	"time"
)

type RepoI interface {
	Create(ctx context.Context, sub string, ttl time.Duration, payload any) (string, error)
	Validate(ctx context.Context, token string) (bool, []byte, error)
}
