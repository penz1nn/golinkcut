// WARNING! Please successfully launch golinkcut/for-tests/postgres-compose.yml
// or these tests will be skipped.
package db

import (
	"context"
	"errors"
	"golinkcut/internal/config"
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
	cfg := config.Config{
		"memory": false,
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
	storage := NewStorage(cfg)
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
	var ref ErrNotExists
	if !errors.As(err, &ref) {
		t.Errorf("Wrong type of error. Got: %T, want: %T", err, ref)
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
	var ref ErrAliasExists
	if !errors.As(err, &ref) {
		t.Errorf("Wrong type of error. Got: %T, want: %T", err, ref)
	}
}
