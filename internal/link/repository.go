package link

import (
	"context"
	"fmt"
	"golinkcut/internal/entity"
)

type Repository interface {
	//TODO: use context or remove it from the interface
	Get(ctx context.Context, shortLink string) (entity.Link, error)
	Save(ctx context.Context, link entity.Link) error
}

type ShortlinkTakenError struct {
	shortLink    string
	originalLink string
}

func (e ShortlinkTakenError) Error() string {
	errorString := fmt.Sprintf("short link alias %s is already occupied", e.shortLink)
	if e.originalLink != "" {
		errorString = errorString + fmt.Sprintf(" by original link %s", e.originalLink)
	}
	return errorString
}

type NotFoundError struct {
	shortLink string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("short link alias %s not found in storage", e.shortLink)
}
