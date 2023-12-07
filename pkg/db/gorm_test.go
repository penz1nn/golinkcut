package db

import (
	"fmt"
	"golinkcut/internal/config"
	"testing"
)

const (
	short1 = "_aBcDeF120"
	orig1  = "https://google.com/?s=sampletext"
	short2 = "0123456789"
	orig2  = "https://github.com/penz1nn"
)

func TestNewStorage(t *testing.T) {
	c := config.Config{"db": "memory"}
	s := NewStorage(c)
	if &s == nil {
		t.Error("&storage is not expected to be nil")
	}
}

func TestStorage_Add(t *testing.T) {
	c := config.Config{"db": "memory"}
	s := NewStorage(c)
	err := s.Add(short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestErrAliasTaken_Error(t *testing.T) {
	c := config.Config{"db": "memory"}
	s := NewStorage(c)
	err := s.Add(short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = s.Add(short1, orig2)
	if err == nil {
		t.Error("Error is expected but it's nil")
	}
	got := err.Error()
	want := fmt.Sprintf("Alias %v already exists in db", short1)
	if got != want {
		t.Errorf("Wrong error string:\nwant:\t%v\ngot:\t%v", want, got)
	}
}

func TestErrLinkExists_Error(t *testing.T) {
	c := config.Config{"db": "memory"}
	s := NewStorage(c)
	err := s.Add(short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = s.Add(short2, orig1)
	if err == nil {
		t.Error("Error is expected but it's nil")
	}
	got := err.Error()
	want := fmt.Sprintf("Link %v already exists in db", orig1)
	if got != want {
		t.Errorf("Wrong error string:\nwant:\t%v\ngot:\t%v", want, got)
	}
}
