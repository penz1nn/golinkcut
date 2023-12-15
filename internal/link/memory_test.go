package link

import (
	"context"
	"errors"
	"golinkcut/internal/entity"
	"testing"
)

func TestNewKvStorage(t *testing.T) {
	s := NewKvStorage()
	if &s == nil {
		t.Error("kvStorage is not expected to be nil")
	}
}

func TestKvStorage_SaveLink(t *testing.T) {
	s := NewKvStorage()
	l := entity.Link{Alias: short1, Original: orig1}
	err := s.SaveLink(context.Background(), l)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestKvStorage_GetLink(t *testing.T) {
	s := NewKvStorage()
	l := entity.Link{Alias: short1, Original: orig1}
	err := s.SaveLink(context.Background(), l)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestKvStorage_GetLink_Error(t *testing.T) {
	s := NewKvStorage()
	_, err := s.GetLink(context.Background(), short1)
	var ref ErrNotExists
	if !errors.As(err, &ref) {
		t.Errorf("Wrong error type. Got: %T, want: %T", err, ref)
	}
}

func TestDbStorage_SaveLink(t *testing.T) {
	s := NewKvStorage()
	l := entity.Link{Alias: short1, Original: orig1}
	err := s.SaveLink(context.Background(), l)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = s.SaveLink(context.Background(), l)
	var ref ErrLinkExists
	if !errors.As(err, &ref) {
		t.Errorf("Wrong error type. Got: %T, want: %T", err, ref)
	}
}
