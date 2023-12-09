package db

import (
	"context"
	"golinkcut/internal/config"
	"golinkcut/pkg/db/memory"
	"golinkcut/pkg/db/postgresql"
	"gorm.io/gorm"
	"log"
)

const (
	ErrAliasTakenSignature1 = "constraint failed: UNIQUE constraint failed: links.alias"
	ErrAliasTakenSignature2 = "duplicate key value violates unique constraint \"idx_links_alias\""
	ErrLinkExistsSignature1 = "constraint failed: UNIQUE constraint failed: links.original"
	ErrLinkExistsSignature2 = "duplicate key value violates unique constraint \"idx_links_original\""
	ErrNotExistsSignature   = "record not found"
)

// Link is a struct representing original link and it's short alias.
// The fields are set up to both have a unique value in DB
type Link struct {
	gorm.Model
	// Alias is the short link alias
	Alias string `gorm:"uniqueIndex":compositindex;type:text;not null`
	// Original is the original link to be found by Alias
	Original string `gorm:"uniqueIndex":compositindex;type:text;not null`
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
	return tx.Error
}

// Get is used to retrieve the original link, given its alias in DB
func (s Storage) Get(ctx context.Context, alias string) (string, error) {
	var link Link
	tx := s.db.WithContext(ctx).First(&link, "alias = ?", alias)
	if tx.Error != nil {
		return "", tx.Error
	}
	return link.Original, nil
}

func (s Storage) deleteData() {
	s.db.Exec("DROP TABLE links")
}

// NewStorage reads config.Config and returns a corresponding Storage object
func NewStorage(cfg config.Config) Storage {
	var db *gorm.DB
	log.Printf("Will connect to database now")
	useMemory := cfg["memory"].(bool)
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
