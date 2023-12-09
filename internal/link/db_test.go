package link

import (
	"context"
	"golinkcut/internal/config"
	"golinkcut/internal/entity"
	"testing"
)

const (
	short1 = "_aBcDeF125"
	orig1  = "https://google.com/?s=sampletext2"
	short2 = "012345678_"
	orig2  = "https://github.com/penz1nn"
)

func TestNewStorage(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewStorage(c)
	if &s == nil {
		t.Error("storage is not expected to be nil")
	}
}

func TestStorage_SaveLink(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewStorage(c)
	err := s.SaveLink(context.Background(), entity.Link{Alias: short1, Original: orig1})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestStorage_GetLink(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewStorage(c)
	link := entity.Link{Alias: short1, Original: orig1}
	err := s.SaveLink(context.Background(), link)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	gotLink, err := s.GetLink(context.Background(), link.Alias)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if gotLink != link {
		t.Errorf("Got wrong data:\ngot:\t%v\nwant:\t%v", gotLink, link)
	}
}

func TestErrLinkExists_Error(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewStorage(c)
	link1 := entity.Link{Alias: short1, Original: orig1}
	link2 := entity.Link{Alias: short2, Original: orig1}
	err := s.SaveLink(context.Background(), link1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = s.SaveLink(context.Background(), link2)
	if err == nil {
		t.Error("Expected error but got nil")
	}
	want := ErrLinkExists{Original: link1.Original}
	got := err
	if want != got {
		t.Errorf("Wrong error:\ngot:\t%v\nwant:\t%v", got, want)
	}
}

func TestErrAliasTaken_Error(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewStorage(c)
	link1 := entity.Link{Alias: short1, Original: orig1}
	link2 := entity.Link{Alias: short1, Original: orig2}
	err := s.SaveLink(context.Background(), link1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = s.SaveLink(context.Background(), link2)
	if err == nil {
		t.Error("Expected error but got nil")
	}
	want := ErrAliasTaken{Alias: short1}
	got := err
	if want != got {
		t.Errorf("Wrong error:\ngot:\t%v\nwant:\t%v", got, want)
	}
}

func TestErrNotExists_Error(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewStorage(c)
	_, err := s.GetLink(context.Background(), short1)
	if err == nil {
		t.Error("Expected error but got nil")
	}
	want := ErrNotExists{alias: short1}
	got := err
	if want != got {
		t.Errorf("Wrong error:\ngot:\t%v\nwant:\t%v", got, want)
	}

	err = s.SaveLink(context.Background(), entity.Link{Alias: short1, Original: orig1})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	_, err = s.GetLink(context.Background(), short2)
	if err == nil {
		t.Error("Expected error but got nil")
	}
	want = ErrNotExists{alias: short2}
	got = err
	if want != got {
		t.Errorf("Wrong error:\ngot:\t%v\nwant:\t%v", got, want)
	}
}
