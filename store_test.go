package shorten

import (
	"testing"
)

func TestMemoryStore(t *testing.T) {
	s := NewMemoryStore()

	key := s.Add("value1")

	if v, ok := s.Load(key); v != "value1" || !ok {
		t.Errorf("s.Load(%d) = (%v, %v), want (value1, true)", key, v, ok)
	}

	if v, ok := s.Load(key + 1); ok {
		t.Errorf("s.Load(%d) did exist: %v", key+1, v)
	}
}
