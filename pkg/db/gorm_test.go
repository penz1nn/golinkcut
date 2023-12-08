package db

import (
	"context"
	"golinkcut/internal/config"
	"strings"
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
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestErrAliasTaken_Error(t *testing.T) {
	c := config.Config{"db": "memory"}
	s := NewStorage(c)
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = s.Add(context.Background(), short1, orig2)
	if err == nil {
		t.Error("Error is expected but it's nil")
	}
	if !strings.Contains(err.Error(), ErrAliasTakenSignature1) {
		t.Errorf("Wrong error format: %s", err.Error())
	}
}

func TestErrLinkExists_Error(t *testing.T) {
	c := config.Config{"db": "memory"}
	s := NewStorage(c)
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = s.Add(context.Background(), short2, orig1)
	if err == nil {
		t.Error("Error is expected but it's nil")
	}
	if !strings.Contains(err.Error(), ErrLinkExistsSignature1) {
		t.Errorf("Wrong error format: %s", err.Error())
	}
}

func TestStorage_Get(t *testing.T) {
	c := config.Config{"db": "memory"}
	s := NewStorage(c)
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	original, err := s.Get(context.Background(), short1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if original != orig1 {
		t.Errorf("Got wrong data: got:%v, want:%v", original, orig1)
	}
	_, err = s.Get(context.Background(), short2)
	if err == nil {
		t.Error("An error expected but got nil")
	}
	if !strings.Contains(err.Error(), ErrNotExistsSignature) {
		t.Errorf("Substring %s not found in error text: %v", ErrNotExistsSignature, err)
	}
}
