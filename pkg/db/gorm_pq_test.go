// WARNING! Please successfully launch golinkcut/testdata/docker-compose.yml
// or these tests will be skipped.
package db

import (
	"context"
	"golinkcut/internal/config"
	"strings"
	"testing"
)

const (
	host     = "localhost"
	user     = "golinkcut_test"
	password = "example"
	dbname   = "golinkcut_test"
	port     = "5432"
	tz       = "Europe/Moscow"
)

func initStorage(t *testing.T) Storage {
	config := config.Config{
		"db": map[string]string{
			"host":     host,
			"user":     user,
			"password": password,
			"dbname":   dbname,
			"port":     port,
			"tz":       tz,
		},
	}
	defer func() {
		if recover() != nil {
			t.Skip("Could not connect to PostgreSQL, skipping")
		}
	}()
	storage := NewStorage(config)
	t.Cleanup(func() {
		storage.deleteData()
	})
	return storage
}

func TestStorage_Add2(t *testing.T) {
	s := initStorage(t)
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
func TestStorage_Get2(t *testing.T) {
	s := initStorage(t)
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

func TestErrAliasTaken_Error2(t *testing.T) {
	s := initStorage(t)
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = s.Add(context.Background(), short1, orig2)
	if err == nil {
		t.Error("Error is expected but it's nil")
	}
	if !strings.Contains(err.Error(), ErrAliasTakenSignature2) {
		t.Errorf("Wrong error format: %s", err.Error())
	}
}

func TestErrLinkExists_Error2(t *testing.T) {
	s := initStorage(t)
	err := s.Add(context.Background(), short1, orig1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = s.Add(context.Background(), short2, orig1)
	if err == nil {
		t.Error("Error is expected but it's nil")
	}
	if !strings.Contains(err.Error(), ErrLinkExistsSignature2) {
		t.Errorf("Wrong error format: %s", err.Error())
	}
}
