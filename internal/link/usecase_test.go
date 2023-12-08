package link

import (
	"context"
	"golinkcut/internal/config"
	"golinkcut/pkg/log"
	"regexp"
	"testing"
)

const (
	url1 = "https://yandex.com/?s=sampletext3"
	url2 = "github.com/penz1nn/golinkcut"
	url3 = "ftp://192.168.1.1"
)

func TestCreateLinkRequest_Validate(t *testing.T) {
	req := CreateLinkRequest{url1}
	if !req.Validate() {
		t.Errorf("Valid url discarded: %s", req.OriginalLink)
	}
	req = CreateLinkRequest{url2}
	if !req.Validate() {
		t.Errorf("Valid url discarded: %s", req.OriginalLink)
	}
	req = CreateLinkRequest{url3}
	if !req.Validate() {
		t.Errorf("Valid url discarded: %s", req.OriginalLink)
	}
	req = CreateLinkRequest{"Not-a-link-at-all"}
	if req.Validate() {
		t.Errorf("Invalid url passed as valid: %s", req.OriginalLink)
	}
}

func TestNewUseCase(t *testing.T) {
	c := config.Config{
		"db":       "memory",
		"validate": false,
		"debug":    false,
	}
	repo := NewStorage(c)
	logger := log.NewWithConfig(c)
	uc := NewUseCase(repo, logger, c)
	if &uc == nil {
		t.Error("create &uc is not expected to be nil")
	}
}

func TestUsecase_Create(t *testing.T) {
	c := config.Config{
		"db":       "memory",
		"validate": false,
		"debug":    false,
	}
	repo := NewStorage(c)
	logger := log.NewWithConfig(c)
	uc := NewUseCase(repo, logger, c)
	url := "https://mos.ru"
	link, err := uc.Create(context.Background(), CreateLinkRequest{OriginalLink: url})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if link.Original != url {
		t.Errorf("Created ink with wrong url:\ngot:\t%v\nwant:\t%v", link.Original, url)
	}
	if !validateAlias(link.Alias) {
		t.Errorf("Invalid alias: %v", link.Alias)
	}
}

func TestUsecase_Get(t *testing.T) {
	c := config.Config{
		"db":       "memory",
		"validate": false,
		"debug":    false,
	}
	repo := NewStorage(c)
	logger := log.NewWithConfig(c)
	uc := NewUseCase(repo, logger, c)
	link1, err := uc.Create(context.Background(), CreateLinkRequest{OriginalLink: url1})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	link2, err := uc.Get(context.Background(), link1.Alias)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if link1 != link2 {
		t.Errorf("Retuned string is not the saved:\ngot:\t%v\nwant:\t%v", link2, link1)
	}
}

func TestErrBadUrl_Error(t *testing.T) {
	c := config.Config{
		"db":       "memory",
		"validate": true,
		"debug":    false,
	}
	repo := NewStorage(c)
	logger := log.NewWithConfig(c)
	uc := NewUseCase(repo, logger, c)
	_, err := uc.Create(context.Background(), CreateLinkRequest{OriginalLink: url1})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	_, err = uc.Create(context.Background(), CreateLinkRequest{OriginalLink: "some-wrong-text"})
	err2 := ErrBadUrl{Url: "some-wrong-text"}
	if err.Error() != err2.Error() {
		t.Errorf("Expected:%v\ngot:%v", err2, err)
	}
}

func TestGen(t *testing.T) {
	were := map[string]bool{}
	for i := 0; i < 1000; i++ {
		alias := generateShortAlias()
		if were[alias] {
			t.Errorf("Repeated alias %v. Please repeat the test and see if it fails more, if not, consider it safe", alias)
		}
		were[alias] = true
	}
}

func validateAlias(alias string) bool {
	if len(alias) != 10 {
		return false
	}
	r := regexp.MustCompile("^[A-Za-z0-9_]{10}$")
	if !r.MatchString(alias) {
		return false
	}
	return true
}
