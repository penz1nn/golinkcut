package db

import (
	"fmt"
	"golinkcut/internal/config"
	"golinkcut/pkg/db/memory"
	"gorm.io/gorm"
	"strings"
)

type Link struct {
	gorm.Model
	Alias    string `gorm:"uniqueIndex":compositindex;type:text;not null`
	Original string `gorm:"uniqueIndex":compositindex;type:text;not null`
}

type Storage struct {
	db *gorm.DB
}

func (s Storage) Add(alias string, original string) error {
	tx := s.db.Create(&Link{Alias: alias, Original: original})
	if tx.Error != nil {
		if strings.Contains(tx.Error.Error(), "constraint failed: UNIQUE constraint failed: links.alias") {
			return ErrAliasTaken{Alias: alias}
		}
		if strings.Contains(tx.Error.Error(), "constraint failed: UNIQUE constraint failed: links.original") {
			return ErrLinkExists{Original: original}
		}
	}
	return tx.Error
}

func (s Storage) Get(alias string) (string, error) {
	var link Link
	tx := s.db.First(&link, "alias = ?", alias)
	if tx.Error != nil {
		return "", tx.Error
	}
	return link.Original, nil
}

func NewStorage(config config.Config) Storage {
	// TODO add code for postgresql
	db := memory.NewDb(config)
	err := db.AutoMigrate(&Link{})
	if err != nil {
		panic(err)
	}
	return Storage{db}
}

type ErrAliasTaken struct {
	Alias string
}

func (e ErrAliasTaken) Error() string {
	return fmt.Sprintf("Alias %v already exists in db", e.Alias)
}

// TODO return code 409 "Conflict" when error occurs
type ErrLinkExists struct {
	Original string
}

func (e ErrLinkExists) Error() string {
	return fmt.Sprintf("Link %v already exists in db", e.Original)
}
