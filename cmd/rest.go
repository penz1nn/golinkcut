package main

import (
	"golinkcut/api/rest"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"log"
)

func runRestApi(uc link.UseCase, cfg config.Config) {
	r := rest.SetupRouter(uc, cfg)
	log.Printf("REST API started at port %s", cfg["httpPort"].(string))
	if err := r.Run(":" + cfg["httpPort"].(string)); err != nil {
		log.Printf("REST API Error: %s", err)
	}
}
