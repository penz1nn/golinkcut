package link

import (
	"context"
	"golinkcut/internal/entity"
	"golinkcut/pkg/memory"
)

type memoryStorage struct {
	*memory.Storage
}

func (ms *memoryStorage) Save(ctx context.Context, link entity.Link) error {
	_, err := ms.Get(ctx, link.Short)
	if err != nil {
		ms.SetKey(ctx, link.Short, link.Original)
		return nil
	}
	return ShortlinkTakenError{
		shortLink:    link.Short,
		originalLink: link.Original,
	}
}

func (ms *memoryStorage) Get(ctx context.Context, shortLink string) (entity.Link, error) {
	originalLink := ms.GetKey(ctx, shortLink)
	if originalLink == "" {
		return entity.Link{}, NotFoundError{shortLink: shortLink}
	}
	return entity.Link{Short: shortLink, Original: originalLink}, nil
}
