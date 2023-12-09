package rest

import (
	"github.com/gin-gonic/gin"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
)

func SetupRouter(uc link.UseCase, cfg config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	restApi := NewRestApi(cfg, uc)
	r.GET("/:alias", restApi.GetLink)
	r.POST("/new", restApi.CreateLink)
	return r
}
