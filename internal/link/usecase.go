package link

import (
	"context"
	"golinkcut/internal/entity"
	"golinkcut/pkg/log"
	"math/rand"
	"time"
	"unsafe"
)

type UseCase interface {
	Get(ctx context.Context, shortLink string) (entity.Link, error)
	Create(ctx context.Context, input CreateLinkRequest) (entity.Link, error)
}

type usecase struct {
	repo   Repository
	logger log.Logger
}

func (uc usecase) Get(ctx context.Context, shortLink string) (entity.Link, error) {
	link, err := uc.repo.Get(ctx, shortLink)
	if err != nil {
		return entity.Link{}, err
	}
	return link, nil
}

// CreateLinkRequest represents a shortened link creation request
type CreateLinkRequest struct {
	OriginalLink string `json:"original_link"`
}

func (uc *usecase) Create(ctx context.Context, req CreateLinkRequest) (entity.Link, error) {
	shortLink := generateShortAlias()
	uc.repo.Save(ctx, req.Original)
}

func NewUseCase(repo Repository, logger log.Logger) UseCase {
	return usecase{repo: repo, logger: logger}
}

func generateShortAlias() string {
	return randStringBytesMaskImprSrcUnsafe(10)
}

// code below is mostly borrowed from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986 and adapted for use case

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrcUnsafe(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
