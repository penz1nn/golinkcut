package link

import (
	"context"
	"fmt"
	"golinkcut/internal/entity"
	"golinkcut/pkg/memory"
	"testing"
)

const (
	short1 = "_aBcDeF120"
	orig1  = "https://google.com/?s=sampletext"
)

func TestMemoryStorage_Save(t *testing.T) {
	ms := memoryStorage{memory.NewStorage()}
	link := entity.Link{
		Short:    short1,
		Original: orig1,
	}
	err := ms.Save(context.Background(), link)
	if err != nil {
		t.Errorf("Error %v when saving new link", err)
	}
}

func TestMemoryStorage_Get(t *testing.T) {
	ms := memoryStorage{memory.NewStorage()}
	link := entity.Link{
		Short:    short1,
		Original: orig1,
	}
	err := ms.Save(context.Background(), link)
	if err != nil {
		t.Errorf("Unexpected rror %v when saving new link", err)
	}
	linkGot, err := ms.Get(context.Background(), short1)
	if err != nil {
		t.Errorf("Unexpected error when getting link: %v", err)
	}
	if linkGot != link {
		t.Errorf("Got a wrong link:\nexpected:\t%v\ngot:\t\t%v", link, linkGot)
	}
}

func TestShortlinkTakenError_Error(t *testing.T) {
	ms := memoryStorage{memory.NewStorage()}
	link := entity.Link{
		Short:    short1,
		Original: orig1,
	}
	err := ms.Save(context.Background(), link)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	err = ms.Save(context.Background(), link)
	if err == nil {
		t.Error("Error expected but got nil")
	}
	got := err.Error()
	want := fmt.Sprintf("short link alias %s is already occupied by original link %s", short1, orig1)
	if got != want {
		t.Errorf("Wrong error string:\nexpected:\t%s\ngot:\t\t%s", want, got)
	}
}

func TestNotFoundError_Error(t *testing.T) {
	ms := memoryStorage{memory.NewStorage()}
	_, err := ms.Get(context.Background(), short1)
	if err == nil {
		t.Error("Error expected but got nil")
	}
	got := err.Error()
	want := fmt.Sprintf("short link alias %s not found in storage", short1)
	if got != want {
		t.Errorf("Wrong expected error:\nwant:\t%v\ngot:\t%v", want, got)
	}
	link := entity.Link{
		Short:    short1,
		Original: orig1,
	}
	err = ms.Save(context.Background(), link)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	_, err = ms.Get(context.Background(), "wrongAlias")
	if err == nil {
		t.Error("Error expected but got nil")
	}
	got = err.Error()
	want = fmt.Sprintf("short link alias %s not found in storage", "wrongAlias")
	if got != want {
		t.Errorf("Wrong expected error:\nwant:\t%v\ngot:\t%v", want, got)
	}
}
