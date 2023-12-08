package db

import (
	"context"
	"fmt"
	"golinkcut/internal/config"
	"golinkcut/pkg/db/memory"
	"golinkcut/pkg/db/postgresql"
	"gorm.io/gorm"
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
func NewStorage(config config.Config) Storage {
	var db *gorm.DB
	dbCfg, ok := config["db"].(string)
	if !ok {
		db = postgresql.NewDb(config)
	} else if dbCfg != "memory" {
		panic(fmt.Sprintf("unknown db config string: config[\"db\"] == %s", dbCfg))
	} else {
		db = memory.NewDb(config)
	}
	err := db.AutoMigrate(&Link{})
	if err != nil {
		panic(err)
	}
	return Storage{db}
}

/*
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
*/
