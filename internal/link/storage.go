package link

import (
	"golinkcut/internal/config"
	"log"
)

func NewStorage(cfg config.Config) Repository {
	if cfg["memory"].(bool) {
		log.Print("Will use in memory key-value store")
		return NewKvStorage()
	} else {
		log.Print("Will use database")
		return NewDbStorage(cfg)
	}
}
