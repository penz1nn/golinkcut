package postgresql

import (
	"fmt"
	"golinkcut/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var (
	dbparams   = [6]string{"host", "user", "password", "dbname", "port", "tz"}
	dbsettings = [6]string{"host=%s ", "user=%s ", "password=%s ", "dbname=%s ", "port=%s ", "sslmode=disable TimeZone=%s"}
)

func NewDb(config config.Config) *gorm.DB {
	connStr := ""
	for index, param := range dbparams {
		if val, ok := config["db"].(map[string]string)[param]; ok {
			connStr = connStr + fmt.Sprintf(dbsettings[index], val)
		} else {
			errorStr := fmt.Sprintf("wrong DB connection data: no parameter %s given", param)
			panic(errorStr)
		}
	}

	cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	if debug, ok := config["debug"]; ok {
		if debug.(bool) {
			cfg = &gorm.Config{}
		}
	}
	db, err := gorm.Open(postgres.Open(connStr), cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database with %v", err))
	}
	log.Printf("Connected to PostgreSQL")
	return db
}
