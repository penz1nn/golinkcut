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
	// letterBytes represent the letters which are allowed to be in the short
	// link alias. This alphabet is also used to create a number in a numeric
	// system with the base of 63 (so there are 63 characters)
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

	// urlPattern represents a string for regex pattern to match a valid URL
	urlPattern = "(?sm)[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"

	// maxAlias is used to define a maximum number in a base 63 numeric system
	// which corresponds to a "maximum" alias from allowed characters
	maxAlias = 63*63*63*63*63*63*63*63*63*63 - 1
)

// urlRegex is a regexp pattern to check if a string is a valid URL
var urlRegex = regexp.MustCompile(urlPattern)

// UseCase is the interface which defines behavior to execute the main
// business logic of application - fetching created and creating  new short links
type UseCase interface {
	Get(ctx context.Context, shortLink string) (entity.Link, error)
	Create(ctx context.Context, input CreateLinkRequest) (entity.Link, error)
}

// usecase struct is used to implement the UseCase interface
type usecase struct {
	repo     Repository
	validate bool
}

// Get fetches the entity.Link instance, given it's alias (the key part of the short link)
func (uc usecase) Get(ctx context.Context, shortLink string) (entity.Link, error) {
	link, err := uc.repo.GetLink(ctx, shortLink)
	return link, err
}

// CreateLinkRequest represents a shortened link creation request
type CreateLinkRequest struct {
	OriginalLink string `json:"original_link"`
}

// Validate checks whether passed string of CreateLinkRequest is a valid URL
func (req CreateLinkRequest) Validate() bool {
	return urlRegex.MatchString(req.OriginalLink)
}

// Create receives a CreateLinkRequest to create and store the entity.Link
func (uc usecase) Create(ctx context.Context, req CreateLinkRequest) (entity.Link, error) {
	if uc.validate && !req.Validate() {
		return entity.Link{}, ErrBadUrl{Url: req.OriginalLink}
	}

	shortLink := generateShortAlias(req.OriginalLink)
	l := entity.Link{Alias: shortLink, Original: req.OriginalLink}
	err := uc.repo.SaveLink(ctx, l)
	return l, err
}

// NewUseCase function is used to create a new instance of usecase
func NewUseCase(repo Repository, config config.Config) UseCase {
	validate := false
	v, ok := config["validate"].(bool)
	if ok && v {
		validate = true
	}
	return usecase{repo: repo, validate: validate}
}

// generateShortAlias receives the url string to create a short 10-character
// alias string from allowed characters (see letterBytes)
func generateShortAlias(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	shaStr := hex.EncodeToString(h.Sum(nil)[12:])
	// NOTE: unhandled error here; test show the code is reliable
	shaHash, _ := strconv.ParseUint(shaStr, 16, 64)
	hash := shaHash % maxAlias
	return toBase63(hash)
}

// toBase63 receives a number to convert it to numeric system with the base of
// 63 and returns a corresponding string. The characters to create the number
// are taken from letterBytes.
// If a created number string contains less then 10 digits, leading zeroes are
// added.
func toBase63(num uint64) string {
	var result []rune
	for num > 0 {
		remainder := num % 63
		result = append(result, rune(letterBytes[remainder]))
		num = num / 63
	}
	for len(result) < 10 {
		for i := 0; i < 10-len(result); i++ {
			result = append(result, rune(letterBytes[0]))
		}
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

// ErrBadUrl represents an error which is thrown when invalid string is passed
// to be used as URL
type ErrBadUrl struct {
	Url string
}

func (e ErrBadUrl) Error() string {
	return fmt.Sprintf("Wrong url format for %s", e.Url)
}
