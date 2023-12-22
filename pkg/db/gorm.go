package db

import (
	"context"
	"errors"
	"golinkcut/internal/config"
	"golinkcut/pkg/db/memory"
	"golinkcut/pkg/db/postgresql"
	"gorm.io/gorm"
	"log"
)

// Link is a struct representing original link and it's short alias.
// The fields are set up to both have a unique value in DB
type Link struct {
	gorm.Model
	// Alias is the short link alias
	Alias string `gorm:"uniqueIndex"`
	// Original is the original link to be found by Alias
	Original string `gorm:"type:text;not null"`
}

// Storage encapsulates logic for operating a DB, reading and writing information
type Storage struct {
	db *gorm.DB
}

// Add is used to add a new Link to the DB.
// Add will return ErrAliasTaken error if a passed alias exists in DB.
// Add will return ErrLinkExists if a passed original link exists in DB (and
// has an assigned alias already)
func (s Storage) Add(ctx context.Context, alias string, original string) error {
	tx := s.db.WithContext(ctx).Create(&Link{Alias: alias, Original: original})
	if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return ErrAliasExists{alias: alias}
	}
	return tx.Error
}

// Get is used to retrieve the original link, given its alias in DB
func (s Storage) Get(ctx context.Context, alias string) (string, error) {
	var link Link
	tx := s.db.WithContext(ctx).First(&link, "alias = ?", alias)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return "", ErrNotExists{alias: alias}
		}
		return "", tx.Error
	}
	return link.Original, nil
}

// NewStorage reads config.Config and returns a corresponding Storage object
func NewStorage(cfg config.Config) Storage {
	var db *gorm.DB
	log.Printf("Will connect to database now")
	useMemory := cfg["memory"].(bool)
	// NOTE: here congig['memory'] field is used "downstream" to define whether
	// to use a PostgreSQL connection or a SQLite in-memory database. If a
	// key-value storage is to be used, the "memory" flag should be used
	// earlier to define it's use.
	if useMemory {
		log.Printf("Will use in memory db")
		db = memory.NewDb(cfg)
	} else {
		log.Printf("Will connect to postgreSQL now")
		db = postgresql.NewDb(cfg)
	}
	err := db.AutoMigrate(&Link{})
	if err != nil {
		panic(err)
	}
	return Storage{db}
}
