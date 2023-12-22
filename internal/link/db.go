package link

import (
	"context"
	"errors"
	"golinkcut/internal/config"
	"golinkcut/internal/entity"
	"golinkcut/pkg/db"
)

type dbStorage struct {
	db db.Storage
}

func (s dbStorage) GetLink(ctx context.Context, alias string) (entity.Link, error) {
	originalLink, err := s.db.Get(ctx, alias)
	if err != nil {
		var ref db.ErrNotExists
		if errors.As(err, &ref) {
			return entity.Link{}, ErrNotExists{alias: alias}
		}
	}
	return entity.Link{Alias: alias, Original: originalLink}, err
}

func (s dbStorage) SaveLink(ctx context.Context, link entity.Link) error {
	err := s.db.Add(ctx, link.Alias, link.Original)
	if err != nil {
		var ref1 db.ErrAliasExists
		if errors.As(err, &ref1) {
			return ErrLinkExists{Original: link.Original}
		}
	}
	return err
}

func NewDbStorage(config config.Config) dbStorage {
	return dbStorage{db.NewStorage(config)}
}
