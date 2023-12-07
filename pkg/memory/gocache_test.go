package memory

import (
	"context"
	"testing"
)

func TestStorage_Set(t *testing.T) {
	s := NewStorage()
	s.SetKey(context.Background(), "key1", "value1")
	s.SetKey(context.Background(), "key2", "value2")
	s.SetKey(context.Background(), "key1", "value3")
}

func TestStorage_Get(t *testing.T) {
	s := NewStorage()
	k1 := "key1"
	k2 := "key2"
	v1 := "value1"
	v2 := "value2"
	v3 := "value3"
	s.SetKey(context.Background(), k1, v1)
	got := s.GetKey(context.Background(), k1)
	want := v1
	if got != want {
		t.Errorf("Expected: %v, got: %v", want, got)
	}
	s.SetKey(context.Background(), k2, v2)
	got = s.GetKey(context.Background(), k2)
	want = v2
	if got != want {
		t.Errorf("Expected: %v, got: %v", want, got)
	}
	s.SetKey(context.Background(), k1, v3)
	got = s.GetKey(context.Background(), k1)
	want = v3
	if got != want {
		t.Errorf("Expected: %v, got: %v", want, got)
	}
}
