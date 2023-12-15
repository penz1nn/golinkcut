package link

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"golinkcut/internal/config"
	"golinkcut/internal/entity"
	"regexp"
	"strconv"
)

const (
	// letterBytes represent the letters which are allowed to be in the short link alias
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	urlPattern  = "(?sm)[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
	maxAlias    = 63*63*63*63*63*63*63*63*63*63 - 1
)

var urlRegex = regexp.MustCompile(urlPattern)

type UseCase interface {
	Get(ctx context.Context, shortLink string) (entity.Link, error)
	Create(ctx context.Context, input CreateLinkRequest) (entity.Link, error)
}

type usecase struct {
	repo     Repository
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
	if uc.validate && !req.Validate() {
		return entity.Link{}, ErrBadUrl{Url: req.OriginalLink}
	}

	shortLink := generateShortAlias(req.OriginalLink)
	l := entity.Link{Alias: shortLink, Original: req.OriginalLink}
	err := uc.repo.SaveLink(ctx, l)
	return l, err
}

func NewUseCase(repo Repository, config config.Config) UseCase {
	validate := false
	v, ok := config["validate"].(bool)
	if ok && v {
		validate = true
	}
	return usecase{repo: repo, validate: validate}
}

func generateShortAlias(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	sha_str := hex.EncodeToString(h.Sum(nil)[12:])
	sha_hash, err := strconv.ParseUint(sha_str, 16, 64)
	if err != nil {
		panic(err)
	}
	hash := sha_hash % maxAlias
	return toBase63(hash)
}

func toBase63(num uint64) string {
	var result []rune
	for num > 0 {
		remainder := num % 63
		result = append(result, rune(letterBytes[remainder]))
		num = num / 63
	}
	if len(result) < 10 {
		for i := 0; i < 10-len(result); i++ {
			result = append(result, rune(letterBytes[0]))
		}
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

type ErrBadUrl struct {
	Url string
}

func (e ErrBadUrl) Error() string {
	return fmt.Sprintf("Wrong url format for %s", e.Url)
}
