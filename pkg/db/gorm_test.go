package db

import (
	"context"
	"errors"
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
	c := config.Config{"memory": true}
	s := NewStorage(c)
	if &s == nil {
		t.Error("&storage is not expected to be nil")
	}
}

func TestStorage_Add(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewStorage(c)
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestErrAliasTaken_Error(t *testing.T) {
	c := config.Config{"memory": true}
	s := NewStorage(c)
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = s.Add(context.Background(), short1, orig2)
	if err == nil {
		t.Error("Error is expected but it's nil")
	}
	var ref ErrAliasExists
	if !errors.As(err, &ref) {
		t.Errorf("Wrong error format: got %T, expected: %T", err, ref)
	}
}

func TestStorage_Get(t *testing.T) {
	c := config.Config{"memory": true}
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
	var ref ErrNotExists
	if !errors.As(err, &ref) {
		t.Errorf("Wrong type of error. Got: %T, want: %T", err, ref)
	}
}

func (s Storage) deleteData() {
	s.db.Exec("DROP TABLE links")
}
