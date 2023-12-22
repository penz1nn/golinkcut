package link

import (
	"context"
	"fmt"
	"golinkcut/internal/entity"
)

type Repository interface {
	GetLink(ctx context.Context, alias string) (entity.Link, error)
	SaveLink(ctx context.Context, link entity.Link) error
}

type ErrNotExists struct {
	alias string
}

func (e ErrNotExists) Error() string {
	return fmt.Sprintf("short link alias %s not found", e.alias)
}

// ErrLinkExists error is returned in case there was an attempt to save a link
// which is already present in DB (and has an assigned short alias)
type ErrLinkExists struct {
	Original string
}

func (e ErrLinkExists) Error() string {
	return fmt.Sprintf("Link %v already exists", e.Original)
}
