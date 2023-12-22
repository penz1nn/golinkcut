package link

import (
	"context"
	"fmt"
	"golinkcut/internal/entity"
)

// Repository interface defines behavior to save and get entity.Link instances
type Repository interface {
	// GetLink receives link's short alias to fetch entity.Link instance
	GetLink(ctx context.Context, alias string) (entity.Link, error)

	// SaveLink receives an entity.Link instance to save it in storage
	SaveLink(ctx context.Context, link entity.Link) error
}

// ErrNotExists represents an error which is thrown when there is no entity.Link
// instance saved in storage which contains given alias
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
