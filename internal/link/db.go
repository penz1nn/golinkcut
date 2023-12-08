package link

import (
	"context"
	"golinkcut/internal/config"
	"golinkcut/internal/entity"
	"golinkcut/pkg/db"
	"strings"
)

type storage struct {
	db db.Storage
}

func (s storage) GetLink(ctx context.Context, alias string) (entity.Link, error) {
	originalLink, err := s.db.Get(ctx, alias)
	if err != nil {
		if strings.Contains(err.Error(), db.ErrNotExistsSignature) {
			return entity.Link{}, ErrNotExists{alias: alias}
		}
	}
	return entity.Link{Alias: alias, Original: originalLink}, err
}

func (s storage) SaveLink(ctx context.Context, link entity.Link) error {
	err := s.db.Add(ctx, link.Alias, link.Original)
	if err != nil {
		if strings.Contains(err.Error(), db.ErrAliasTakenSignature1) {
			return ErrAliasTaken{Alias: link.Alias}
		}
		if strings.Contains(err.Error(), db.ErrAliasTakenSignature2) {
			return ErrAliasTaken{Alias: link.Alias}
		}
		if strings.Contains(err.Error(), db.ErrLinkExistsSignature1) {
			return ErrLinkExists{Original: link.Original}
		}
		if strings.Contains(err.Error(), db.ErrLinkExistsSignature2) {
			return ErrLinkExists{Original: link.Original}
		}
	}
	return err
}

func NewStorage(config config.Config) storage {
	return storage{db.NewStorage(config)}
}
