package link

import (
	"context"
	"golinkcut/internal/entity"
)

// Service encapsulates usecase logic for links
type Service interface {
	Get(ctx context.Context, shortLink string) (Link, error)
	Create(ctx context.Context, input CreateLinkRequest) (Link, error)
}

// Link represents a single pair of shortened and full links
type Link struct {
	entity.Link
}

type service struct {
	repo Repository
}

//func NewService(repo Repository, logger log.Logger)
