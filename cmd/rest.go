package main

import (
	"github.com/gin-gonic/gin"
	"golinkcut/api/rest"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"golinkcut/pkg/log"
)

func runRestApi(uc link.UseCase, logger log.Logger, cfg config.Config) {
	r := gin.Default()
	restApi := rest.NewRestApi(cfg, uc)
	if !cfg["debug"].(bool) {
		gin.SetMode(gin.ReleaseMode)
	}
	r.GET("/", restApi.GetLink)
	r.POST("/:alias", restApi.CreateLink)
	logger.Infof("REST API started at port %s", cfg["httpPort"].(string))
	if err := r.Run(":" + cfg["httpPort"].(string)); err != nil {
		logger.Errorf("REST API Error: %s", err)
	}
}
