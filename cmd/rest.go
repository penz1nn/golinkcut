package main

import (
	"golinkcut/api/rest"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"golinkcut/pkg/log"
)

func runRestApi(uc link.UseCase, logger log.Logger, cfg config.Config) {
	r := rest.SetupRouter(uc, cfg)
	logger.Infof("REST API started at port %s", cfg["httpPort"].(string))
	if err := r.Run(":" + cfg["httpPort"].(string)); err != nil {
		logger.Errorf("REST API Error: %s", err)
	}
}
