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
	return fmt.Sprintf("short link alias %s not found in storage", e.alias)
}

// ErrAliasTaken error is returned in case there was an attempt to save a link
// with a short alias which exists in the DB
type ErrAliasTaken struct {
	Alias string
}

func (e ErrAliasTaken) Error() string {
	return fmt.Sprintf("Alias %v already exists in db", e.Alias)
}

// TODO return code 409 "Conflict" when error occurs
// ErrLinkExists error is returned in case there was an attempt to save a link
// which is already present in DB (and has an assigned short alias)
type ErrLinkExists struct {
	Original string
}

func (e ErrLinkExists) Error() string {
	return fmt.Sprintf("Link %v already exists in db", e.Original)
}
