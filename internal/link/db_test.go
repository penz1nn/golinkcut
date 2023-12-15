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
)

func TestNewStorage(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewDbStorage(c)
	if &s == nil {
		t.Error("dbStorage is not expected to be nil")
	}
}

func TestStorage_SaveLink(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewDbStorage(c)
	err := s.SaveLink(context.Background(), entity.Link{Alias: short1, Original: orig1})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestStorage_GetLink(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewDbStorage(c)
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

func TestErrNotExists_Error(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewDbStorage(c)
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
