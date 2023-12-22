// Package memory contains code for use of in-memory key-value storage
package memory

import (
	"context"
	"github.com/eko/gocache/lib/v4/cache"
	gocachestore "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
	"log"
	"time"
)

// Storage struct encapsulates an instance of in-memory key-value storage
type Storage struct {
	cm *cache.Cache[[]byte]
}

// NewStorage creates a new Storage instance
func NewStorage() *Storage {
	// Client of patrickmn/go-cache package
	gocacheClient := gocache.New(5*time.Minute, 10*time.Minute)
	// Store based on patrickmn/go-cache for eko-gocache
	gocacheStore := gocachestore.NewGoCache(gocacheClient)
	cm := cache.New[[]byte](gocacheStore)
	return &Storage{cm}
}

// SetKey receives string values of key and value to create such record in the
// storage. It does not check whether the key already exists, so such check has
// to be performed upstream.
func (s *Storage) SetKey(ctx context.Context, key string, value string) {
	err := s.cm.Set(ctx, key, []byte(value))
	if err != nil {
		log.Printf("Unexpected error when setting key %s to value %s in storage: %s", key, value, err)
	}
}

// GetKey looks up a provided key in the storage and returns its string value.
// If key does not exist, it returns an empty string.
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
