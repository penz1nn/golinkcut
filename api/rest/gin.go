package rest

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"net/http"
)

type restApi struct {
	uc       link.UseCase
	redirect bool
	host     string
}

type CreateLinkRequest struct {
	Url string `json:"url,omitempty" binding:"required"`
}

func (rest restApi) CreateLink(c *gin.Context) {
	var input CreateLinkRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := link.CreateLinkRequest{OriginalLink: input.Url}
	l, err := rest.uc.Create(c, req)
	if err != nil {
		var errLinkExists link.ErrLinkExists
		if errors.As(err, &errLinkExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		var errBadUrl link.ErrBadUrl
		if errors.As(err, &errBadUrl) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	shortLink := fmt.Sprintf("%s/%s", rest.host, l.Alias)
	c.JSON(http.StatusCreated, gin.H{"shortLink": shortLink})
}

func (rest restApi) GetLink(c *gin.Context) {
	alias := c.Request.URL.Path
	l, err := rest.uc.Get(c, alias)
	if err != nil {
		var errNotExists link.ErrNotExists
		if errors.As(err, &errNotExists) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rest.redirect {
		c.Redirect(http.StatusFound, l.Original)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"url": l.Original})
		return
	}
}

func NewRestApi(cfg config.Config, uc link.UseCase) restApi {
	rest := restApi{
		uc:       uc,
		redirect: cfg["redirect"].(bool),
		host:     cfg["httpHost"].(string),
	}
	return rest
}
