// Package memory (SQLite / in-memory SQL database) is used for tests only, this code is otherwise deprecated
package memory

import (
	"github.com/glebarez/sqlite"
	"golinkcut/internal/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDb creates a new SQLite in-memory database and returns a connection to it
func NewDb(config config.Config) *gorm.DB {
	cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent), TranslateError: true}
	if debug, ok := config["debug"]; ok {
		if debug.(bool) {
			cfg = &gorm.Config{}
		}
	}
	db, err := gorm.Open(sqlite.Open(":memory:"), cfg)
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}
