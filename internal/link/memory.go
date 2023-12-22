package link

import (
	"context"
	"golinkcut/internal/entity"
	"golinkcut/pkg/memory"
)

type kvStorage struct {
	s *memory.Storage
}

func (kv kvStorage) GetLink(ctx context.Context, alias string) (entity.Link, error) {
	result := kv.s.GetKey(ctx, alias)
	if result == "" {
		return entity.Link{}, ErrNotExists{alias: alias}
	}
	l := entity.Link{Alias: alias, Original: result}
	return l, nil
}

func (kv kvStorage) SaveLink(ctx context.Context, link entity.Link) error {
	saved := kv.s.GetKey(ctx, link.Alias)
	if saved != "" {
		return ErrLinkExists{Original: link.Original}
	}
	kv.s.SetKey(ctx, link.Alias, link.Original)
	return nil
}

func NewKvStorage() kvStorage {
	return kvStorage{memory.NewStorage()}
}
