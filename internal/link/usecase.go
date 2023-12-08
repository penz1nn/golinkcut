package link

import (
	"context"
	"errors"
	"fmt"
	"golinkcut/internal/config"
	"golinkcut/internal/entity"
	"golinkcut/pkg/log"
	"math/rand"
	"regexp"
)

const (
	// letterBytes represent the letters which are allowed to be in the short link alias
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"
	// n is the number of letters in the short string alias
	n          = 10
	urlPattern = "(?sm)[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
)

var urlRegex = regexp.MustCompile(urlPattern)

type UseCase interface {
	Get(ctx context.Context, shortLink string) (entity.Link, error)
	Create(ctx context.Context, input CreateLinkRequest) (entity.Link, error)
}

type usecase struct {
	repo     Repository
	logger   log.Logger
	validate bool
}

func (uc usecase) Get(ctx context.Context, shortLink string) (entity.Link, error) {
	link, err := uc.repo.GetLink(ctx, shortLink)
	return link, err
}

// CreateLinkRequest represents a shortened link creation request
type CreateLinkRequest struct {
	OriginalLink string `json:"original_link"`
}

func (req CreateLinkRequest) Validate() bool {
	if urlRegex.MatchString(req.OriginalLink) {
		return true
	}
	return false
}

func (uc usecase) Create(ctx context.Context, req CreateLinkRequest) (entity.Link, error) {
	link := entity.Link{}
	if uc.validate && !req.Validate() {
		return link, ErrBadUrl{Url: req.OriginalLink}
	}
	saved := false
	tries := 0
	var err error
	for !saved && tries < 100 {
		tries++
		shortLink := generateShortAlias()
		link = entity.Link{Alias: shortLink, Original: req.OriginalLink}
		err = uc.repo.SaveLink(ctx, link)
		if errors.Is(err, ErrAliasTaken{}) {
			uc.logger.Errorf("error when trying to save new link: %v, retrying", err)
			err = nil
		} else if err != nil {
			break
		}
		saved = true
	}
	if tries >= 100 {
		panic("could not save a link after 100 tries!")
	}
	return link, err
}

func NewUseCase(repo Repository, logger log.Logger, config config.Config) UseCase {
	validate := false
	v, ok := config["validate"].(bool)
	if ok && v {
		validate = true
	}
	return usecase{repo: repo, logger: logger, validate: validate}
}

// TODO: use hashing instead
func generateShortAlias() string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type ErrBadUrl struct {
	Url string
}

func (e ErrBadUrl) Error() string {
	return fmt.Sprintf("Wrong url format for %s", e.Url)
}
