package service

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo RepoI
}

func New(repo RepoI) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, sub string, ttl time.Duration, payload any) (string, error) {
	result, err := s.repo.Create(ctx, sub, ttl, payload)
	if err != nil {
		return "", fmt.Errorf("repo.Create: %w", err)
	}

	return result, nil
}

func (s *Service) Validate(ctx context.Context, token string) (bool, []byte, error) {
	valid, claims, err := s.repo.Validate(ctx, token)
	if err != nil {
		return false, nil, fmt.Errorf("repo.Validate: %w", err)
	}

	return valid, claims, nil
}
