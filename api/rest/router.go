package rest

import (
	"github.com/gin-gonic/gin"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
)

func SetupRouter(uc link.UseCase, cfg config.Config) *gin.Engine {
	r := gin.Default()
	restApi := NewRestApi(cfg, uc)
	if !cfg["debug"].(bool) {
		gin.SetMode(gin.ReleaseMode)
	}
	r.GET("/:alias", restApi.GetLink)
	r.POST("/new", restApi.CreateLink)
	return r
}
