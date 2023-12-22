package link

import (
	"context"
	"errors"
	"golinkcut/internal/config"
	"math/rand"
	"regexp"
	"testing"
)

const (
	url1 = "https://yandex.com/?s=sampletext3"
	url2 = "github.com/penz1nn/golinkcut"
	url3 = "ftp://192.168.1.1"
)

var lettersTest = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+?./")

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
		"memory":   true,
		"validate": false,
		"debug":    false,
	}
	repo := NewStorage(c)
	uc := NewUseCase(repo, c)
	if &uc == nil {
		t.Error("create &uc is not expected to be nil")
	}
}

func TestUsecase_Create(t *testing.T) {
	c := config.Config{
		"memory":   true,
		"validate": false,
		"debug":    false,
	}
	repo := NewStorage(c)
	uc := NewUseCase(repo, c)
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
		"memory":   true,
		"validate": false,
		"debug":    false,
	}
	repo := NewStorage(c)
	uc := NewUseCase(repo, c)
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
		"memory":   true,
		"validate": true,
		"debug":    false,
	}
	repo := NewStorage(c)
	uc := NewUseCase(repo, c)
	_, err := uc.Create(context.Background(), CreateLinkRequest{OriginalLink: url1})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	_, err = uc.Create(context.Background(), CreateLinkRequest{OriginalLink: "some-wrong-text"})
	err2 := ErrBadUrl{Url: "some-wrong-text"}
	if err.Error() != err2.Error() {
		t.Errorf("Expected:%v\ngot:%v", err2, err)
	}
	var errBadUrl ErrBadUrl
	if !errors.As(err, &errBadUrl) {
		t.Errorf("Wrong error format: %T", err)
	}
}

func TestToBase63(t *testing.T) {
	got := toBase63(uint64(63*63*63*63*63*63*63*63*63*63 - 1))
	want := "__________"
	if got != want {
		t.Errorf("Wrong result. Got: %v, Want: %v", got, want)
	}
	var num1 uint64 = 63 * 63 * 63 * 63 * 63 * 63 * 63 * 63 * 63 * 62
	num1 = num1 + 17*63*63*63*63*63*63*63*63
	num1 = num1 + 14*63*63*63*63*63*63*63
	num1 = num1 + 21*63*63*63*63*63*63
	num1 = num1 + 21*63*63*63*63*63
	num1 = num1 + 24*63*63*63*63
	got = toBase63(num1)
	want = "_hello0000"
	if got != want {
		t.Errorf("Wrong result. Got: %v, Want: %v", got, want)
	}
}

func TestGenAlias(t *testing.T) {

	for i := 0; i < 1000000; i++ {
		b := make([]rune, rand.Intn(30))
		for i := range b {
			b[i] = lettersTest[rand.Intn(len(lettersTest))]
		}
		str := string(b)
		alias := generateShortAlias(str)
		if !validateAlias(alias) {
			t.Errorf("Wrong alias format: %s", alias)
		}
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
