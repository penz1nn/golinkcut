package memory

import (
	"context"
	"github.com/eko/gocache/lib/v4/cache"
	gocachestore "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
	"log"
	"time"
)

type Storage struct {
	cm *cache.Cache[[]byte]
}

func NewStorage() *Storage {
	// Client of patrickmn/go-cache package
	gocacheClient := gocache.New(5*time.Minute, 10*time.Minute)
	// Store based on patrickmn/go-cache for eko-gocache
	gocacheStore := gocachestore.NewGoCache(gocacheClient)
	cm := cache.New[[]byte](gocacheStore)
	return &Storage{cm}
}

func (s *Storage) SetKey(ctx context.Context, key string, value string) {
	err := s.cm.Set(ctx, key, []byte(value))
	if err != nil {
		log.Printf("Unexpected error when setting key %s to value %s in storage: %s", key, value, err)
	}
}

func (s *Storage) GetKey(ctx context.Context, key string) string {
	value, err := s.cm.Get(ctx, key)
	if err != nil {
		if err.Error() == "value not found in store" {
			return ""
		}
		log.Printf("Unexpected error when fetching key %s from cache: %s", key, err)
	}
	return string(value)
}
